package cli

import (
	"fmt"
	"strconv"

	"github.com/gkarthikreddi/tcp/pkg/network"
	"github.com/gkarthikreddi/tcp/pkg/stack"
	"github.com/gkarthikreddi/tcp/tools"
	"github.com/gkarthikreddi/tcp/tools/cmdparser"
)

func showHandler(param *cmdparser.Param, buff *cmdparser.SerBuff) bool {
	code := cmdparser.ExtractCmdCode(buff)
	buff = buff.Next

	var node *network.Node
	if buff != nil && buff.Data.Id == "node-name" {
		node, _ = network.GetNodeByNodeName(graph, buff.Data.Value)
	}
	switch code {
	case SHOW_TOPO:
		dumpGraph(graph)
		return true
	case ARP_TABLE:
		dumpArpTable(node)
		return true
	case MAC_TABLE:
		dumpMacTable(node)
		return true
	case RT_TABLE:
		dumpRoutingTable(node)
		return true
	}
	return false
}

func arpHandler(param *cmdparser.Param, buff *cmdparser.SerBuff) bool {
	code := cmdparser.ExtractCmdCode(buff)
	buff = buff.Next

	switch code {
	case ARP_HANDLER:
		var node *network.Node
		var ip *network.Ip
		for curr := buff; curr != nil; curr = curr.Next {
			if curr.Data.Id == "node-name" {
				node, _ = network.GetNodeByNodeName(graph, curr.Data.Value)
			}
			if curr.Data.Id == "ip-addr" {
				ip = &network.Ip{Addr: tools.ConvertStrToIp(curr.Data.Value)}
			}
		}
		err := stack.SendArpBroadcast(node, nil, ip)
		if err != nil {
			fmt.Println(err)
		}
		return true
	}
	return false
}

func l3ConfigHandler(param *cmdparser.Param, buff *cmdparser.SerBuff) bool {
	code := cmdparser.ExtractCmdCode(buff)
	buff = buff.Next

	switch code {
	case L3_HANDLER:
		var node *network.Node
		var dstIp string
		var mask uint8
		var gatewayIp string
		var outIntf string

		for curr := buff; curr != nil; curr = curr.Next {
			if curr.Data.Id == "node-name" {
				node, _ = network.GetNodeByNodeName(graph, curr.Data.Value)
			} else if curr.Data.Id == "dst" {
				dstIp = curr.Data.Value
			} else if curr.Data.Id == "mask" {
				num, _ := strconv.Atoi(curr.Data.Value)
				mask = uint8(num)
			} else if curr.Data.Id == "gw-ip" {
				gatewayIp = curr.Data.Value
			} else if curr.Data.Id == "out-intf" {
				outIntf = curr.Data.Value
			}
		}

		if intf, err := network.GetIntfByIntfName(node, outIntf); err == nil {
			if network.IsIntfIp(intf) {
				ip := network.Ip{Addr: tools.ConvertStrToIp(dstIp), Mask: mask}
				tmp := network.Ip{Addr: tools.ConvertStrToIp(gatewayIp)}
				entry := network.RoutEntry{DstIpAddr: &ip,
					IsDirect:  false,
					GatewayIp: &tmp,
					OutIntf:   outIntf}

				stack.AddRoutingTableEntry(node, &entry)
				return true
			}
		}
	}
	return false
}

func pingHandler(param *cmdparser.Param, buff *cmdparser.SerBuff) bool {
	code := cmdparser.ExtractCmdCode(buff)
	buff = buff.Next
	switch code {
	case PING_HANDLER:
		var node *network.Node
		var dstIp [4]byte

		for curr := buff; curr != nil; curr = curr.Next {
			if curr.Data.Id == "node-name" {
				node, _ = network.GetNodeByNodeName(graph, curr.Data.Value)
			} else if curr.Data.Id == "ip-addr" {
				dstIp = tools.ConvertStrToIp(curr.Data.Value)
			}
		}

        stack.Ping(node, dstIp)
	}
	return false
}
