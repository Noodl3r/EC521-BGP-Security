/*
 * main/display/menu.go
 * All CLI visualisation functions, the interactive menu, per-option
 * action handlers, main run loop.
 */
package display

import (
	"bgpsim/lib"
	"bgpsim/util"
	"fmt"
)

// ── ShowTopology ──────────────────────────────────────────────────────────────

// ShowTopology prints the ASCII network diagram and the AS summary table.
func ShowTopology() {
	util.Hdr("NETWORK TOPOLOGY")
	fmt.Printf(lib.BOLD + "  Topology Diagram:\n\n" + lib.RST)
	fmt.Printf(
		"           " + lib.GRN + "[AS3 Enterprise-A]\n" + lib.RST +
			"                  |\n" +
			"                  | eBGP\n" +
			"                  |\n" +
			lib.CYN + "  [AS1 TierOne ISP]" + lib.RST +
			"  <----eBGP---->  " +
			lib.CYN + "[AS2 Regional ISP]" + lib.RST +
			"  <----eBGP---->  " +
			lib.GRN + "[AS4 Enterprise-B]\n" + lib.RST +
			"        |        \\\n" +
			"      eBGP         \\----eBGP----\\         eBGP (via AS2)\n" +
			"        |                        \\            |\n" +
			lib.YLW + "  [AS5 DataCenter/CDN]" + lib.RST +
			"  <-----------eBGP-----------/\n\n")

	fmt.Printf(lib.BOLD + "  AS Summary Table:\n\n" + lib.RST)
	fmt.Printf(lib.CYN +
		"  +--------+--------------------+------------+--------------------------+\n" +
		"  |  ASN   | Name               | Router-ID  | Owned Prefixes           |\n" +
		"  +--------+--------------------+------------+--------------------------+\n" + lib.RST)

	for i := range lib.Network {
		statusColor := lib.GRN
		if !lib.Network[i].Alive {
			statusColor = lib.RED
		}
		firstPfx := ""
		if len(lib.Network[i].Prefixes) > 0 {
			firstPfx = lib.Network[i].Prefixes[0]
		}
		fmt.Printf(lib.CYN+"  | "+lib.RST+"%sAS%-5d%s"+lib.CYN+" | "+lib.RST+
			"%-18s"+lib.CYN+" | "+lib.RST+"%-10s"+lib.CYN+" | "+lib.RST+
			"%-24s"+lib.CYN+" |\n"+lib.RST,
			statusColor, lib.Network[i].ASN, lib.RST,
			lib.Network[i].Name, lib.Network[i].RID, firstPfx)
		for p := 1; p < len(lib.Network[i].Prefixes); p++ {
			fmt.Printf(lib.CYN+
				"  |        |                    |            | "+lib.RST+
				"%-24s"+lib.CYN+" |\n"+lib.RST, lib.Network[i].Prefixes[p])
		}
		if !lib.Network[i].Alive {
			fmt.Printf(lib.CYN +
				"  |        |                    |            | " + lib.RED +
				"** OFFLINE / BLACKHOLED **  " + lib.RST + lib.CYN + "|\n" + lib.RST)
		}
	}
	fmt.Printf(lib.CYN +
		"  +--------+--------------------+------------+--------------------------+\n" + lib.RST)
}

// ── ShowSessions ──────────────────────────────────────────────────────────────

// ShowSessions prints the BGP session status for every peer relationship.
func ShowSessions() {
	util.Hdr("BGP SESSION STATUS")
	fmt.Printf(lib.CYN +
		"  +--------------------+--------------------+---------+" +
		"---------------+------+------+\n" +
		"  | Local AS           | Peer AS            | Type    |" +
		" State         | Sent | Recv |\n" +
		"  +--------------------+--------------------+---------+" +
		"---------------+------+------+\n" + lib.RST)

	for i := range lib.Network {
		for p := range lib.Network[i].Peers {
			pe := &lib.Network[i].Peers[p]
			peerName := "Unknown"
			if j := util.FindAS(pe.ASN); j >= 0 {
				peerName = lib.Network[j].Name
			}
			typeStr := "iBGP"
			if pe.EBGP {
				typeStr = "eBGP"
			}
			fmt.Printf(lib.CYN+"  | "+lib.RST+"%-18s"+lib.CYN+" | "+lib.RST+
				"%-18s"+lib.CYN+" | "+lib.RST+"%-7s"+lib.CYN+" | "+lib.RST,
				lib.Network[i].Name, peerName, typeStr)
			fmt.Printf("%s%-13s%s"+lib.CYN+" | "+lib.RST+"%-4d"+lib.CYN+
				" | "+lib.RST+"%-4d"+lib.CYN+" |\n"+lib.RST,
				util.StateColor(pe.State), util.StateStr(pe.State), lib.RST,
				pe.Sent, pe.Recv)
		}
	}
	fmt.Printf(lib.CYN +
		"  +--------------------+--------------------+---------+" +
		"---------------+------+------+\n" + lib.RST)
	fmt.Printf("\n  " + lib.BOLD + "State colours: " + lib.RST +
		lib.RED + "IDLE " + lib.RST + "| " + lib.YLW + "CONNECT/ACTIVE " + lib.RST +
		"| " + lib.MAG + "OPENSENT/OPENCONFIRM " + lib.RST +
		"| " + lib.GRN + "ESTABLISHED\n" + lib.RST)
}

