package stack

import (
	"fmt"
	"net"
	"sync"

	"github.com/gkarthikreddi/tcp/pkg/network"
	"github.com/gkarthikreddi/tcp/tools"
)

const (
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
)

var port = 40000

type packet struct {
	Intf       string
	EtherFrame ethernetHeader
}

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

func sendPkt(etherFrame *ethernetHeader, intf *network.Interface) error {
	dstNode, err := network.GetNbrNode(intf)
	if err != nil {
		return err
	}

	dstIntf := network.GetNbrIntf(intf)
	pkt := packet{Intf: dstIntf, EtherFrame: *etherFrame}
	msg, err := tools.StructToByte(pkt)
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
	return nil
}

func receivePkt(node *network.Node, data []byte) error {
	if pkt, err := tools.ByteToStruct(data, packet{}); err == nil {
		if intf, err := network.GetIntfByIntfName(node, pkt.Intf); err == nil {
			layer2FrameRecieve(node, intf, &pkt.EtherFrame)
			return nil
		}
	}
	return fmt.Errorf("Error while trasferring the received packet to layer2 of node: %s", node.Name)
}
