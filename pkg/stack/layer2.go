package stack

import (
	"fmt"
	"time"

	"github.com/gkarthikreddi/tcp/pkg/network"
	"github.com/gkarthikreddi/tcp/tools"
)

type ethernetHeader struct {
	DstMacAddr [6]byte
	SrcMacAddr [6]byte
	Tagged     *vlan8021qHeader
	EtherType  uint16
	Payload    [500]byte // Ideally it could be between 45 - 1500 bytes
	Fcs        uint32
}

func assignPayload(etherFrame *ethernetHeader, arpFrame *arpHeader) error {
	if msg, err := tools.StructToByte(arpFrame); err == nil {
		copy(etherFrame.Payload[:], msg)
		return nil
	} else {
		return err
	}
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
	if network.IsIntfIp(intf) {
		if ether.Tagged == nil && (ether.DstMacAddr == network.GetIntfMac(intf).Addr || isBroadcastAddr(ether.DstMacAddr)) {
			return true
		}
		return false
	}

	mode := network.GetIntfL2Mode(intf)
	if mode == network.ACCESS {
		vlan := network.GetIntfVlanMembership(intf)[0]
		if ether.Tagged != nil {
			if vlan == ether.Tagged.Id {
				return true
			}
			return false
		} else {
			if vlan != 0 {
				ether.Tagged = &vlan8021qHeader{TPID: 0x800, Id: vlan}
				return true
			}
			return false
		}
	} else if mode == network.TRUNK {
		if ether.Tagged != nil {
			for _, val := range network.GetIntfVlanMembership(intf) {
				if val == 0 {
					break
				}
				if val == ether.Tagged.Id {
					return true
				}
			}
		}
		return false
	}
	return false
}

func layer2FrameRecieve(node *network.Node, intf *network.Interface, etherFrame *ethernetHeader) {
	if !validL2Intf(intf, etherFrame) {
		return
	}

	if network.IsIntfIp(intf) {
		promotePktToLayer2(node, intf, etherFrame)
	} else if mode := network.GetIntfL2Mode(intf); mode == "access" || mode == "trunk" {
		l2switchReceiveFrame(intf, etherFrame)
	}
}

func promotePktToLayer2(node *network.Node, intf *network.Interface, etherFrame *ethernetHeader) {
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
		break
	case ETH_IP:
		promotePktToLayer3(node, intf, etherFrame)
		break
	}
}

func demotePktToLayer2(node *network.Node, nextHopIp *network.Ip, outIntf string, ipFrame *ipHeader, protocol uint16) error {
	if protocol == ETH_IP {
		etherFrame := &ethernetHeader{EtherType: ETH_IP}
		if msg, err := tools.StructToByte(ipFrame); err == nil {
			copy(etherFrame.Payload[:], msg)
			return l2ForwardIpPkt(node, nextHopIp, outIntf, etherFrame)
		} else {
			return fmt.Errorf("Can't assign IP payload into EtherFrame")
		}
	}
	return nil
}

func l2ForwardIpPkt(node *network.Node, nextHopIp *network.Ip, outIntf string, etherFrame *ethernetHeader) error {
	var intf *network.Interface
	var err error
	if outIntf != "NA" {
		if intf, err = network.GetIntfByIntfName(node, outIntf); err != nil {
			return err
		}
	} else {
		if isLocalDelivery(node, nextHopIp) {
			promotePktToLayer3(node, nil, etherFrame)
			return nil
		} else {
			if intf, err = network.NodeGetMatchingSubnet(node, nextHopIp); err != nil {
				return err
			}
		}
	}

	entry := arpTableLookup(network.GetNodeArpTable(node), nextHopIp)
	if entry == nil {
		go SendArpBroadcast(node, intf, nextHopIp)
		time.Sleep(time.Millisecond * 100)
		entry = arpTableLookup(network.GetNodeArpTable(node), nextHopIp)
		if entry == nil {
			fmt.Println("quit")
			return nil
		}
	}
	// if entry != nil {
	etherFrame.DstMacAddr = entry.MacAddr.Addr
	etherFrame.SrcMacAddr = network.GetIntfMac(intf).Addr
	sendPkt(etherFrame, intf)
	// } else {
	//        return fmt.Errorf("no entry in arp")
	//    }

	return nil
}