// ── ShowRoutes ────────────────────────────────────────────────────────────────

// ShowRoutes prints the full BGP RIB (Routing Information Base) for every AS.
func ShowRoutes() {
	util.Hdr("BGP ROUTING TABLES")

	for i := range lib.Network {
		aliveColor := lib.GRN
		if !lib.Network[i].Alive {
			aliveColor = lib.RED
		}
		fmt.Printf("\n  "+lib.BOLD+"%s[AS%d] %s%s\n"+lib.RST,
			aliveColor, lib.Network[i].ASN, lib.Network[i].Name, lib.RST)
		if !lib.Network[i].Alive {
			fmt.Printf(lib.RED +
				"    ** OFFLINE: blackholed / unreachable **\n" + lib.RST)
			continue
		}
		fmt.Printf(lib.CYN +
			"  +---------------------+------------------------------+" +
			"--------+-------+\n" +
			"  | Prefix              | AS Path                      |" +
			" Source | LPref |\n" +
			"  +---------------------+------------------------------+" +
			"--------+-------+\n" + lib.RST)

		shown := 0
		for r := range lib.Network[i].Routes {
			rt := &lib.Network[i].Routes[r]
			if !rt.Valid {
				continue
			}
			shown++
			srcColor := lib.CYN
			switch rt.Src {
			case "LOCAL":
				srcColor = lib.GRN
			case "HIJACK":
				srcColor = lib.RED
			}
			fmt.Printf(lib.CYN+"  | "+lib.RST+"%-19s"+lib.CYN+" | "+lib.RST+
				"%-28s"+lib.CYN+" | %s%-6s%s"+lib.CYN+" | "+lib.RST+
				"%-5d"+lib.CYN+" |\n"+lib.RST,
				rt.Prefix, util.PathToStr(rt.Path),
				srcColor, rt.Src, lib.RST, rt.LocalPref)
		}
		if shown == 0 {
			fmt.Printf("  |  (no routes — run BGP advertisement first)" +
				"                              |\n")
		}
		fmt.Printf(lib.CYN +
			"  +---------------------+------------------------------+" +
			"--------+-------+\n" + lib.RST)
	}
	fmt.Printf(lib.DIM + "\n  Source: " + lib.RST +
		lib.GRN + "LOCAL" + lib.RST + " = owned  |  " +
		lib.CYN + "eBGP" + lib.RST + " = learned externally  |  " +
		lib.RED + "HIJACK" + lib.RST + " = fraudulent announcement\n" + lib.RST)
}

// ── ShowLogs ──────────────────────────────────────────────────────────────────

// ShowLogs prints the most recent 40 entries from the simulation event log.
func ShowLogs() {
	util.Hdr("SIMULATION EVENT LOG")
	start := 0
	if len(lib.EventLog) > 40 {
		start = len(lib.EventLog) - 40
	}
	for i := start; i < len(lib.EventLog); i++ {
		e := &lib.EventLog[i]
		fmt.Printf("  %s[T%03d][%s]%s %s\n",
			util.LogColor(e.Level), e.T,
			util.LogLabel(e.Level), lib.RST, e.Msg)
	}
	fmt.Printf(lib.DIM+
		"\n  Showing %d of %d events  |  Levels: "+lib.RST+
		lib.CYN+"INFO "+lib.RST+
		lib.YLW+"WARN "+lib.RST+
		lib.RED+"ERROR "+lib.RST+
		lib.GRN+"OK "+lib.RST+
		lib.BLU+"BGP\n"+lib.RST,
		len(lib.EventLog)-start, len(lib.EventLog))
}

// ── ShowPkt ───────────────────────────────────────────────────────────────────

