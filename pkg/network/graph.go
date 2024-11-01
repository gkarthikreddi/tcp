package network

import (
	"fmt"
)

const MAX_INTF_SIZE int = 10
const MAX_INTF_PER_NODE int = 10

type link struct {
	intf1 Interface
	intf2 Interface
	cost  uint
}

type Interface struct {
	Name     string
	Att_node *Node
	conn     *link
	prop     intfProp
}

type Node struct {
	Name string
	Intf [MAX_INTF_SIZE]*Interface
	Prev *Node
	Next *Node
	prop nodeProp
}

type Graph struct {
	Name string
	List *Node
}

func getNodeIntfAvailableSlot(node *Node) (int, error) {
	for i := 0; i < MAX_INTF_PER_NODE; i++ {
		if node.Intf[i] == nil {
			return i, nil
		}
	}
	return -1, fmt.Errorf("No available slots in the node")
}

func getIntfByIntfName(node *Node, name string) (*Interface, error) {
	for i := 0; i < MAX_INTF_PER_NODE; i++ {
		if node.Intf[i] == nil {
			break
		}
		if node.Intf[i].Name == name {
			return node.Intf[i], nil
		}
	}
	return nil, fmt.Errorf("No interface with the given name: %s", name)
}

func GetNbrNode(intf *Interface) (*Node, error) {
	if intf.Att_node == nil || intf.conn == nil {
		return nil, fmt.Errorf("Either att_node or wire is not there")
	}

	if intf.conn.intf1 == *intf {
		return intf.conn.intf2.Att_node, nil
	}
	return intf.conn.intf1.Att_node, nil
}

func GetNodeByNodeName(graph *Graph, name string) (*Node, error) {
	for node := graph.List; node != nil; node = node.Next {
		if node.Name == name {
			return node, nil
		}
	}
	return nil, fmt.Errorf("No interface with the given name: %s", name)
}

func CreateNewGraph(name string) *Graph {
	graph := Graph{Name: name, List: nil}
	return &graph
}

func CreateGraphNode(graph *Graph, name string) *Node {
	node := Node{Name: name}

	if graph.List == nil {
		graph.List = &node
		return &node
	}
	for curr := graph.List; curr != nil; curr = curr.Next {
		if curr.Next == nil {
			node.Prev = curr
			curr.Next = &node
			break
		}
	}
	return &node
}

func InsertLinkBetweenNodes(node1, node2 *Node, fromIntfNode, toIntfNode string, cost uint) error {
	intf1 := Interface{Name: fromIntfNode}
	intf2 := Interface{Name: toIntfNode}

	wire := link{intf1: intf1, intf2: intf2}

	// Setting back pointers
	wire.intf1.conn = &wire
	wire.intf2.conn = &wire

	wire.intf1.Att_node = node1
	wire.intf2.Att_node = node2
	wire.cost = cost

	if i, err := getNodeIntfAvailableSlot(node1); err == nil {
		node1.Intf[i] = &wire.intf1
	} else {
		return fmt.Errorf("Node available slots in node1")
	}
	if i, err := getNodeIntfAvailableSlot(node2); err == nil {
		node2.Intf[i] = &wire.intf2
	} else {
		return fmt.Errorf("Node available slots in node2")
	}

	return nil
}
