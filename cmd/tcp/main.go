package main

import (
	"github.com/gkarthikreddi/tcp/pkg/cli"
	"github.com/gkarthikreddi/tcp/tools/cmdparser"
)

func main() {
	cli.InitNwCli()
	cmdparser.CommandParser()
}