// ShowPkt prints the full result of a single packet-routing trace.
func ShowPkt(p *lib.Packet) {
	util.Hdr("PACKET ROUTING TRACE")
	srcName := "?"
	if si := util.FindAS(p.SrcASN); si >= 0 {
		srcName = lib.Network[si].Name
	}
	fmt.Printf("  Source IP    : "+lib.YLW+"%s"+lib.RST+"  (AS%d – %s)\n",
		p.Src, p.SrcASN, srcName)
	fmt.Printf("  Destination  : "+lib.YLW+"%s\n"+lib.RST, p.Dst)
	fmt.Printf("\n  Path taken   : ")
	for i, hop := range p.Hops {
		fmt.Printf(lib.CYN+"AS%d"+lib.RST, hop)
		if idx := util.FindAS(hop); idx >= 0 {
			fmt.Printf("(%s)", lib.Network[idx].Name)
		}
		if i < len(p.Hops)-1 {
			fmt.Printf(lib.YLW + " ---> " + lib.RST)
		}
	}
	fmt.Printf("\n\n  Result       : ")
	if p.Dropped {
		fmt.Printf(lib.RED+lib.BOLD+"DROPPED    "+lib.RST+
			lib.RED+"%s\n"+lib.RST, p.Reason)
	} else {
		fmt.Printf(lib.GRN + lib.BOLD + "DELIVERED SUCCESSFULLY\n" + lib.RST)
	}
	fmt.Printf("\n")
}

// ── ShowMenu ──────────────────────────────────────────────────────────────────

// ShowMenu renders the numbered option list to stdout.
func ShowMenu() {
	fmt.Printf(lib.BOLD + lib.CYN +
		"\n  +=============================================================+\n" +
		"  |              BGP SIMULATION  --  MAIN MENU                 |\n" +
		"  +=============================================================+\n" + lib.RST)

	fmt.Printf(lib.CYN + "  |" + lib.RST + lib.BOLD +
		"  VIEW ──────────────────────────────────────────────────" + lib.CYN + "|\n" + lib.RST)
	fmt.Printf(lib.CYN + "  |" + lib.RST +
		"   1)  Network Topology                                " + lib.CYN + "|\n" + lib.RST)
	fmt.Printf(lib.CYN + "  |" + lib.RST +
		"   2)  BGP Session Status                              " + lib.CYN + "|\n" + lib.RST)
	fmt.Printf(lib.CYN + "  |" + lib.RST +
		"   3)  BGP Routing Tables.                             " + lib.CYN + "|\n" + lib.RST)
	fmt.Printf(lib.CYN + "  |" + lib.RST +
		"   4)  Event Log                                       " + lib.CYN + "|\n" + lib.RST)

	fmt.Printf(lib.CYN + "  |" + lib.RST + lib.BOLD +
		"  BGP OPERATIONS ────────────────────────────────────────" + lib.CYN + "|\n" + lib.RST)
	fmt.Printf(lib.CYN + "  |" + lib.RST +
		"   5)  Establish BGP Sessions                          " + lib.CYN + "|\n" + lib.RST)
	fmt.Printf(lib.CYN + "  |" + lib.RST +
		"   6)  Run BGP Route Advertisement                     " + lib.CYN + "|\n" + lib.RST)

	fmt.Printf(lib.CYN + "  |" + lib.RST + lib.BOLD +
		"  PACKET SIMULATION ──────────────────────────────────────" + lib.CYN + "|\n" + lib.RST)
	fmt.Printf(lib.CYN + "  |" + lib.RST +
		"   7)  Send a Packet                                   " + lib.CYN + "|\n" + lib.RST)

	fmt.Printf(lib.CYN + "  |" + lib.RST + lib.BOLD +
		"  ATTACK SIMULATIONS ─────────────────────────────────────" + lib.CYN + "|\n" + lib.RST)
	fmt.Printf(lib.CYN + "  |" + lib.RST + lib.RED +
		"   8)  Simulate BGP Route Hijack                       " + lib.CYN + "|\n" + lib.RST)
	fmt.Printf(lib.CYN + "  |" + lib.RST + lib.YLW +
		"   9)  Simulate BGP Blackhole                          " + lib.CYN + "|\n" + lib.RST)

	fmt.Printf(lib.CYN + "  |" + lib.RST + lib.BOLD +
		"  SYSTEM ──────────────────────────────────────────" + lib.CYN + "|\n" + lib.RST)
	fmt.Printf(lib.CYN + "  |" + lib.RST +
		"  10)  Restore Network to Normal State                 " + lib.CYN + "|\n" + lib.RST)
	fmt.Printf(lib.CYN + "  |" + lib.RST +
		"  11)  Exit                                            " + lib.CYN + "|\n" + lib.RST)
	fmt.Printf(lib.CYN +
		"  +=============================================================+\n" + lib.RST)
	fmt.Printf(lib.BOLD + "Enter choice: " + lib.RST)
}

// ── Interactive Action Handlers ───────────────────────────────────────────────

