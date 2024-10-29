package pkg

import (
	"fmt"
)

const MAX_INTF_SIZE int = 10
const MAX_INTF_PER_NODE int = 10

type Link struct {
	Intf1 Interface
	Intf2 Interface
	Cost  uint
}

type Interface struct {
	Name     string
	Att_node *Node
	Wire     *Link
	Prop     IntfProp
}

type Node struct {
	Name string
	Intf [MAX_INTF_SIZE]*Interface
	Prev *Node
	Next *Node
	Prop NodeProp
}

type Graph struct {
	Name string
	List *Node
}

// Helper functions
func getNbrNode(intf *Interface) (*Node, error) {
	if intf.Att_node == nil || intf.Wire == nil {
		return nil, fmt.Errorf("Either att_node or wire is not there")
	}

	if intf.Wire.Intf1 == *intf {
		return intf.Wire.Intf2.Att_node, nil
	}
	return intf.Wire.Intf1.Att_node, nil
}

func getNodeIntfAvailableSlot(node *Node) (int, error) {
	for i := 0; i < MAX_INTF_PER_NODE; i++ {
		if node.Intf[i] == nil {
			return i, nil
		}
	}
	return -1, fmt.Errorf("No available slots in the node")
}

func getNodeByIntfName(node *Node, name string) (*Interface, error) {
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

func getNodeByNodeName(graph *Graph, name string) (*Node, error) {
	for node := graph.List; node != nil; node = node.Next {
		if node.Name == name {
			return node, nil
		}
	}
	return nil, fmt.Errorf("No interface with the given name: %s", name)
}

// Main functions
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

	wire := Link{Intf1: intf1, Intf2: intf2}

	// Setting back pointers
	wire.Intf1.Wire = &wire
	wire.Intf2.Wire = &wire

	wire.Intf1.Att_node = node1
	wire.Intf2.Att_node = node2
	wire.Cost = cost

	if i, err := getNodeIntfAvailableSlot(node1); err == nil {
		node1.Intf[i] = &wire.Intf1
	} else {
		return fmt.Errorf("Node available slots in node1")
	}
	if i, err := getNodeIntfAvailableSlot(node2); err == nil {
		node2.Intf[i] = &wire.Intf2
	} else {
		return fmt.Errorf("Node available slots in node2")
	}

	return nil
}
