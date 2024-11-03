package cli

import (
	"fmt"
	"github.com/gkarthikreddi/tcp/pkg/network"
	"github.com/gkarthikreddi/tcp/pkg/stack"
	"github.com/gkarthikreddi/tcp/tools"
	"github.com/gkarthikreddi/tcp/tools/cmdparser"
)

func showHandler(param *cmdparser.Param, buff *cmdparser.SerBuff) bool {
	code := cmdparser.ExtractCmdCode(buff)
	buff = buff.Next

	var node *network.Node
	if buff.Data.Id == "node-name" {
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
