package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gkarthikreddi/tcp/pkg/network"
	"github.com/gkarthikreddi/tcp/pkg/stack"
	"github.com/gkarthikreddi/tcp/tools"
	"github.com/jedib0t/go-pretty/v6/table"
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
	stack.InitRoutingTable(graph)
	return graph
}

func buildL2SwitchGraph() *network.Graph {
	graph := network.CreateNewGraph("Simpel L2 switch demo graph")
	h1 := network.CreateGraphNode(graph, "H1")
	h2 := network.CreateGraphNode(graph, "H2")
	h3 := network.CreateGraphNode(graph, "H3")
	h4 := network.CreateGraphNode(graph, "H4")
	l2sw := network.CreateGraphNode(graph, "L2SW")

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

	stack.InitNetworkListening(graph)
	stack.InitRoutingTable(graph)
	return graph
}

func buildDualSwitchGraph() *network.Graph {
	graph := network.CreateNewGraph("Dual Switch Topo")
	H1 := network.CreateGraphNode(graph, "H1")
	network.NodeSetLbAddr(H1, "122.1.1.1")
	H2 := network.CreateGraphNode(graph, "H2")
	network.NodeSetLbAddr(H2, "122.1.1.2")
	H3 := network.CreateGraphNode(graph, "H3")
	network.NodeSetLbAddr(H3, "122.1.1.3")
	H4 := network.CreateGraphNode(graph, "H4")
	network.NodeSetLbAddr(H4, "122.1.1.4")
	H5 := network.CreateGraphNode(graph, "H5")
	network.NodeSetLbAddr(H5, "122.1.1.5")
	H6 := network.CreateGraphNode(graph, "H6")
	network.NodeSetLbAddr(H6, "122.1.1.6")

	L2SW1 := network.CreateGraphNode(graph, "L2SW1")
	L2SW2 := network.CreateGraphNode(graph, "L2SW2")

	network.InsertLinkBetweenNodes(H1, L2SW1, "eth0/1", "eth0/2", 1)
	network.InsertLinkBetweenNodes(H2, L2SW1, "eth0/3", "eth0/7", 1)
	network.InsertLinkBetweenNodes(H3, L2SW1, "eth0/4", "eth0/6", 1)
	network.InsertLinkBetweenNodes(L2SW1, L2SW2, "eth0/5", "eth0/7", 1)
	network.InsertLinkBetweenNodes(H5, L2SW2, "eth0/8", "eth0/9", 1)
	network.InsertLinkBetweenNodes(H4, L2SW2, "eth0/11", "eth0/12", 1)
	network.InsertLinkBetweenNodes(H6, L2SW2, "eth0/11", "eth0/10", 1)

	network.NodeSetIntfIpAddr(H1, "eth0/1", "10.1.1.1", 24)
	network.NodeSetIntfIpAddr(H2, "eth0/3", "10.1.1.2", 24)
	network.NodeSetIntfIpAddr(H3, "eth0/4", "10.1.1.3", 24)
	network.NodeSetIntfIpAddr(H4, "eth0/11", "10.1.1.4", 24)
	network.NodeSetIntfIpAddr(H5, "eth0/8", "10.1.1.5", 24)
	network.NodeSetIntfIpAddr(H6, "eth0/11", "10.1.1.6", 24)

	network.NodeSetIntfL2Mode(L2SW1, "eth0/2", network.ACCESS)
	network.NodeSetIntfVlanMembership(L2SW1, "eth0/2", 10)
	network.NodeSetIntfL2Mode(L2SW1, "eth0/7", network.ACCESS)
	network.NodeSetIntfVlanMembership(L2SW1, "eth0/7", 10)
	network.NodeSetIntfL2Mode(L2SW1, "eth0/5", network.TRUNK)
	network.NodeSetIntfVlanMembership(L2SW1, "eth0/5", 10)
	network.NodeSetIntfVlanMembership(L2SW1, "eth0/5", 11)
	network.NodeSetIntfL2Mode(L2SW1, "eth0/6", network.ACCESS)
	network.NodeSetIntfVlanMembership(L2SW1, "eth0/6", 11)

	network.NodeSetIntfL2Mode(L2SW2, "eth0/7", network.TRUNK)
	network.NodeSetIntfVlanMembership(L2SW2, "eth0/7", 10)
	network.NodeSetIntfVlanMembership(L2SW2, "eth0/7", 11)
	network.NodeSetIntfL2Mode(L2SW2, "eth0/9", network.ACCESS)
	network.NodeSetIntfVlanMembership(L2SW2, "eth0/9", 10)
	network.NodeSetIntfL2Mode(L2SW2, "eth0/10", network.ACCESS)
	network.NodeSetIntfVlanMembership(L2SW2, "eth0/10", 10)
	network.NodeSetIntfL2Mode(L2SW2, "eth0/12", network.ACCESS)
	network.NodeSetIntfVlanMembership(L2SW2, "eth0/12", 11)

	stack.InitNetworkListening(graph)
	stack.InitRoutingTable(graph)

	return graph
}
func buildLinear3NodeTopo() *network.Graph {
	graph := network.CreateNewGraph("3 node linear topo")
	R1 := network.CreateGraphNode(graph, "R1")
	R2 := network.CreateGraphNode(graph, "R2")
	R3 := network.CreateGraphNode(graph, "R3")

	network.InsertLinkBetweenNodes(R1, R2, "eth0/1", "eth0/2", 1)
	network.InsertLinkBetweenNodes(R2, R3, "eth0/3", "eth0/4", 1)

	network.NodeSetLbAddr(R1, "122.1.1.1")
	network.NodeSetLbAddr(R2, "122.1.1.2")
	network.NodeSetLbAddr(R3, "122.1.1.3")

	network.NodeSetIntfIpAddr(R1, "eth0/1", "10.1.1.1", 24)
	network.NodeSetIntfIpAddr(R2, "eth0/2", "10.1.1.2", 24)
	network.NodeSetIntfIpAddr(R2, "eth0/3", "11.1.1.2", 24)
	network.NodeSetIntfIpAddr(R3, "eth0/4", "11.1.1.1", 24)

	stack.InitNetworkListening(graph)
	stack.InitRoutingTable(graph)

	return graph
}