// DoSend handles option 7: interactive hop-by-hop packet trace.
func DoSend() {
	util.Hdr("SEND A PACKET")
	fmt.Printf("  Active ASes (possible sources):\n")
	for i := range lib.Network {
		if lib.Network[i].Alive {
			fmt.Printf("    AS%d  %-22s  Router-ID: %s\n",
				lib.Network[i].ASN, lib.Network[i].Name, lib.Network[i].RID)
		}
	}
	fmt.Printf("\n  Enter source AS number: ")
	src, err := util.ReadInt()
	if err != nil {
		fmt.Printf(lib.RED + "  Invalid input — please enter a valid AS number.\n" + lib.RST)
		return
	}
	fmt.Printf("\n  Available destination prefixes:\n")
	for i := range lib.Network {
		if lib.Network[i].Alive {
			for _, pfx := range lib.Network[i].Prefixes {
				fmt.Printf("    [AS%d]  %s\n", lib.Network[i].ASN, pfx)
			}
		}
	}
	fmt.Printf("\n  Enter destination prefix (e.g. 192.168.2.0/24): ")
	dst := util.ReadString()
	fmt.Println()
	pkt := util.SendPacket(src, dst)
	ShowPkt(&pkt)
}

// DoHijack handles option 8: interactive BGP route-hijack simulation.
func DoHijack() {
	util.Hdr("BGP ROUTE HIJACK")
	fmt.Printf("  Active ASes:\n")
	for i := range lib.Network {
		if lib.Network[i].Alive {
			fmt.Printf("    AS%d  %s\n", lib.Network[i].ASN, lib.Network[i].Name)
		}
	}
	fmt.Printf("\n  Enter attacker AS number: ")
	atk, err := util.ReadInt()
	if err != nil {
		fmt.Printf(lib.RED + "  Invalid input.\n" + lib.RST)
		return
	}
	fmt.Printf("  Enter victim prefix to hijack (e.g. 203.0.113.0/24): ")
	pfx := util.ReadString()
	util.SimHijack(atk, pfx)
	fmt.Printf(lib.RED +
		"\n  Hijack injected.  Use option 3 to inspect routing tables.\n" + lib.RST)
}

// DoBlackhole handles option 9: interactive BGP blackhole simulation.
func DoBlackhole() {
	util.Hdr("BGP BLACKHOLE")
	fmt.Printf("  Active ASes:\n")
	for i := range lib.Network {
		if lib.Network[i].Alive {
			fmt.Printf("    AS%d  %s\n", lib.Network[i].ASN, lib.Network[i].Name)
		}
	}
	fmt.Printf("\n  Enter AS number to blackhole: ")
	tgt, err := util.ReadInt()
	if err != nil {
		fmt.Printf(lib.RED + "  Invalid input.\n" + lib.RST)
		return
	}
	util.SimBlackhole(tgt)
	fmt.Printf(lib.RED+
		"\n  Blackhole applied.  AS%d is now unreachable.\n"+lib.RST, tgt)
}

// ── Run ───────────────────────────────────────────────────────────────────────

// Run starts the interactive simulation loop.
// It does not return until the user selects option 12 (Exit).
func Run() {
	for {
		ShowMenu()

		choice, err := util.ReadInt()
		if err != nil {
			fmt.Printf(lib.RED +
				"  Invalid input — please enter a number between 1 and 12.\n" + lib.RST)
			util.PauseKey()
			continue
		}
		util.Clr()

		switch choice {
		case 1:
			ShowTopology()
			util.PauseKey()
		case 2:
			ShowSessions()
			util.PauseKey()
		case 3:
			ShowRoutes()
			util.PauseKey()
		case 4:
			ShowLogs()
			util.PauseKey()
		case 5:
			util.BGPEstablish()
			ShowSessions()
			fmt.Printf(lib.GRN + "\n  BGP sessions established!\n" + lib.RST)
			util.PauseKey()
		case 6:
			util.BGPAdvertise()
			ShowRoutes()
			fmt.Printf(lib.GRN +
				"\n  Route advertisement complete.  Network converged.\n" + lib.RST)
			util.PauseKey()
		case 7:
			DoSend()
			util.PauseKey()
		case 8:
			DoHijack()
			util.PauseKey()
		case 9:
			DoBlackhole()
			util.PauseKey()
		case 10:
			util.Restore()
			fmt.Printf(lib.GRN + "\n  Network restored to default state.\n" + lib.RST)
			util.PauseKey()
		case 11:
			fmt.Printf(lib.CYN + "\n  Goodbye!  BGP Simulation terminated.\n\n" + lib.RST)
			return
		default:
			fmt.Printf(lib.RED +
				"  Unknown option\n" + lib.RST)
			util.PauseKey()
		}
	}
}
