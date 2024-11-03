package cli

import (
	"fmt"
	"github.com/gkarthikreddi/tcp/pkg/network"
	"github.com/gkarthikreddi/tcp/pkg/stack"
	"github.com/gkarthikreddi/tcp/tools"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
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

func buildGraph() *network.Graph {
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

func buildL2SwitchGraph() *network.Graph{
    /*
                                       +-----------+
                                       |  H4       |
                                       | 122.1.1.4 |
                                       +----+------+
                                            |eth0/7 - 10.1.1.3/24       
                                            |       
                                            |eth0/1
                                       +----+----+                        +---------+
       +---------+                     |         |                        |         |
       |         |10.1.1.2/24          |   L2Sw  |eth0/2       10.1.1.1/24|  H3     |
       |  H1     +---------------------+         +------------------------+122.1.1.3|
       |122.1.1.1|eth0/5         eth0/4|         |                 eth0/6 |         |
       + --------+                     |         |                        |         |
                                       +----+----+                        +---------+
                                            |eth0/3     
                                            |
                                            |
                                            |
                                            |10.1.1.4/24
                                            |eth0/8
                                      +----++------+
                                      |            |
                                      |   H2       |
                                      |122.1.1.2   |
                                      |            |
                                      +------------+


    */
	topo := network.CreateNewGraph("Simpel L2 switch demo graph")
	h1 := network.CreateGraphNode(topo, "H1")
	h2 := network.CreateGraphNode(topo, "H2")
	h3 := network.CreateGraphNode(topo, "H3")
	h4 := network.CreateGraphNode(topo, "H4")
	l2sw := network.CreateGraphNode(topo, "L2SW")

	network.InsertLinkBetweenNodes(h1, l2sw, "eth0/5", "eth0/4", 1)
	network.InsertLinkBetweenNodes(h2, l2sw, "eth0/8", "eth0/3", 1)
	network.InsertLinkBetweenNodes(h3, l2sw, "eth0/6", "eth0/2", 1)
	network.InsertLinkBetweenNodes(h4, l2sw, "eth0/7", "eth0/1", 1)

	network.NodeSetLbAddr(h1, "122.1.1.1")
	network.NodeSetIntfIpAddr(h1, "eth0/5", "10.1.1.2", 24)

	network.NodeSetLbAddr(h2, "122.1.1.2")
	network.NodeSetIntfIpAddr(h2, "eth0/8", "10.1.1.4", 24)

	network.NodeSetLbAddr(h3, "122.1.1.3")
	network.NodeSetIntfIpAddr(h3, "eth0/6", "10.1.1.1", 24)

	network.NodeSetLbAddr(h4, "122.1.1.4")
	network.NodeSetIntfIpAddr(h4, "eth0/5", "10.1.1.3", 24)

	network.NodeSetIntfL2Mode(l2sw, "eth0/1", network.ACCESS)
	network.NodeSetIntfL2Mode(l2sw, "eth0/2", network.ACCESS)
	network.NodeSetIntfL2Mode(l2sw, "eth0/3", network.ACCESS)
	network.NodeSetIntfL2Mode(l2sw, "eth0/4", network.ACCESS)

    stack.InitNetworkListening(topo)
    return topo
}

var graph = buildL2SwitchGraph()

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

func dumpArpTable(node *network.Node) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"IP", "MAC", "Interface"})
	for curr := network.GetNodeArpTable(node); curr != nil; curr = curr.Next {
		t.AppendRow(table.Row{
			tools.ConvertAddrToStr(curr.IpAddr.Addr[:]),
			tools.ConvertAddrToStr(curr.MacAddr.Addr[:]),
			curr.Name})
	}
	t.Render()
}

func dumpMacTable(node *network.Node) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"MAC", "Interface"})
	for curr := network.GetNodeMacTable(node); curr != nil; curr = curr.Next {
		t.AppendRow(table.Row{
			tools.ConvertAddrToStr(curr.MacAddr.Addr[:]),
			curr.Name})
	}
	t.Render()
}
