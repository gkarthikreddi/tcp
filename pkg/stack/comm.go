package stack

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gkarthikreddi/tcp/pkg/network"
	"github.com/gkarthikreddi/tcp/tools"
)

var port = 40000

func initUdpSocket(node *network.Node) error {
	port++
	network.AssignNodePort(node, port)

	if socket, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port)); err == nil {
		network.AssignNodeSocket(node, socket)
		return nil
	}
	return fmt.Errorf("Can't bind upd port to node: %s", node.Name)
}

func InitNetworkListening(graph *network.Graph) {
	nodes := []*network.Node{}
	for curr := graph.List; curr != nil; curr = curr.Next {
		initUdpSocket(curr)
		nodes = append(nodes, curr)
	}

	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Add(1)
		go startListening(node, &wg)
	}
	wg.Wait()
}

func startListening(node *network.Node, wg *sync.WaitGroup) error {
	conn, err := net.ListenUDP("udp", network.GetNodeSocket(node))
	if err != nil {
		return fmt.Errorf("Error while establishing connection on node: %s with port: %d", node.Name, network.GetNodePort(node))
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	wg.Done()
	for {
		if n, _, err := conn.ReadFromUDP(buffer); err == nil {
			if err = receivePkt(node, buffer[:n]); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("Error while receiving data on node: %s on port: %d", node.Name, network.GetNodePort(node))
		}
	}
}

func sendPkt(etherFame *ethernetHeader, intf *network.Interface) error {
	dstNode, err := network.GetNbrNode(intf)
	if err != nil {
		return err
	}

	msg, err := tools.StructToByte(etherFame)
	if err != nil {
		return err
	}
	dstPort := network.GetNodePort(dstNode)
	dstaddr := net.UDPAddr{Port: dstPort, IP: net.ParseIP("127.0.0.1")}
	if conn, err := net.DialUDP("udp", nil, &dstaddr); err == nil {
		conn.Write(msg)
		conn.Close()
	} else {
		return fmt.Errorf("Can't estrablish connection with DestinationNode: %s, Port: %d", dstNode.Name, dstPort)
	}
	time.Sleep(time.Millisecond * 5)
	return nil
}

func receivePkt(node *network.Node, data []byte) error {
	if etherFrame, err := tools.ByteToStruct(data, ethernetHeader{}); err == nil {
		if arpFrame, e := tools.ByteToStruct(etherFrame.Payload[:], arpHeader{}); e == nil {
			ip := &network.Ip{Addr: arpFrame.DstProtocolAddr}
			if intf, err := network.NodeGetMatchingSubnet(node, ip); err == nil {
				layer2FrameRecieve(node, intf, etherFrame)
				return nil
			}
		}
	}
	return fmt.Errorf("Error while trasferring the received packet to layer2 of node: %s", node.Name)
}

// Assignment
func sendPktFlood(node *network.Node, excludeIntf *network.Interface, msg string) {
	for i := 0; i < network.MAX_INTF_PER_NODE; i++ {
		if node.Intf[i] != nil && node.Intf[i] != excludeIntf {
			// sendPkt(msg, node.Intf[i])
		}
	}
}
