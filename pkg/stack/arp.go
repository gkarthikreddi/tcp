package stack

import (
	"fmt"

	"github.com/gkarthikreddi/tcp/pkg/network"
	"github.com/gkarthikreddi/tcp/tools"
)

const (
	ARP_BROAD_REQ = 1
	ARP_RPLY      = 2
	ARP_MSG       = 806
)

type arpHeader struct {
	HardwareType    uint16 // Always 1 for ethernet cable
	ProtocolType    uint16 // 0x0800 for IPv4
	HardwareLength  uint8  // 6 for network.Mac Address
	ProtocolLength  uint8  // 4 for IPv4
	Operation       uint16 // request:1 or reply:2
	SrcMacAddr      [6]byte
	SrcProtocolAddr [4]byte
	DstMacAddr      [6]byte
	DstProtocolAddr [4]byte
}

func arpTableLookup(arpTable *network.ArpEntry, ip *network.Ip) *network.ArpEntry {
	for entry := arpTable; entry != nil; entry = entry.Next {
		if entry.IpAddr.Addr == ip.Addr {
			return entry
		}
	}

	return nil
}

func deleteArpTableEntry(node *network.Node, ip *network.Ip) {
	arpTable := network.GetNodeArpTable(node)

	for entry := arpTable; entry != nil; entry = entry.Next {
		if *entry.IpAddr == *ip {
			if entry.Prev != nil && entry.Next != nil {
				entry.Prev.Next = entry.Next
				entry.Next.Prev = entry.Prev
			} else if entry.Prev == nil && entry.Next == nil {
				network.AssignNodeArpTable(node, nil)
			} else if entry.Prev == nil {
				entry.Next.Prev = nil
				network.AssignNodeArpTable(node, entry.Next)
			} else {
				entry.Prev.Next = nil
			}
			break
		}
	}
}

func addArpTableEntry(node *network.Node, entry *network.ArpEntry) {
	arpTable := network.GetNodeArpTable(node)
	if arpTable == nil {
		network.AssignNodeArpTable(node, entry)
		return
	}

	oldEntry := arpTableLookup(arpTable, entry.IpAddr)
	if oldEntry != nil {
		if entry.MacAddr.Addr == oldEntry.MacAddr.Addr {
			return
		} else {
			deleteArpTableEntry(node, oldEntry.IpAddr)
			addArpTableEntry(node, entry)
			return
		}
	}

	entry.Next = arpTable.Next
	entry.Prev = arpTable
	arpTable.Next = entry
}

func updateArpTableFromArpReply(node *network.Node, arpReply *arpHeader, localIntf *network.Interface) {
	if arpReply.Operation != ARP_RPLY {
		return
	}

	entry := network.ArpEntry{IpAddr: &network.Ip{Addr: arpReply.SrcProtocolAddr},
		MacAddr: &network.Mac{Addr: arpReply.SrcMacAddr},
		Name:    localIntf.Name}

	addArpTableEntry(node, &entry)
}

func SendArpBroadcast(node *network.Node, outIntf *network.Interface, ip *network.Ip) error {
	var err error
	if outIntf == nil {
		outIntf, err = network.NodeGetMatchingSubnet(node, ip)
		if err != nil {
			return fmt.Errorf("No outgoing interface and also no matching subnet interfaces")
		}
	}

	etherFrame := ethernetHeader{SrcMacAddr: network.GetIntfMac(outIntf).Addr,
		EtherType: ARP_MSG,
		Fcs:       0, // You shouldn't do this!
	}

	arpFrame := arpHeader{HardwareType: 1,
		ProtocolType:    0x0800,
		HardwareLength:  6,
		ProtocolLength:  4,
		Operation:       ARP_BROAD_REQ,
		SrcMacAddr:      network.GetIntfMac(outIntf).Addr,
		SrcProtocolAddr: network.GetIntfIp(outIntf).Addr,
		DstProtocolAddr: ip.Addr}
	fillBroadcastAddr(&etherFrame.DstMacAddr)

	if err = assignPayload(&etherFrame, &arpFrame); err != nil {
		return err
	}

	if err = sendPkt(&etherFrame, outIntf); err != nil {
		return err
	}
	return nil
}

func processArpBroadcast(node *network.Node, localIntf *network.Interface, etherFrame *ethernetHeader) error {
	fmt.Printf("ARP braodcast msg recieved on interface %s of node %s\n", localIntf.Name, node.Name)

	if arpFrame, err := tools.ByteToStruct(etherFrame.Payload[:], arpHeader{}); err == nil {
		ip := arpFrame.DstProtocolAddr
		if ip == network.GetIntfIp(localIntf).Addr {
			sendArpReply(etherFrame, localIntf)
		}
	} else {
		return fmt.Errorf("Can't get arpFrame from ethernetFrame")
	}
	return nil
}

func sendArpReply(etherFrame *ethernetHeader, outIntf *network.Interface) error {
    var err error
    if arpFrame, err := tools.ByteToStruct(etherFrame.Payload[:], arpHeader{}); err == nil {
		arpReplyFrame := arpHeader{HardwareType: 1,
			ProtocolType:    0x0800,
			HardwareLength:  6,
			ProtocolLength:  4,
			Operation:       ARP_RPLY,
			SrcMacAddr:      network.GetIntfMac(outIntf).Addr,
			SrcProtocolAddr: network.GetIntfIp(outIntf).Addr,
			DstMacAddr:      arpFrame.SrcMacAddr,
			DstProtocolAddr: arpFrame.SrcProtocolAddr,
		}
		etherReplyFrame := ethernetHeader{DstMacAddr: arpFrame.SrcMacAddr,
			SrcMacAddr: network.GetIntfMac(outIntf).Addr,
			EtherType:  ARP_MSG,
			Fcs:        0,
		}

		if err = assignPayload(&etherReplyFrame, &arpReplyFrame); err == nil {
            if err = sendPkt(&etherReplyFrame, outIntf); err == nil {
                return nil
            }
		}
	} 
	return err
}

func processArpReply(node *network.Node, localIntf *network.Interface, etherFrame *ethernetHeader) error {
	fmt.Printf("ARP reply msg recieved on interface %s of node %s\n", localIntf.Name, node.Name)

	if arpFrame, err := tools.ByteToStruct(etherFrame.Payload[:], arpHeader{}); err == nil {
		updateArpTableFromArpReply(node, arpFrame, localIntf)
	} else {
		return err
	}
	return nil
}