func buildSquareTopo() *network.Graph {
	topo := network.CreateNewGraph("Square Topo")
	R1 := network.CreateGraphNode(topo, "R1")
	R2 := network.CreateGraphNode(topo, "R2")
	R3 := network.CreateGraphNode(topo, "R3")
	R4 := network.CreateGraphNode(topo, "R4")

	network.InsertLinkBetweenNodes(R1, R2, "eth0/0", "eth0/1", 1)
	network.InsertLinkBetweenNodes(R2, R3, "eth0/2", "eth0/3", 1)
	network.InsertLinkBetweenNodes(R3, R4, "eth0/4", "eth0/5", 1)
	network.InsertLinkBetweenNodes(R4, R1, "eth0/6", "eth0/7", 1)

	network.NodeSetLbAddr(R1, "122.1.1.1")
	network.NodeSetIntfIpAddr(R1, "eth0/0", "10.1.1.1", 24)
	network.NodeSetIntfIpAddr(R1, "eth0/7", "40.1.1.2", 24)

	network.NodeSetLbAddr(R2, "122.1.1.2")
	network.NodeSetIntfIpAddr(R2, "eth0/1", "10.1.1.2", 24)
	network.NodeSetIntfIpAddr(R2, "eth0/2", "20.1.1.1", 24)

	network.NodeSetLbAddr(R3, "122.1.1.3")
	network.NodeSetIntfIpAddr(R3, "eth0/3", "20.1.1.2", 24)
	network.NodeSetIntfIpAddr(R3, "eth0/4", "30.1.1.1", 24)

	network.NodeSetLbAddr(R4, "122.1.1.4")
	network.NodeSetIntfIpAddr(R4, "eth0/5", "30.1.1.2", 24)
	network.NodeSetIntfIpAddr(R4, "eth0/6", "40.1.1.1", 24)

	stack.InitNetworkListening(topo)
	stack.InitRoutingTable(topo)

	return topo
}

var graph = buildSquareTopo()

