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
	BROADCAST_MAC = 0xFFFFFFFFFFFF
)

type ethernetHeader struct {
	DstMacAddr [6]byte
	SrcMacAddr [6]byte
	EtherType  uint16
	Payload    [500]byte // Ideally it could be between 45 - 1500 bytes
	Fcs        uint32
}

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

func assignPayload(ether *ethernetHeader, msg []byte) {
	copy(ether.Payload[:], msg)
}

func fillBroadcastAddr(mac *[6]byte) {
	for i := 0; i < 6; i++ {
		(*mac)[i] = 255
	}
}

func isBroadcastAddr(mac [6]byte) bool {
	return mac == [6]byte{255, 255, 255, 255, 255, 255}
}

func validL2Intf(intf *network.Interface, ether *ethernetHeader) bool {
	if network.IsIntfIp(intf) && (ether.DstMacAddr == network.GetIntfMac(intf).Addr || isBroadcastAddr(ether.DstMacAddr)) {
		return true
	}
	return false
}

func arpTableLookup(arpTable *network.ArpEntry, ip *network.Ip) *network.ArpEntry {
	for entry := arpTable; entry != nil; entry = entry.Next {
		if *entry.IpAddr == *ip {
			return entry
		}
	}

	return nil
}

func deleteArpTableEntry(arpTable *network.ArpEntry, ip *network.Ip) {
	for entry := arpTable; entry != nil; entry = entry.Next {
		if *entry.IpAddr == *ip {
			entry.Prev.Next = entry.Next
			entry.Next.Prev = entry.Prev
			return
		}
	}
}

func addArpTableEntry(arpTable *network.ArpEntry, entry *network.ArpEntry) {
	oldEntry := arpTableLookup(arpTable, entry.IpAddr)
	if oldEntry != nil {
		if entry.MacAddr.Addr == oldEntry.MacAddr.Addr {
			return
		} else {
			deleteArpTableEntry(arpTable, entry.IpAddr)
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

	if arpTable := network.GetNodeArpTable(node); arpTable == nil {
		network.AssignNodeArpTable(node, &entry)
	} else {
		addArpTableEntry(arpTable, &entry)
	}
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

	if arr, err := tools.StructToByte(arpFrame); err == nil {
		copy(etherFrame.Payload[:], arr)
	} else {
		return err
	}

	if err := sendPkt(&etherFrame, outIntf); err != nil {
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
		if arr, err := tools.StructToByte(arpReplyFrame); err == nil {
			copy(etherReplyFrame.Payload[:], arr)
		} else {
			return err
		}

		sendPkt(&etherReplyFrame, outIntf)
	} else {
		return err
	}
	return nil
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

func layer2FrameRecieve(node *network.Node, intf *network.Interface, etherFrame *ethernetHeader) {
	if !validL2Intf(intf, etherFrame) {
		return
	}

	switch etherFrame.EtherType {
	case ARP_MSG:
		if arpFrame, err := tools.ByteToStruct(etherFrame.Payload[:], arpHeader{}); err == nil {
			switch arpFrame.Operation {
			case ARP_BROAD_REQ:
				processArpBroadcast(node, intf, etherFrame)
				break
			case ARP_RPLY:
				processArpReply(node, intf, etherFrame)
				break
			}
		}
	}
}
