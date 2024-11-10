package stack

import (
	"fmt"

	"github.com/gkarthikreddi/tcp/pkg/network"
)

// for now will only implement ping functionality
func Ping(node *network.Node, dstIPAddr [4]byte) {
	ip := &network.Ip{Addr: dstIPAddr}
    if err := demotePktToLayer3(node, ip, ICMP_PRO); err != nil {
        fmt.Println(err)
    }
}
