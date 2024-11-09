package stack

import "github.com/gkarthikreddi/tcp/pkg/network"

type vlan8021qHeader struct {
	TPID uint16
	Id   uint16 // This should be 12 bits has there are PRI (3bits) and CFI (1bits)
}

func addMacTableEntry(node *network.Node, macEntry *network.MacEntry) {
	macTable := network.GetNodeMacTable(node)
	if macTable == nil {
		network.AssignNodeMacTable(node, macEntry)
		return
	}
	if oldEntry := macTableLookup(macTable, macEntry.MacAddr.Addr); oldEntry != nil {
		if oldEntry.Name == macEntry.Name {
			return
		} else {
			deleteMacTableEntry(node, oldEntry.MacAddr.Addr)
			addMacTableEntry(node, macEntry)
			return
		}
	}
	macEntry.Next = macTable.Next
	macEntry.Prev = macTable
	macTable.Next = macEntry
}

func macTableLookup(macTable *network.MacEntry, macAddr [6]byte) *network.MacEntry {
	for entry := macTable; entry != nil; entry = entry.Next {
		if entry.MacAddr.Addr == macAddr {
			return entry
		}
	}

	return nil
}

func deleteMacTableEntry(node *network.Node, macAddr [6]byte) {
	macTable := network.GetNodeMacTable(node)
	for entry := macTable; entry != nil; entry = entry.Next {
		if entry.MacAddr.Addr == macAddr {
			if entry.Prev != nil && entry.Next != nil {
				entry.Prev.Next = entry.Next
				entry.Next.Prev = entry.Prev
			} else if entry.Prev == nil && entry.Next == nil {
				network.AssignNodeMacTable(node, nil)
			} else if entry.Prev == nil {
				entry.Next.Prev = nil
				network.AssignNodeMacTable(node, entry.Next)
			} else {
				entry.Prev.Next = nil
			}
		}
		break
	}
}

func l2switchReceiveFrame(localIntf *network.Interface, etherFrame *ethernetHeader) {
	node := localIntf.Att_node
	srcMac := etherFrame.SrcMacAddr

	// Perfrom Mac Learning
	macEntry := network.MacEntry{MacAddr: &network.Mac{Addr: srcMac}, Name: localIntf.Name}
	addMacTableEntry(node, &macEntry)

	// Forward frame
	l2switchForwardFrame(node, localIntf, etherFrame)
}

func l2switchForwardFrame(node *network.Node, localIntf *network.Interface, etherFrame *ethernetHeader) {
	if isBroadcastAddr(etherFrame.DstMacAddr) {
		l2sendPktFlood(node, localIntf, etherFrame)
		return
	}
	if macEntry := macTableLookup(network.GetNodeMacTable(node), etherFrame.DstMacAddr); macEntry != nil {
		if intf, _ := network.GetIntfByIntfName(node, macEntry.Name); intf != nil {
			l2switchSendPkt(etherFrame, intf)
			return
		}
	}
	l2sendPktFlood(node, localIntf, etherFrame)
}

func l2sendPktFlood(node *network.Node, excludeIntf *network.Interface, etherFrame *ethernetHeader) {
	for i := 0; i < network.MAX_INTF_PER_NODE; i++ {
		intf := node.Intf[i]
		if intf != nil && intf != excludeIntf {
			l2switchSendPkt(etherFrame, intf)
		}
	}
}

func l2switchSendPkt(etherFrame *ethernetHeader, outintf *network.Interface) {
	if network.IsIntfIp(outintf) {
		return
	}

	mode := network.GetIntfL2Mode(outintf)
	if mode == network.ACCESS {
		vlan := network.GetIntfVlanMembership(outintf)[0]
		if etherFrame.Tagged == nil {
			if vlan == 0 {
				sendPkt(etherFrame, outintf)
			}
		} else {
			if etherFrame.Tagged.Id == vlan {
				// untagging the frame before forwading it
				untaggedFrame := *etherFrame
				untaggedFrame.Tagged = nil
				sendPkt(&untaggedFrame, outintf)
			}
		}
	} else if mode == network.TRUNK && etherFrame.Tagged != nil {
		for _, val := range network.GetIntfVlanMembership(outintf) {
			if val == 0 {
				break
			}
			if val == etherFrame.Tagged.Id {
				sendPkt(etherFrame, outintf)
			}
		}
	}
}
