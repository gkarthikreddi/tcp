package pkg

import (
	"fmt"
	"net"
	"sync"
)

var port = 40000

func initUdpSocket(node *Node) error {
	port++
	node.Prop.Port = port

	if socket, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port)); err == nil {
		node.Prop.Socket = socket
		fmt.Printf("port: %d is up on node: %s\n", port, node.Name)
		return nil
	}
	return fmt.Errorf("Can't bind upd port to node: %s", node.Name)
}

func InitNetworkListening(graph *Graph) {
	nodes := []*Node{}
	for curr := graph.List; curr != nil; curr = curr.Next {
		nodes = append(nodes, curr)
	}

	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Add(1)
		go startListening(node, &wg)
	}
    wg.Wait()
}

func startListening(node *Node, wg *sync.WaitGroup) error {
	conn, err := net.ListenUDP("udp", node.Prop.Socket)
	if err != nil {
		return fmt.Errorf("Error while establishing connection on node: %s with port: %d", node.Name, node.Prop.Port)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	fmt.Printf("Listing on node: %s on port: %d\n", node.Name, node.Prop.Port)
	wg.Done()
	for {
		if n, _, err := conn.ReadFromUDP(buffer); err == nil {
			receivePkt(node, string(buffer[:n]))
		} else {
			fmt.Printf("Error while receiving data on node: %s on port: %d", node.Name, node.Prop.Port)
		}
	}
}

func sendPkt(data string, intf *Interface) error {
	dstNode, err := getNbrNode(intf)
	if err != nil {
		return err
	}

	var dstIntf Interface
	if intf.Wire.Intf1 == *intf {
		dstIntf = intf.Wire.Intf2
	} else {
		dstIntf = intf.Wire.Intf1
	}

	data = fmt.Sprintf("%s:%s", dstIntf.Name, data)
	msg := []byte(data)
	dstPort := dstNode.Prop.Port

	dstaddr := net.UDPAddr{Port: dstPort, IP: net.ParseIP("127.0.0.1")}
	if conn, err := net.DialUDP("udp", nil, &dstaddr); err == nil {
		conn.Write(msg)
		conn.Close()
	} else {
		return fmt.Errorf("Can't estrablish connection with DestinationNode: %s, Port: %d", dstNode.Name, dstPort)
	}
	return nil
}

func receivePkt(node *Node, msg string) error {
	fmt.Println(msg)
	return nil
}

// Assignment
func sendPktFlood(node * Node, excludeIntf *Interface, msg string) {
    for i := 0; i < MAX_INTF_PER_NODE; i++ {
        if node.Intf[i] != nil && node.Intf[i] != excludeIntf{
            sendPkt(msg, node.Intf[i])
        }
    }
}
