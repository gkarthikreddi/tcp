package cli

import "github.com/gkarthikreddi/tcp/tools/cmdparser"

const (
	SHOW_TOPO   = 1
	ARP_HANDLER = 2
	ARP_TABLE   = 3
)

func InitNwCli() {
	cmdparser.InitLibcli()

	show := cmdparser.GetShowHook()
	run := cmdparser.GetRunHook()

	{
		var topo cmdparser.Param
		cmdparser.InitParam(&topo, // param
			cmdparser.CMD,                  // type of param
			"topology",                     // name of param, nil for leaf param
			showTopology,                   // callback handler
			nil,                            // validationn handler
			cmdparser.INVALID,              // leaftype
			"",                             // id of leaf, nil for cmd param
			"Dump entire network topology") // help string
		cmdparser.LibcliRegisterParam(show, &topo)
		cmdparser.SetParamCmdCode(&topo, SHOW_TOPO)
	}
	{
		var node cmdparser.Param
		cmdparser.InitParam(&node,
			cmdparser.CMD,
			"node",
			nil,
			nil,
			cmdparser.INVALID,
			"",
			"Given a node name and operation it performs that operation")
		cmdparser.LibcliRegisterParam(show, &node)

		{
			var nodeName cmdparser.Param
			cmdparser.InitParam(&nodeName,
				cmdparser.LEAF,
				"",
				nil,
				validNodeName,
				cmdparser.STRING,
				"node-name",
				"Name of a node in the topology")
			cmdparser.LibcliRegisterParam(&node, &nodeName)

			{
				var arp cmdparser.Param
				cmdparser.InitParam(&arp,
					cmdparser.CMD,
					"arp",
					dumpArpTable,
					nil,
					cmdparser.INVALID,
					"",
					"Adress Resolution Protocol")
				cmdparser.LibcliRegisterParam(&nodeName, &arp)
				cmdparser.SetParamCmdCode(&arp, ARP_TABLE)
			}
		}
	}
	{
		var node cmdparser.Param
		cmdparser.InitParam(&node,
			cmdparser.CMD,
			"node",
			nil,
			nil,
			cmdparser.INVALID,
			"",
			"Given a node name and operation it performs that operation")
		cmdparser.LibcliRegisterParam(run, &node)

		{
			var nodeName cmdparser.Param
			cmdparser.InitParam(&nodeName,
				cmdparser.LEAF,
				"",
				nil,
				validNodeName,
				cmdparser.STRING,
				"node-name",
				"Name of a node in the topology")
			cmdparser.LibcliRegisterParam(&node, &nodeName)

			{
				var resolveArp cmdparser.Param
				cmdparser.InitParam(&resolveArp,
					cmdparser.CMD,
					"resolve-arp",
					nil,
					nil,
					cmdparser.INVALID,
					"",
					"resolves arp of a node")
				cmdparser.LibcliRegisterParam(&nodeName, &resolveArp)

				{
					var ipAddr cmdparser.Param
					cmdparser.InitParam(&ipAddr,
						cmdparser.LEAF,
						"",
						arpHandler,
						validIPAddr,
						cmdparser.STRING,
						"ip-addr",
						"Takes ip addr of node i.e loopback addr")
					cmdparser.LibcliRegisterParam(&resolveArp, &ipAddr)
					cmdparser.SetParamCmdCode(&ipAddr, ARP_HANDLER)
				}
			}

		}

	}
}
