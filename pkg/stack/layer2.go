package stack

import (
	"github.com/gkarthikreddi/tcp/pkg/network"
	"github.com/gkarthikreddi/tcp/tools"
)

type ethernetHeader struct {
	DstMacAddr [6]byte
	SrcMacAddr [6]byte
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
	if network.IsIntfIp(intf) && (ether.DstMacAddr == network.GetIntfMac(intf).Addr || isBroadcastAddr(ether.DstMacAddr)) {
		return true
	}
	return false
}

func layer2FrameRecieve(node *network.Node, intf *network.Interface, etherFrame *ethernetHeader) {
	if !validL2Intf(intf, etherFrame) {
		return
	}

	if network.IsIntfIp(intf) {
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
	} else if mode := network.GetIntfL2Mode(intf); mode == "access" || mode == "trunk" {
		l2switchReceiveFrame(intf, etherFrame)
	}
}
