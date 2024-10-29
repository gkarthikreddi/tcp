package main

import (
	"github.com/gkarthikreddi/tcp/pkg"
	"github.com/gkarthikreddi/tcp/tools/cmdparser"
)

func main() {
    pkg.InitNwCli()
    cmdparser.CommandParser()
}
