package main

import (
	"fmt"

	"github.com/gkarthikreddi/tcp/tools/cmdparser"
)

func helper(param *cmdparser.Param, buff *cmdparser.SerBuff) int {
	code := cmdparser.ExtractCmdCode(buff)

	switch code {
	case 1:
		fmt.Println("hello")
	case 0:
		fmt.Println("world")
	default:
		fmt.Println("JOhn")

	}
	return -1
}

func main() {

	cmdparser.InitLibcli()
    show := cmdparser.GetShowHook()
	{
		var topo cmdparser.Param
		cmdparser.InitParam(&topo,
			cmdparser.CMD,
			"topology",
			helper,
			nil,
			cmdparser.INVALID,
			"",
			"Dump complete network")
        cmdparser.LibcliRegisterParam(show, &topo)
        cmdparser.SetParamCmdCode(&topo, 1)
	}

	cmdparser.CommandParser()
}
