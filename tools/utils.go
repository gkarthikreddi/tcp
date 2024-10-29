package tools

import (
	"strings"
    "strconv"
)

func ApplyMask(ip string, mask int) string {
	addr := strings.Split(ip, ".")
	subnet := splitSubnet(strings.Repeat("1", int(mask)) + strings.Repeat("0", 32-int(mask)))

    var network []string
    for i := 0; i < len(addr); i++ {
        a, _ := strconv.Atoi(addr[i])
        b := subnet[i]
        network = append(network, strconv.Itoa(a & b))
    }

    return strings.Join(network, ".")
}

func splitSubnet(subnet string) []int {
    var ans []int
    for i := 0; i < len(subnet); i += 8 {
        end := i+8
        num, _ := strconv.ParseInt(subnet[i:end], 2, 8)
        ans = append(ans, int(num))
    }

    return ans
}
