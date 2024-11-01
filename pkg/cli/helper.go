package cli

import (
	"fmt"
	"os"
    "github.com/jedib0t/go-pretty/v6/table"
	"github.com/gkarthikreddi/tcp/pkg/network"
	"github.com/gkarthikreddi/tcp/pkg/stack"
	"github.com/gkarthikreddi/tcp/tools"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
)

/*
	                 +----------+
	             0/4 |          |0/0
	+----------------+   R0     +---------------------------+
	|     40.1.1.1/24| 122.1.1.0|20.1.1.1/24                |
	|                +----------+                           |
	|                                                       |
	|                                                       |
	|                                                       |
	|40.1.1.2/24                                            |20.1.1.2/24
	|0/5                                                    |0/1

+---+---+                                              +----+-----+
|       |0/3                                        0/2|          |
| R2    +----------------------------------------------+    R1    |
|       |30.1.1.2/24                        30.1.1.1/24|          |
+-------+                                              +----------+
*/
func buildGraph() *network.Graph {
	graph := network.CreateNewGraph("Topology")
	node1 := network.CreateGraphNode(graph, "r0")
	node2 := network.CreateGraphNode(graph, "r1")
	node3 := network.CreateGraphNode(graph, "r2")

	network.InsertLinkBetweenNodes(node1, node2, "eth00", "eth01", 1)
	network.InsertLinkBetweenNodes(node2, node3, "eth02", "eth03", 1)
	network.InsertLinkBetweenNodes(node3, node1, "eth05", "eth04", 1)

	network.NodeSetLbAddr(node1, "122.1.1.0")
	network.NodeSetIntfIpAddr(node1, "eth04", "40.1.1.1", 24)
	network.NodeSetIntfIpAddr(node1, "eth00", "20.1.1.1", 24)

	network.NodeSetIntfIpAddr(node2, "eth02", "30.1.1.1", 24)
	network.NodeSetIntfIpAddr(node2, "eth01", "20.1.1.2", 24)

	network.NodeSetIntfIpAddr(node3, "eth03", "30.1.1.2", 24)
	network.NodeSetIntfIpAddr(node3, "eth05", "40.1.1.2", 24)

	stack.InitNetworkListening(graph)
	return graph
}

var graph = buildGraph()

func dumpGraph(graph *network.Graph) {
	fmt.Println("Name: " + Cyan + graph.Name + Reset)
	for curr := graph.List; curr != nil; curr = curr.Next {
		dumpNode(curr)
	}
}

func dumpNode(node *network.Node) {
	fmt.Println("Node name: " + node.Name + "\nLb addr: " + Yellow + tools.ConvertAddrToStr(network.GetNodeIp(node).Addr[:]) + Reset)
	for i := 0; node.Intf[i] != nil; i++ {
		dumpInterface(node.Intf[i])
	}
}

func dumpInterface(intf *network.Interface) {
	fmt.Println("\tInterface name: " + Cyan + intf.Name + Reset)
	nbrNode, _ := network.GetNbrNode(intf)
	fmt.Println("\t\tLocalNode: " + Cyan + intf.Att_node.Name + Reset + ", Nbr Node: " + Cyan + nbrNode.Name + Reset)
	fmt.Println("\t\tIp addr: " + Yellow + tools.ConvertAddrToStr(network.GetIntfIp(intf).Addr[:]) + Reset + " Mac addr: " + Yellow + tools.ConvertAddrToStr(network.GetIntfMac(intf).Addr[:]) + Reset)
}

func dumpArp(node *network.Node) {
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.AppendHeader(table.Row{"IP", "MAC", "Interface"})
	for curr := network.GetNodeArpTable(node); curr != nil; curr = curr.Next {
        t.AppendRow(table.Row{
        tools.ConvertAddrToStr(curr.IpAddr.Addr[:]),
        tools.ConvertAddrToStr(curr.MacAddr.Addr[:]),
        curr.Name })
	}
    t.Render()
}
