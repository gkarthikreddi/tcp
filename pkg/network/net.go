package network

import (
	"fmt"
	"github.com/gkarthikreddi/tcp/tools"
	"net"
)

type Ip struct {
	Addr [4]byte
}

type Mac struct {
	Addr [6]byte
}

type L2Mode string

const (
	ACCESS  L2Mode = "access"
	TRUNK   L2Mode = "trunk"
	UNKNOWN L2Mode = "unknown"
)

type nodeProp struct {
	// L3 properties
	isLbAddr bool
	lbAddr   Ip

	// L2 properties
	arpTable *ArpEntry
	macTable *MacEntry

	port   int
	socket *net.UDPAddr
}

type intfProp struct {
	macAddr  Mac
	isIpAddr bool
	ipAddr   Ip
	mask     int

	// L2 properties
	l2Mode L2Mode
}

type ArpEntry struct {
	IpAddr  *Ip
	MacAddr *Mac
	Name    string
	Next    *ArpEntry
	Prev    *ArpEntry
}

type MacEntry struct {
	MacAddr *Mac
	Name    string
	Next    *MacEntry
	Prev    *MacEntry
}

// Encapsulation
func IsIntfIp(intf *Interface) bool {
	return intf.prop.isIpAddr
}

func GetIntfIp(intf *Interface) *Ip {
	return &intf.prop.ipAddr
}

func GetIntfMac(intf *Interface) *Mac {
	return &intf.prop.macAddr
}

func GetIntfL2Mode(intf *Interface) L2Mode {
	return intf.prop.l2Mode
}

func GetNodeIp(node *Node) *Ip {
	return &node.prop.lbAddr
}

func GetNodePort(node *Node) int {
	return node.prop.port
}

func GetNodeSocket(node *Node) *net.UDPAddr {
	return node.prop.socket
}

func GetNodeArpTable(node *Node) *ArpEntry {
	return node.prop.arpTable
}

func GetNodeMacTable(node *Node) *MacEntry {
	return node.prop.macTable
}

func AssignNodePort(node *Node, num int) {
	node.prop.port = num
}

func AssignNodeSocket(node *Node, socket *net.UDPAddr) {
	node.prop.socket = socket
}

func AssignNodeArpTable(node *Node, arpEntry *ArpEntry) {
	node.prop.arpTable = arpEntry
}

func AssignNodeMacTable(node *Node, macEntry *MacEntry) {
	node.prop.macTable = macEntry
}

// --------------------

func NodeSetLbAddr(node *Node, addr string) bool {
	if addr == "" {
		return false
	}

	node.prop.isLbAddr = true
	GetNodeIp(node).Addr = tools.ConvertStrToIp(addr)

	return true
}

func NodeSetIntfIpAddr(node *Node, name, addr string, mask int) bool {
	intf, err := GetIntfByIntfName(node, name)
	if err != nil {
		return false
	}

	if err := intfAssignMacAddr(intf); err != nil {
		return false
	}
	intf.prop.isIpAddr = true
	GetIntfIp(intf).Addr = tools.ConvertStrToIp(addr)
	intf.prop.mask = mask

	return true
}

func NodeUnsetIntfIpAddr(node *Node, name string) bool {
	intf, err := GetIntfByIntfName(node, name)
	if err != nil {
		return false
	}

	intf.prop = intfProp{}

	return true
}

func NodeGetMatchingSubnet(node *Node, ip *Ip) (*Interface, error) {
	for i := 0; i < MAX_INTF_PER_NODE; i++ {
		intf := node.Intf[i]
		if intf == nil {
			break
		}
		if intf.prop.isIpAddr == false {
			continue
		}

		addr := GetIntfIp(intf)
		mask := intf.prop.mask

		n1 := applyMask(addr, mask)
		n2 := applyMask(ip, mask)

		if n1 == n2 {
			return intf, nil
		}
	}
	return nil, fmt.Errorf("No matching subnet for the given node")
}

func NodeSetIntfL2Mode(node *Node, name string, mode L2Mode) bool {
	intf, err := GetIntfByIntfName(node, name)
	if err != nil {
		return false
	}

    intf.prop.l2Mode = mode
	return true
}

func intfAssignMacAddr(intf *Interface) error {
	if mac, err := tools.RandomMacAddr(); err == nil {
		for i, val := range mac {
			intf.prop.macAddr.Addr[i] = val
		}
		return nil
	} else {
		return err
	}
}

func applyMask(ip *Ip, mask int) [4]byte {
	subnet := tools.GetSubnetFromMask(mask)
	var ans [4]byte
	for i := 0; i < 4; i++ {
		ans[i] = subnet[i] & ip.Addr[i]
	}

	return ans
}
