package cli

import (
	"strconv"
	"strings"
)

func validNodeName(str string) bool {
	for curr := graph.List; curr != nil; curr = curr.Next {
		if curr.Name == str {
			return true
		}
	}
	return false
}

func validIPAddr(str string) bool {
	addr := strings.Split(str, ".")
	if len(addr) != 4 {
		return false
	}
	for _, val := range addr {
		if i, err := strconv.Atoi(val); err == nil {
			if i > 256 || i < 0 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}
