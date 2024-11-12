package cli

import "github.com/gkarthikreddi/tcp/tools/cmdparser"

const (
	SHOW_TOPO      = 1
	ARP_HANDLER    = 2
	ARP_TABLE      = 3
	MAC_TABLE      = 4
	RT_TABLE       = 5
	L3_HANDLER     = 6
	PING_HANDLER   = 7
	ARPALL_HANDLER = 8
)

func InitNwCli() {
	cmdparser.InitLibcli()

	show := cmdparser.GetShowHook()
	run := cmdparser.GetRunHook()
	config := cmdparser.GetConfigHook()

	{
		var topo cmdparser.Param
		cmdparser.InitParam(&topo, // param
			cmdparser.CMD,                  // type of param
			"topology",                     // name of param, nil for leaf param
			showHandler,                    // callback handler
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
				var mac cmdparser.Param
				cmdparser.InitParam(&mac,
					cmdparser.CMD,
					"mac",
					showHandler,
					nil,
					cmdparser.INVALID,
					"",
					"Mac Table lookup of a node")
				cmdparser.LibcliRegisterParam(&nodeName, &mac)
				cmdparser.SetParamCmdCode(&mac, MAC_TABLE)
			}

			{
				var arp cmdparser.Param
				cmdparser.InitParam(&arp,
					cmdparser.CMD,
					"arp",
					showHandler,
					nil,
					cmdparser.INVALID,
					"",
					"Adress Resolution Protocol")
				cmdparser.LibcliRegisterParam(&nodeName, &arp)
				cmdparser.SetParamCmdCode(&arp, ARP_TABLE)
			}
			{
				var routingTable cmdparser.Param
				cmdparser.InitParam(&routingTable,
					cmdparser.CMD,
					"rt",
					showHandler,
					nil,
					cmdparser.INVALID,
					"",
					"L3 Routing Table")
				cmdparser.LibcliRegisterParam(&nodeName, &routingTable)
				cmdparser.SetParamCmdCode(&routingTable, RT_TABLE)
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
				{
					var all cmdparser.Param
					cmdparser.InitParam(&all,
						cmdparser.CMD,
						"all",
						arpHandler,
						nil,
						cmdparser.INVALID,
						"",
						"Resolves arp on all interfaces")
					cmdparser.LibcliRegisterParam(&resolveArp, &all)
					cmdparser.SetParamCmdCode(&all, ARPALL_HANDLER)
				}
			}

			{
				var ping cmdparser.Param
				cmdparser.InitParam(&ping,
					cmdparser.CMD,
					"ping",
					nil,
					nil,
					cmdparser.INVALID,
					"",
					"Ping function")
				cmdparser.LibcliRegisterParam(&nodeName, &ping)

				{
					var ipAddr cmdparser.Param
					cmdparser.InitParam(&ipAddr,
						cmdparser.LEAF,
						"",
						pingHandler,
						validIPAddr,
						cmdparser.STRING,
						"ip-addr",
						"Dst IPaddr for ping functionality")
					cmdparser.LibcliRegisterParam(&ping, &ipAddr)
					cmdparser.SetParamCmdCode(&ipAddr, PING_HANDLER)
				}
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
		cmdparser.LibcliRegisterParam(config, &node)

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
				var route cmdparser.Param
				cmdparser.InitParam(&route,
					cmdparser.CMD,
					"route",
					nil,
					nil,
					cmdparser.INVALID,
					"",
					"Routing Table")
				cmdparser.LibcliRegisterParam(&nodeName, &route)

				{
					var dst cmdparser.Param
					cmdparser.InitParam(&dst,
						cmdparser.LEAF,
						"",
						nil,
						validIPAddr,
						cmdparser.STRING,
						"dst",
						"Destination Ip Addr")
					cmdparser.LibcliRegisterParam(&route, &dst)

					{
						var mask cmdparser.Param
						cmdparser.InitParam(&mask,
							cmdparser.LEAF,
							"",
							nil,
							validMask,
							cmdparser.STRING,
							"mask",
							"Mask of Ip Addr")
						cmdparser.LibcliRegisterParam(&dst, &mask)

						{
							var gatewayIP cmdparser.Param
							cmdparser.InitParam(&gatewayIP,
								cmdparser.LEAF,
								"",
								nil,
								validIPAddr,
								cmdparser.STRING,
								"gw-ip",
								"Gateway IP Addr")
							cmdparser.LibcliRegisterParam(&mask, &gatewayIP)

							{
								var outIntf cmdparser.Param
								cmdparser.InitParam(&outIntf,
									cmdparser.LEAF,
									"",
									l3ConfigHandler,
									nil,
									cmdparser.STRING,
									"out-intf",
									"Outgoing interface")
								cmdparser.LibcliRegisterParam(&gatewayIP, &outIntf)
								cmdparser.SetParamCmdCode(&outIntf, L3_HANDLER)
							}
						}
					}
				}
			}
		}
	}
}