func dumpGraph(graph *network.Graph) {
	fmt.Println("Name: " + Cyan + graph.Name + Reset)
	if graph.Name == "Topology" {
		fmt.Println(`
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
        `)
	} else if graph.Name == "Dual Switch Topo" {
		fmt.Println(`
	                                    +---------+                               +----------+
	                                    |         |                               |          |
	                                    |  H2     |                               |  H5      |
	                                    |122.1.1.2|                               |122.1.1.5 |
	                                    +---+-----+                               +-----+----+
	                                        |10.1.1.2/24                                +10.1.1.5/24
	                                        |eth0/3                                     |eth0/8
	                                        |                                           |
	                                        |eth0/7,AC,V10                              |eth0/9,AC,V10
	                                  +-----+----+                                +-----+---+
	                                  |          |                                |         |
	   +------+---+                   |          |                                |         |                         +--------+
	   |  H1      |10.1.1.1/24        |   L2SW1  |eth0/5                    eth0/7| L2SW2   |eth0/10           eth0/11|  H6    |
	   |122.1.1.1 +-------------------|          |+-------------------------------|         +-------------+----------+122.1.1.6|
	   +------+---+ eth0/1      eth0/2|          |TR,V10,V11            TR,V10,V11|         |AC,V10        10.1.1.6/24|        |
	                            AC,V10|          |                                |         |                         +-+------+
	                                  +-----+----+                                +----+----+
	                                        |eth0/6                                    |eth0/12
	                                        |AC,V11                                    |AC,V11
	                                        |                                          |
	                                        |                                          |
	                                        |                                          |
	                                        |                                          |eth0/11
	                                        |eth0/4                                    |10.1.1.4/24
	                                        |10.1.1.3/24                             +--+-----+
	                                   +----+---+|                                   | H4     |
	                                   |  H3     |                                   |        |
	                                   |122.1.1.3|                                   |122.1.1.4|
	                                   +--------+|                                   +--------+ `)
	} else if graph.Name == "" {
		fmt.Println(`
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
	                                  +------------+ `)
	} else if graph.Name == "3 node linear topo" {
		fmt.Println(`
                                        +---------+                                  +----------+
+--------+                              |         |                                  |R3        |
|R1      |eth0/1                  eth0/2|R2       |eth0/3                      eth0/4|122.1.1.3 |
|122.1.1.1+-----------------------------+122.1.1.2|+----------------------------------+         |
|        |10.1.1.1/24        10.1.1.2/24|         |11.1.1.2/24            11.1.1.1/24|          |
+--------+                              +-------+-|                                  +----------+ `)
	} else if graph.Name == "Square Topo" {
		fmt.Println(`
  +-----------+                      +--------+                            +--------+
  |           |eth0/0     10.1.1.2/24|        | eth0/2               eth0/3|        |
  | R1        +----------------------|  R2    +----------------------------+   R3   |
  |122.1.1.1  |10.1.1.1/24     eth0/1|122.1.1.2|  20.1.1.1/24   20.1.1.2/24| 122.1.1.3|
  +---+--+----+                      |        |                            +-       +
         |eth0/7                     +--------+                            +----+---+
         | 40.1.1.2/24                                                          | eth0/4   
         |                                                                      |30.1.1.1/24
         |                                                                      |
         |                                                                      |
         |                                                                      |
         |                                                                      |
         |                                                                      |
         |                          +-----------+                               |
         |                          |           |                               |
         |                  eth0/6  |  R4       |                               |
         +--------------------------+ 122.1.1.4 |                               |
                         40.1.1.1/24|           +-------------------------------+
                                    |           |eth0/5
                                    +-----------+30.1.1.2/24 `)
	}
	for curr := graph.List; curr != nil; curr = curr.Next {
		dumpNode(curr)
	}
}

func dumpNode(node *network.Node) {
	fmt.Println("Node name: " + node.Name + "\nLb addr: " + Yellow + tools.ConvertAddrToStr(network.GetNodeIp(node).Addr[:]) + Reset + ", UDP Port: " + strconv.Itoa(network.GetNodePort(node)))
	for i := 0; node.Intf[i] != nil; i++ {
		dumpInterface(node.Intf[i])
	}
}

func dumpInterface(intf *network.Interface) {
	fmt.Println("\tInterface name: " + Cyan + intf.Name + Reset)
	nbrNode, _ := network.GetNbrNode(intf)
	fmt.Println("\t\tLocalNode: " + Cyan + intf.Att_node.Name + Reset + ", Nbr Node: " + Cyan + nbrNode.Name + Reset)

	if network.IsIntfIp(intf) {
		fmt.Println("\t\tIp addr: " + Yellow + tools.ConvertAddrToStr(network.GetIntfIp(intf).Addr[:]) + Reset + " Mac addr: " + Yellow + tools.ConvertAddrToStr(network.GetIntfMac(intf).Addr[:]) + Reset)
	} else {
		fmt.Printf("\t\tL2 Mode: %v\t Vlan Membership: ", network.GetIntfL2Mode(intf))
		for _, val := range network.GetIntfVlanMembership(intf) {
			if val == 0 {
				break
			} else {
				fmt.Printf("%d ", val)
			}
		}
		fmt.Println()
	}
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

func dumpRoutingTable(node *network.Node) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Dst IpAddr", "Mask", "Direct", "Gateway IpAddr", "Outgoing Intf"})
	for curr := network.GetNodeRoutingTable(node); curr != nil; curr = curr.Next {
		addr := "NA"
		if curr.GatewayIp != nil {
			addr = tools.ConvertAddrToStr(curr.GatewayIp.Addr[:])
		}
		t.AppendRow(table.Row{
			tools.ConvertAddrToStr(curr.DstIpAddr.Addr[:]),
			curr.DstIpAddr.Mask,
			curr.IsDirect,
			addr,
			curr.OutIntf,
		})
	}
	t.Render()
}
