/*
 * util/input.go
 * Everything that is not a type definition, a colour code, or a display
 * function lives here.  Sections are clearly delimited:
 */
package util

import (
	"bgpsim/lib"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// ── Input / IO ────────────────────────────────────────────────────────────────

// stdinReader is the shared buffered reader for all interactive input.
var stdinReader = bufio.NewReader(os.Stdin)

// ReadInt reads one line from stdin and parses it as an integer.
// Returns a non-nil error when the line cannot be parsed.
func ReadInt() (int, error) {
	line, _ := stdinReader.ReadString('\n')
	return strconv.Atoi(strings.TrimSpace(line))
}

// ReadString reads and returns one trimmed line from stdin.
func ReadString() string {
	line, _ := stdinReader.ReadString('\n')
	return strings.TrimSpace(line)
}

// PauseKey prints a prompt and blocks until the user presses Enter.
func PauseKey() {
	fmt.Print(lib.DIM + "\n  [Press ENTER to continue...]" + lib.RST)
	stdinReader.ReadString('\n') //nolint:errcheck
}

// Clr clears the terminal screen on both Unix and Windows.
func Clr() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

// ── Logging ───────────────────────────────────────────────────────────────────

// LogAdd appends a new entry to the global event log.
// When the log exceeds LogCap the oldest entry is silently discarded.
func LogAdd(msg string, level int) {
	if len(lib.EventLog) >= lib.LogCap {
		lib.EventLog = lib.EventLog[1:]
	}
	lib.EventLog = append(lib.EventLog, lib.LogEntry{
		Msg: msg, Level: level, T: lib.Tick,
	})
	lib.Tick++
}

// LogColor returns the ANSI colour code for a given log level.
func LogColor(level int) string {
	switch level {
	case 0:
		return lib.CYN // INFO
	case 1:
		return lib.YLW // WARN
	case 2:
		return lib.RED // ERROR
	case 3:
		return lib.GRN // OK
	case 4:
		return lib.BLU // BGP
	default:
		return lib.WHT
	}
}

// LogLabel returns a fixed-width 5-character label for a given log level.
func LogLabel(level int) string {
	switch level {
	case 0:
		return "INFO "
	case 1:
		return "WARN "
	case 2:
		return "ERROR"
	case 3:
		return " OK  "
	case 4:
		return " BGP "
	default:
		return "  ?  "
	}
}

// ── Formatting ────────────────────────────────────────────────────────────────

// PathToStr converts an integer AS_PATH slice to a human-readable string
// such as "AS1->AS2->AS5".  Returns "(origin)" for an empty path.
func PathToStr(path []int) string {
	if len(path) == 0 {
		return "(origin)"
	}
	parts := make([]string, len(path))
	for i, asn := range path {
		parts[i] = fmt.Sprintf("AS%d", asn)
	}
	return strings.Join(parts, "->")
}

// Hdr prints a centred, ANSI-framed section-header box to stdout.
func Hdr(title string) {
	const w = 68
	tl := len(title)
	lp := (w - 2 - tl) / 2
	rp := w - 2 - tl - lp
	bar := strings.Repeat("-", w-2)
	fmt.Printf(lib.BOLD+lib.CYN+"\n  +%s+\n  |%s%s%s|\n  +%s+\n"+lib.RST+"\n",
		bar,
		strings.Repeat(" ", lp), title, strings.Repeat(" ", rp),
		bar)
}

// StateStr returns a fixed-width 11-character name for a BGP FSM state.
func StateStr(s lib.BGPState) string {
	switch s {
	case lib.Idle:
		return "IDLE       "
	case lib.Connect:
		return "CONNECT    "
	case lib.Active:
		return "ACTIVE     "
	case lib.OpenSent:
		return "OPENSENT   "
	case lib.OpenConfirm:
		return "OPENCONFIRM"
	case lib.Established:
		return "ESTABLISHED"
	default:
		return "UNKNOWN    "
	}
}

// StateColor returns the ANSI colour code associated with a BGP FSM state.
func StateColor(s lib.BGPState) string {
	switch s {
	case lib.Idle:
		return lib.RED
	case lib.Connect, lib.Active:
		return lib.YLW
	case lib.OpenSent, lib.OpenConfirm:
		return lib.MAG
	case lib.Established:
		return lib.GRN
	default:
		return lib.WHT
	}
}

// ── Network Lookup ────────────────────────────────────────────────────────────

// FindAS returns the slice index of the live (Alive=true) AS with the given
// ASN.  Returns -1 if the AS does not exist or is currently offline.
func FindAS(asn int) int {
	for i := range lib.Network {
		if lib.Network[i].ASN == asn && lib.Network[i].Alive {
			return i
		}
	}
	return -1
}

// FindPeer returns the index of the peer slot in Network[ai] that represents
// a session to peerASN.  Returns -1 if the peer is not configured.
func FindPeer(ai, peerASN int) int {
	for p := range lib.Network[ai].Peers {
		if lib.Network[ai].Peers[p].ASN == peerASN {
			return p
		}
	}
	return -1
}

// FindRoute returns the index of the first valid route matching prefix pfx in
// Network[ai].Routes.  Returns -1 if no valid route is found.
func FindRoute(ai int, pfx string) int {
	for r := range lib.Network[ai].Routes {
		if lib.Network[ai].Routes[r].Prefix == pfx &&
			lib.Network[ai].Routes[r].Valid {
			return r
		}
	}
	return -1
}

// InPath returns true if asn appears anywhere in the path slice.
// Used for AS_PATH loop prevention during BGP UPDATE processing.
func InPath(path []int, asn int) bool {
	for _, v := range path {
		if v == asn {
			return true
		}
	}
	return false
}

// ── Network Builder ───────────────────────────────────────────────────────────

// AddPeer appends a new IDLE peer entry to the AS at index ai.
func AddPeer(ai, peerASN int, ebgp bool) {
	lib.Network[ai].Peers = append(lib.Network[ai].Peers, lib.Peer{
		ASN: peerASN, State: lib.Idle, EBGP: ebgp,
	})
}

// AddLocalRoute installs a LOCAL-origin route into Network[ai]'s RIB.
// A LOCAL route has an empty AS_PATH; the AS itself is the origin.
func AddLocalRoute(ai int, pfx string) {
	lib.Network[ai].Routes = append(lib.Network[ai].Routes, lib.Route{
		Prefix:    pfx,
		Path:      []int{},
		NextHop:   lib.Network[ai].ASN,
		LocalPref: 100,
		Valid:     true,
		OriginAS:  lib.Network[ai].ASN,
		Src:       "LOCAL",
	})
}

// ── Network Init ──────────────────────────────────────────────────────────────

// InitSim resets all global state and rebuilds the 5-AS topology.
// Safe to call multiple times — used by Restore() after attacks.
//
// Topology
// ─────────────────────────────────────────────────────────────
//
//	[AS3 Enterprise-A]
//	       |
//	     eBGP
//	       |
//
// [AS1 TierOne] <─eBGP─> [AS2 Regional] <─eBGP─> [AS4 Enterprise-B]
//
//	  |       \                  |
//	eBGP       \──eBGP──\      eBGP
//	  |                  \       |
//
// [AS5 DataCenter/CDN] <──eBGP──/
func InitSim() {
	lib.Network = nil
	lib.EventLog = nil
	lib.Tick = 0

	// ── AS1: Tier-1 ISP ─────────────────────────────────────────────────
	lib.Network = append(lib.Network, lib.AutonomousSystem{
		ASN: 1, Name: "AS1-TierOne-ISP", RID: "1.1.1.1",
		Desc: "Tier-1 Internet Service Provider", Alive: true,
		Prefixes: []string{"10.0.0.0/8"},
	})
	AddLocalRoute(0, "10.0.0.0/8")

	// ── AS2: Regional ISP ────────────────────────────────────────────────
	lib.Network = append(lib.Network, lib.AutonomousSystem{
		ASN: 2, Name: "AS2-Regional-ISP", RID: "2.2.2.2",
		Desc: "Regional Internet Service Provider", Alive: true,
		Prefixes: []string{"172.16.0.0/16"},
	})
	AddLocalRoute(1, "172.16.0.0/16")

	// ── AS3: Enterprise A ────────────────────────────────────────────────
	lib.Network = append(lib.Network, lib.AutonomousSystem{
		ASN: 3, Name: "AS3-Enterprise-A", RID: "3.3.3.3",
		Desc: "Enterprise network, customer of AS1", Alive: true,
		Prefixes: []string{"192.168.1.0/24"},
	})
	AddLocalRoute(2, "192.168.1.0/24")

	// ── AS4: Enterprise B ────────────────────────────────────────────────
	lib.Network = append(lib.Network, lib.AutonomousSystem{
		ASN: 4, Name: "AS4-Enterprise-B", RID: "4.4.4.4",
		Desc: "Enterprise network, customer of AS2", Alive: true,
		Prefixes: []string{"192.168.2.0/24"},
	})
	AddLocalRoute(3, "192.168.2.0/24")

	// ── AS5: Data Centre / CDN ───────────────────────────────────────────
	lib.Network = append(lib.Network, lib.AutonomousSystem{
		ASN: 5, Name: "AS5-DataCenter-CDN", RID: "5.5.5.5",
		Desc: "Content Delivery Network / Data Centre", Alive: true,
		Prefixes: []string{"203.0.113.0/24", "198.51.100.0/24"},
	})
	AddLocalRoute(4, "203.0.113.0/24")
	AddLocalRoute(4, "198.51.100.0/24")

	// ── eBGP peering pairs (both directions must be registered) ──────────
	AddPeer(0, 2, true)
	AddPeer(1, 1, true) // AS1 <-> AS2
	AddPeer(0, 3, true)
	AddPeer(2, 1, true) // AS1 <-> AS3
	AddPeer(0, 5, true)
	AddPeer(4, 1, true) // AS1 <-> AS5
	AddPeer(1, 4, true)
	AddPeer(3, 2, true) // AS2 <-> AS4
	AddPeer(1, 5, true)
	AddPeer(4, 2, true) // AS2 <-> AS5
}

// ── BGP Session Establishment ─────────────────────────────────────────────────

// BGPEstablish steps every IDLE peer pair through the BGP FSM:
//
//	IDLE → CONNECT  (TCP SYN/SYN-ACK)
//	CONNECT → OPENSENT  (BGP OPEN message)
//	OPENSENT → ESTABLISHED  (KEEPALIVE exchange)
//
// Sessions that are already ESTABLISHED are silently skipped.
func BGPEstablish() {
	LogAdd("=== BGP Session Establishment Phase ===", 0)

	for i := range lib.Network {
		if !lib.Network[i].Alive {
			continue
		}
		for p := range lib.Network[i].Peers {
			pe := &lib.Network[i].Peers[p]
			if pe.State != lib.Idle {
				continue
			}
			j := FindAS(pe.ASN)
			if j < 0 {
				continue
			}
			q := FindPeer(j, lib.Network[i].ASN)
			if q < 0 {
				continue
			}

			// ── IDLE → CONNECT ────────────────────────────────────────────
			pe.State = lib.Connect
			lib.Network[j].Peers[q].State = lib.Connect
			LogAdd(fmt.Sprintf("[TCP] AS%d initiated TCP connection to AS%d",
				lib.Network[i].ASN, pe.ASN), 0)

			// ── CONNECT → OPENSENT ────────────────────────────────────────
			pe.State = lib.OpenSent
			LogAdd(fmt.Sprintf("[OPEN] AS%d --> AS%d  Router-ID:%s  HoldTime:180s",
				lib.Network[i].ASN, pe.ASN, lib.Network[i].RID), 4)
			pe.Sent++

			lib.Network[j].Peers[q].State = lib.OpenSent
			LogAdd(fmt.Sprintf("[OPEN] AS%d --> AS%d  Router-ID:%s  HoldTime:180s",
				pe.ASN, lib.Network[i].ASN, lib.Network[j].RID), 4)
			lib.Network[j].Peers[q].Sent++
			pe.Recv++

			// ── OPENSENT → ESTABLISHED ────────────────────────────────────
			pe.State = lib.Established
			lib.Network[j].Peers[q].State = lib.Established
			pe.Sent++
			pe.Recv++
			lib.Network[j].Peers[q].Sent++
			lib.Network[j].Peers[q].Recv++
			LogAdd(fmt.Sprintf("[KEEPALIVE] AS%d <--> AS%d  ==> ESTABLISHED",
				lib.Network[i].ASN, pe.ASN), 3)
		}
	}
	LogAdd("=== All BGP sessions are now ESTABLISHED ===", 3)
}

// ── BGP Route Advertisement ───────────────────────────────────────────────────

// BGPSend attempts to propagate a BGP UPDATE for src from AS[fi] to AS[ti].
//
// Four rules are enforced:
//  1. AS_PATH loop prevention — drop if ti's ASN is already in the path
//  2. Session guard           — session must be ESTABLISHED
//  3. Best-path selection     — keep existing if strictly shorter AS_PATH
//  4. AS_PATH prepend         — sender prepends its own ASN
func BGPSend(fi, ti int, src lib.Route) {
	// Rule 1: loop prevention
	if InPath(src.Path, lib.Network[ti].ASN) {
		return
	}
	// Rule 2: session guard
	fp := FindPeer(fi, lib.Network[ti].ASN)
	if fp < 0 || lib.Network[fi].Peers[fp].State != lib.Established {
		return
	}
	newPathLen := len(src.Path) + 1
	// Rule 3: best-path check
	ri := FindRoute(ti, src.Prefix)
	if ri >= 0 && len(lib.Network[ti].Routes[ri].Path) < newPathLen {
		return
	}
	// Rule 4: build updated route with sender's ASN prepended
	newPath := make([]int, newPathLen)
	newPath[0] = lib.Network[fi].ASN
	copy(newPath[1:], src.Path)

	srcLabel := "eBGP"
	if !lib.Network[fi].Peers[fp].EBGP {
		srcLabel = "iBGP"
	}
	updated := lib.Route{
		Prefix: src.Prefix, Path: newPath,
		NextHop: lib.Network[fi].ASN, LocalPref: 100,
		Valid: true, OriginAS: src.OriginAS, Src: srcLabel,
	}
	if ri < 0 {
		lib.Network[ti].Routes = append(lib.Network[ti].Routes, updated)
	} else {
		lib.Network[ti].Routes[ri] = updated
	}
	lib.Network[fi].Peers[fp].Sent++
	if rp := FindPeer(ti, lib.Network[fi].ASN); rp >= 0 {
		lib.Network[ti].Peers[rp].Recv++
	}
	LogAdd(fmt.Sprintf("[UPDATE] AS%d --> AS%d  prefix:%-22s  path:[%s]",
		lib.Network[fi].ASN, lib.Network[ti].ASN,
		src.Prefix, PathToStr(newPath)), 4)
}

// BGPAdvertise runs the full route-advertisement cycle until convergence.
//
//   - Phase 1    : each AS announces its LOCAL-origin prefixes to peers
//   - Phases 2–5 : repeated passes propagate learned routes further
func BGPAdvertise() {
	LogAdd("=== BGP Route Advertisement ===", 0)

	for i := range lib.Network {
		if !lib.Network[i].Alive {
			continue
		}
		snapshot := make([]lib.Route, len(lib.Network[i].Routes))
		copy(snapshot, lib.Network[i].Routes)
		for _, r := range snapshot {
			if !r.Valid || r.Src != "LOCAL" {
				continue
			}
			for p := range lib.Network[i].Peers {
				if lib.Network[i].Peers[p].State != lib.Established {
					continue
				}
				if j := FindAS(lib.Network[i].Peers[p].ASN); j >= 0 {
					BGPSend(i, j, r)
				}
			}
		}
	}

	for pass := 0; pass < 5; pass++ {
		for i := range lib.Network {
			if !lib.Network[i].Alive {
				continue
			}
			snapshot := make([]lib.Route, len(lib.Network[i].Routes))
			copy(snapshot, lib.Network[i].Routes)
			for _, r := range snapshot {
				if !r.Valid {
					continue
				}
				for p := range lib.Network[i].Peers {
					if lib.Network[i].Peers[p].State != lib.Established {
						continue
					}
					if j := FindAS(lib.Network[i].Peers[p].ASN); j >= 0 {
						BGPSend(i, j, r)
					}
				}
			}
		}
	}
	LogAdd("=== Network converged: all prefixes reachable ===", 3)
}

// ── Packet Routing ────────────────────────────────────────────────────────────

// bestRoute returns the index of the best valid route for dst in Network[ai].
// Selection criteria: highest LOCAL_PREF first, then shortest AS_PATH.
// Returns -1 if no valid route exists.
func bestRoute(ai int, dst string) int {
	best, bestLen, bestLP := -1, 9999, -1
	for r := range lib.Network[ai].Routes {
		rt := &lib.Network[ai].Routes[r]
		if !rt.Valid || rt.Prefix != dst {
			continue
		}
		pLen := len(rt.Path)
		if rt.LocalPref > bestLP ||
			(rt.LocalPref == bestLP && pLen < bestLen) {
			best, bestLen, bestLP = r, pLen, rt.LocalPref
		}
	}
	return best
}

// SendPacket simulates forwarding a packet from srcASN toward dst.
// Returns a Packet recording the full path taken and the final outcome.
func SendPacket(srcASN int, dst string) lib.Packet {
	pkt := lib.Packet{SrcASN: srcASN, Dst: dst}

	ci := FindAS(srcASN)
	if ci < 0 {
		pkt.Dropped = true
		pkt.Reason = "Source AS not found or offline"
		return pkt
	}
	pkt.Src = lib.Network[ci].RID
	pkt.Hops = append(pkt.Hops, srcASN)
	LogAdd(fmt.Sprintf("[PKT] New packet  src:AS%d(%s)  dst:%s",
		srcASN, lib.Network[ci].Name, dst), 0)

	for hop := 0; hop < lib.MaxPath; hop++ {
		// ── Delivery: current AS owns this prefix ─────────────────────────
		for _, pfx := range lib.Network[ci].Prefixes {
			if pfx == dst {
				LogAdd(fmt.Sprintf("[PKT] DELIVERED to AS%d (%s) in %d hop(s)",
					lib.Network[ci].ASN, lib.Network[ci].Name, hop), 3)
				return pkt
			}
		}
		// ── Route lookup ──────────────────────────────────────────────────
		ri := bestRoute(ci, dst)
		if ri < 0 {
			LogAdd(fmt.Sprintf("[PKT] DROPPED at AS%d: no route to %s",
				lib.Network[ci].ASN, dst), 2)
			pkt.Dropped = true
			pkt.Reason = fmt.Sprintf("No route to %s at AS%d", dst, lib.Network[ci].ASN)
			return pkt
		}
		nh := lib.Network[ci].Routes[ri].NextHop
		// ── Next-hop = self: hijack capture or loop ───────────────────────
		if nh == lib.Network[ci].ASN {
			if lib.Network[ci].Routes[ri].Src == "HIJACK" {
				LogAdd(fmt.Sprintf(
					"[PKT] CAPTURED at hijacker AS%d — traffic blackholed!",
					lib.Network[ci].ASN), 2)
				pkt.Reason = fmt.Sprintf(
					"Captured by hijacker AS%d (prefix not legitimately owned)",
					lib.Network[ci].ASN)
			} else {
				LogAdd(fmt.Sprintf(
					"[PKT] DROPPED at AS%d: routing loop (next-hop = self)",
					lib.Network[ci].ASN), 2)
				pkt.Reason = "Routing loop: next-hop equals current AS"
			}
			pkt.Dropped = true
			return pkt
		}
		// ── Forward to next hop ───────────────────────────────────────────
		LogAdd(fmt.Sprintf("[PKT] AS%d --> AS%d  (route-source: %s)",
			lib.Network[ci].ASN, nh, lib.Network[ci].Routes[ri].Src), 0)
		ni := FindAS(nh)
		if ni < 0 {
			LogAdd(fmt.Sprintf("[PKT] DROPPED: AS%d is unreachable / offline", nh), 2)
			pkt.Dropped = true
			pkt.Reason = fmt.Sprintf("AS%d is unreachable or offline", nh)
			return pkt
		}
		pkt.Hops = append(pkt.Hops, nh)
		ci = ni
	}

	pkt.Dropped = true
	pkt.Reason = "TTL exceeded (too many hops)"
	LogAdd("[PKT] DROPPED: TTL exceeded", 2)
	return pkt
}

// ── Attacks ───────────────────────────────────────────────────────────────────

// SimHijack causes atkASN to falsely originate victimPfx and propagate the
// fraudulent BGP UPDATE to all ESTABLISHED peers.
//
// The injected route carries an empty AS_PATH — making the attacker appear
// as the direct origin.  Without RPKI / ROA, neighbours accept it.
//
// Real-world mitigation: deploy RPKI and publish ROA records.
func SimHijack(atkASN int, victimPfx string) {
	ai := FindAS(atkASN)
	if ai < 0 {
		LogAdd("[HIJACK] Attacker AS not found or offline", 2)
		return
	}
	fake := lib.Route{
		Prefix: victimPfx, Path: []int{}, NextHop: atkASN,
		LocalPref: 100, Valid: true, OriginAS: atkASN, Src: "HIJACK",
	}
	if ri := FindRoute(ai, victimPfx); ri < 0 {
		lib.Network[ai].Routes = append(lib.Network[ai].Routes, fake)
	} else {
		lib.Network[ai].Routes[ri] = fake
	}

	for p := range lib.Network[ai].Peers {
		if lib.Network[ai].Peers[p].State != lib.Established {
			continue
		}
		j := FindAS(lib.Network[ai].Peers[p].ASN)
		if j < 0 {
			continue
		}
		BGPSend(ai, j, fake)
		LogAdd(fmt.Sprintf("[HIJACK] Fake UPDATE for %-22s propagated to AS%d",
			victimPfx, lib.Network[ai].Peers[p].ASN), 2)
	}
	LogAdd("[HIJACK] Done. Use option 3 to inspect routing tables.", 2)
}

// SimBlackhole takes tgtASN offline by sending BGP NOTIFICATION CEASE to all
// peers, withdrawing every route, and marking the AS as unreachable.
func SimBlackhole(tgtASN int) {
	ti := FindAS(tgtASN)
	if ti < 0 {
		LogAdd("[BLACKHOLE] Target AS not found or already offline", 2)
		return
	}

	for p := range lib.Network[ti].Peers {
		pe := &lib.Network[ti].Peers[p]
		if pe.State != lib.Established {
			continue
		}
		j := FindAS(pe.ASN)
		if j < 0 {
			continue
		}
		q := FindPeer(j, tgtASN)
		LogAdd(fmt.Sprintf("[NOTIFICATION] AS%d --> AS%d  CEASE (session terminated)",
			tgtASN, pe.ASN), 2)
		pe.State = lib.Idle
		if q >= 0 {
			lib.Network[j].Peers[q].State = lib.Idle
		}
		for r := range lib.Network[j].Routes {
			rt := &lib.Network[j].Routes[r]
			if rt.NextHop == tgtASN && rt.Valid {
				LogAdd(fmt.Sprintf("[WITHDRAW] AS%d removing %-22s (was via AS%d)",
					lib.Network[j].ASN, rt.Prefix, tgtASN), 1)
				rt.Valid = false
			}
		}
	}

	for i := range lib.Network {
		if i == ti {
			continue
		}
		for r := range lib.Network[i].Routes {
			for _, tgtPfx := range lib.Network[ti].Prefixes {
				if lib.Network[i].Routes[r].Prefix == tgtPfx {
					lib.Network[i].Routes[r].Valid = false
				}
			}
		}
	}
	lib.Network[ti].Alive = false
	LogAdd(fmt.Sprintf(
		"[BLACKHOLE] AS%d",
		tgtASN), 2)
}

// Restore reinitialises the simulation and reconverges the network from
// scratch, returning every AS to a fully operational state.
func Restore() {
	LogAdd("=== Restoring network to default state ===", 1)
	InitSim()
	BGPEstablish()
	BGPAdvertise()
}
