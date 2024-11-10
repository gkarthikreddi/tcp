package stack

import (
	"fmt"

	"github.com/gkarthikreddi/tcp/pkg/network"
	"github.com/gkarthikreddi/tcp/tools"
)

func AddRoutingTableEntry(node *network.Node, routEntry *network.RoutEntry) {
	routingTable := network.GetNodeRoutingTable(node)
	if routingTable == nil {
		network.AssignNodeRoutingTable(node, routEntry)
	} else {
		if oldEntry := findDuplicateEntry(routingTable, routEntry.DstIpAddr); oldEntry != nil {
			oldEntry.IsDirect = routEntry.IsDirect
			oldEntry.GatewayIp = routEntry.GatewayIp
			oldEntry.OutIntf = routEntry.OutIntf
		} else {
			routEntry.Next = routingTable.Next
			routEntry.Prev = routingTable
			routingTable.Next = routEntry
		}
	}
}

func findDuplicateEntry(routingTable *network.RoutEntry, dstIp *network.Ip) *network.RoutEntry {
	for entry := routingTable; entry != nil; entry = entry.Next {
		if entry.DstIpAddr.Addr == dstIp.Addr {
			return entry
		}
	}
	return nil
}

func routingTableLookup(routingTable *network.RoutEntry, dstIp *network.Ip) *network.RoutEntry {
	var ans *network.RoutEntry
	var lpm uint8 // Longest Prefix Match
	lpm = 0
	for entry := routingTable; entry != nil; entry = entry.Next {
        dstIp.Mask = entry.DstIpAddr.Mask
		ip := network.ApplyMask(dstIp)
		if ip == entry.DstIpAddr.Addr && entry.DstIpAddr.Mask > lpm {
			lpm = entry.DstIpAddr.Mask
			ans = entry
		}
	}
	return ans
}

func InitRoutingTable(graph *network.Graph) {
	for node := graph.List; node != nil; node = node.Next {
		if network.IsNodeIp(node) {
			routEntry := network.RoutEntry{DstIpAddr: network.GetNodeIp(node), IsDirect: true, GatewayIp: nil, OutIntf: "NA"}
			network.AssignNodeRoutingTable(node, &routEntry)
		}
		for _, intf := range node.Intf {
			if intf == nil {
				break
			}
			if network.IsIntfIp(intf) {
				newEntry := network.RoutEntry{DstIpAddr: network.GetIntfIp(intf), IsDirect: true, GatewayIp: nil, OutIntf: "NA"}
				AddRoutingTableEntry(node, &newEntry)
			}
		}
	}
}

func promotePktToLayer3(node *network.Node, intf *network.Interface, etherFrame *ethernetHeader) error {
	switch etherFrame.EtherType {
	case ETH_IP:
		if ipFrame, err := tools.ByteToStruct(etherFrame.Payload[:], ipHeader{}); err == nil {
			l3recieveFrame(node, intf, ipFrame)
		} else {
			return fmt.Errorf("Error while extracting IP payload from etherFrame")
		}
		break
	}
	return nil
}

func demotePktToLayer3(node *network.Node, dstIp *network.Ip, protocol uint8) error {
	ipFrame := newIpHeader()
	ipFrame.Protocol = protocol
	ipFrame.SrcIpAddr = network.GetNodeIp(node).Addr
	ipFrame.DstIpAddr = dstIp.Addr

	routingTable := network.GetNodeRoutingTable(node)
	var nextHopIp *network.Ip
	if route := routingTableLookup(routingTable, dstIp); route != nil {
		if isDirectRoute(route) {
			nextHopIp = dstIp
		} else {
			nextHopIp = route.GatewayIp
		}
        return demotePktToLayer2(node, nextHopIp, route.OutIntf, &ipFrame, ETH_IP)
	} else {
        return fmt.Errorf("Coudn't transfer packet received from application layer")
    }

}

func l3recieveFrame(node *network.Node, intf *network.Interface, ipFrame *ipHeader) error {
	routingTable := network.GetNodeRoutingTable(node)
	ip := &network.Ip{Addr: ipFrame.DstIpAddr}
	if route := routingTableLookup(routingTable, ip); route != nil {
		if isDirectRoute(route) {
			if isLocalDelivery(node, ip) {
				switch ipFrame.Protocol {
				case ICMP_PRO:
					fmt.Println("Ip Addr: " + tools.ConvertAddrToStr(ipFrame.DstIpAddr[:]) + " ping successful")
					break
				}
			} else {
				demotePktToLayer2(node, nil, "", ipFrame, ETH_IP)
			}
		} else {
			ipFrame.TTL -= 1
			if ipFrame.TTL == 0 {
				return fmt.Errorf("Max TTL reached")
			}

			demotePktToLayer2(node, route.GatewayIp, route.OutIntf, ipFrame, ETH_IP)
		}
	} else {
		return fmt.Errorf("Cound't forward packet has there is no route in the routing table")
	}
	return nil
}

func isDirectRoute(route *network.RoutEntry) bool {
	return route.IsDirect
}

func isLocalDelivery(node *network.Node, dstIp *network.Ip) bool {
	if network.GetNodeIp(node).Addr == dstIp.Addr {
		return true
	}
	for _, intf := range node.Intf {
		if network.GetIntfIp(intf).Addr == dstIp.Addr {
			return true
		}
	}

	return false
}
