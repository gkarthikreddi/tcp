package pkg

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/gkarthikreddi/tcp/tools/cmdparser"
)

const (
	SHOW_TOPO   = 1
	ARP_HANDLER = 2
)

func buildGraph() *Graph {
	graph := CreateNewGraph("Topology")
	node1 := CreateGraphNode(graph, "r0")
	node2 := CreateGraphNode(graph, "r1")
	node3 := CreateGraphNode(graph, "r2")

	InsertLinkBetweenNodes(node1, node2, "eth00", "eth01", 1)
	InsertLinkBetweenNodes(node2, node3, "eth02", "eth03", 1)
	InsertLinkBetweenNodes(node3, node1, "eth05", "eth04", 1)

	NodeSetLbAddr(node1, "122.1.1.0")
	NodeSetIntfIpAddr(node1, "eth04", "40.1.1.1", 24)
	NodeSetIntfIpAddr(node1, "eth00", "20.1.1.1", 24)

	NodeSetLbAddr(node2, "122.1.1.1")
	NodeSetIntfIpAddr(node2, "eth02", "80.1.1.1", 24)
	NodeSetIntfIpAddr(node2, "eth01", "90.1.1.1", 24)

	NodeSetLbAddr(node3, "122.1.1.2")
	NodeSetIntfIpAddr(node3, "eth03", "70.1.1.1", 24)
	NodeSetIntfIpAddr(node3, "eth05", "10.1.1.1", 24)

	return graph
}

var graph = buildGraph()

func showTopology(param *cmdparser.Param, buff *cmdparser.SerBuff) bool {
	code := cmdparser.ExtractCmdCode(buff)
	buff = buff.Next

	switch code {
	case SHOW_TOPO:
		DumpGraph(graph)
		return true
	}
	return false
}

func arpHandler(param *cmdparser.Param, buff *cmdparser.SerBuff) bool {
	code := cmdparser.ExtractCmdCode(buff)
	buff = buff.Next

	switch code {
	case ARP_HANDLER:
		for curr := buff; curr != nil; curr = curr.Next {
			if curr.Data.Id == "node-name" {
				fmt.Println("Node name: ", curr.Data.Value)
			}
			if curr.Data.Id == "ip-addr" {
				fmt.Println("IP addr: ", curr.Data.Value)
			}
		}
	}
	return false
}

func validNodeName(str string) bool {
	for curr := graph.List; curr != nil; curr = curr.Next {
		if curr.Name == str {
			return true
		}
	}
	return false
}

func validIPAddr(str string) bool {
	addr := strings.Split(str, ".")
	if len(addr) != 4 {
		return false
	}
	for _, val := range addr {
		if i, err := strconv.Atoi(val); err == nil {
			if i > 256 || i < 0 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func InitNwCli() {
	cmdparser.InitLibcli()

	show := cmdparser.GetShowHook()
	run := cmdparser.GetRunHook()

	{
		var topo cmdparser.Param
		cmdparser.InitParam(&topo, // param
			cmdparser.CMD,                  // type of param
			"topology",                     // name of param, nil for leaf param
			showTopology,                   // callback handler
			nil,                            // validationn handler
			cmdparser.INVALID,              // leaftype
			"",                             // id of leaf, nil for cmd param
			"Dump entire network topology") // help string
		cmdparser.LibcliRegisterParam(show, &topo)
		cmdparser.SetParamCmdCode(&topo, SHOW_TOPO)
	}
	{
		var node cmdparser.Param
		cmdparser.InitParam(&node,
			cmdparser.CMD,
			"node",
			nil,
			nil,
			cmdparser.INVALID,
			"",
			"Given a node name and operation it performs that operation")
		cmdparser.LibcliRegisterParam(run, &node)

		{
			var nodeName cmdparser.Param
			cmdparser.InitParam(&nodeName,
				cmdparser.LEAF,
				"",
				nil,
				validNodeName,
				cmdparser.STRING,
				"node-name",
				"Name of a node in the topology")
			cmdparser.LibcliRegisterParam(&node, &nodeName)

			{
				var resolveArp cmdparser.Param
				cmdparser.InitParam(&resolveArp,
					cmdparser.CMD,
					"resolve-arp",
					nil,
					nil,
					cmdparser.INVALID,
					"",
					"resolves arp of a node")
				cmdparser.LibcliRegisterParam(&nodeName, &resolveArp)

				{
					var ipAddr cmdparser.Param
					cmdparser.InitParam(&ipAddr,
						cmdparser.LEAF,
						"",
						arpHandler,
						validIPAddr,
						cmdparser.STRING,
						"ip-addr",
						"Takes ip addr of node i.e loopback addr")
					cmdparser.LibcliRegisterParam(&resolveArp, &ipAddr)
					cmdparser.SetParamCmdCode(&ipAddr, ARP_HANDLER)
				}
			}

		}

	}
}
