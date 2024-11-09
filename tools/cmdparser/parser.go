package cmdparser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var buff *SerBuff

func FindMatchingParam(param *Param, name string) (*Param, error) {
	idx := -1
	for i := 0; i < MAX_CHILDREN; i++ {

		if param.children[i] == nil {
			break
		}
		child := param.children[i].kind
		if child.leaf != nil {
			idx = i
			continue
		}
		if child.cmd.name == name {
			return param.children[i], nil
		}
	}

	if idx > -1 {
		return param.children[idx], nil
	}
	return nil, fmt.Errorf("Can't find child with the given name: %s", name)
}

func buildTlvBuffer(param *Param, value string) {
	tlv := Tlv{}
	leaf := GetLeaf(param)

	tlv.LeafType = leaf.leafType
	tlv.Id = leaf.id
	tlv.Value = value
	leaf.value = value

	collectTlv(&tlv)
}

func collectTlv(tlv *Tlv) {
	if buff == nil {
		buff = &SerBuff{Data: tlv}
	} else {
		for curr := buff; curr != nil; curr = curr.Next {
			if curr.Next == nil {
				curr.Next = &SerBuff{Data: tlv}
				break
			}
		}
	}
}

func ExtractCmdCode(buff *SerBuff) int {
	ans, _ := strconv.Atoi(buff.Data.Value)
	return ans
}

func parser(tokens []string) error {
	parent := GetRootHook()

	for i := 0; i < len(tokens); i++ {
		if param, err := FindMatchingParam(parent, tokens[i]); err == nil {
			if param.kind.leaf != nil {
				leaf := GetLeaf(param)
				if leaf.fn != nil {
					if !leaf.fn(tokens[i]) {
						return fmt.Errorf("Given <param-value> '%s' is not suitalbe for %s", tokens[i], leaf.id)
					}
				}
				buildTlvBuffer(param, tokens[i])
			}
			parent = param
		} else {
			return fmt.Errorf("Incorrect usage of command")
		}
	}

	if parent != nil {
		tlv := Tlv{Value: strconv.Itoa(parent.code)}
		buff = &SerBuff{Data: &tlv, Next: buff}
	}

	if parent.fn != nil {
		parent.fn(parent, buff)
	} else {
		return fmt.Errorf("Incomplete Command")
	}

	return nil
}

func CommandParser() {
	for true {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("router> ")
		str, _ := reader.ReadString('\n')
		tokens := strings.Split(strings.TrimSpace(str), " ")

		if err := parser(tokens); err != nil {
			fmt.Println(err)
		}
		buff = nil
		time.Sleep(time.Millisecond * 50)
	}
}
