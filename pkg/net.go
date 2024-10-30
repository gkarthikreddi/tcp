package pkg

import (
	"fmt"
	"net"
	"github.com/gkarthikreddi/tcp/tools"
)

type Ip struct {
	Addr string
}

type Mac struct {
	Addr string
}

type NodeProp struct {
	IsLbAddr bool
	LbAddr   Ip
	Port     int
	Socket   *net.UDPAddr
}

type IntfProp struct {
	MacAddr  Mac
	IsIpAddr bool
	IpAddr   Ip
	Mask     int
}

// some macros
func intfMac(intf *Interface) *string {
	return &intf.Prop.MacAddr.Addr
}

func intfIp(intf *Interface) *string {
	return &intf.Prop.IpAddr.Addr
}

func isIntfIp(intf *Interface) *bool {
	return &intf.Prop.IsIpAddr
}

func nodeIp(node *Node) *string {
	return &node.Prop.LbAddr.Addr
}

func intfAssignMacAddr(intf *Interface) {
	*intfMac(intf) = intf.Att_node.Name + intf.Name
}

func NodeSetLbAddr(node *Node, addr string) bool {
	if addr == "" {
		return false
	}

	node.Prop.IsLbAddr = true
	*nodeIp(node) = addr

	return true
}

func NodeSetIntfIpAddr(node *Node, name, addr string, mask int) bool {
	intf, err := getIntfByIntfName(node, name)
	if err != nil {
		return false
	}

	intfAssignMacAddr(intf)
	intf.Prop.IsIpAddr = true
	*intfIp(intf) = addr
	intf.Prop.Mask = mask

	return true
}

func NodeUnsetIntfIpAddr(node *Node, name string) bool {
	intf, err := getIntfByIntfName(node, name)
	if err != nil {
		return false
	}

	intf.Prop = IntfProp{}

	return true
}

func NodeGetMatchingSubnet(node *Node, ip string) (*Interface, error) {
	for i := 0; i < MAX_INTF_PER_NODE; i++ {
		intf := node.Intf[i]
		if intf == nil {
			break
		}
		if intf.Prop.IsIpAddr == false {
			continue
		}

		addr := *intfIp(intf)
		mask := intf.Prop.Mask

		n1 := tools.ApplyMask(addr, mask)
		n2 := tools.ApplyMask(ip, mask)

		if n1 == n2 {
			return intf, nil
		}
	}
	return nil, fmt.Errorf("No matching subnet for the given node")
}
