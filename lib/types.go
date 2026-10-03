/*
 * lib/types.go
 * All domain types and the BGP finite state machine enum
 */
package lib

// BGPSTA represents one of the six BGP session states defined in RFC 4271.
type BGPSTA int

const (
	Idle        BGPSTA = iota // No active connection attempt
	Connect                   // TCP connection being initiated
	Active                    // Retrying TCP connection
	OpenSent                  // BGP OPEN message sent
	OpenConfirm               // OPEN received, awaiting KEEPALIVE
	Established               // Session up — routes may be exchanged
)

// Route is one entry in a BGP Routing Information Base
type Route struct {
	Prefix    string // destination prefix  e.g. "203.0.113.0/24"
	Path      []int  // ASNs from nearest hop to origin
	NextHop   int    // next-hop AS number
	LocalPref int    // LOCAL_PREF
	Valid     bool   // false = withdrawn / inactive
	OriginAS  int    // ASN that first originated this prefix
	Src       string // "LOCAL" | "eBGP" | "iBGP" | "HIJACK"
}

// Peer describes one BGP peering session configured on an AS
type Peer struct {
	ASN   int    // peer Autonomous System Number
	State BGPSTA // current position in the BGP FSM
	EBGP  bool   // true = eBGP session, false = iBGP
	Sent  int    // BGP messages sent to this peer
	Recv  int    // BGP messages received from this peer
}

// AutonomousSystem models a single BGP speaking network
type AutonomousSystem struct {
	ASN      int      // globally unique AS Number
	Name     string   // human-readable label e.g. "AS1-TierOne-ISP"
	RID      string   // BGP Router-ID (dotted-IPv4 format)
	Prefixes []string // IP prefixes legitimately owned by this AS
	Peers    []Peer   // all configured BGP peering sessions
	Routes   []Route  // BGP Routing Information Base (RIB)
	Alive    bool     // false = offline / blackholed
	Desc     string   // short human-readable description
}

// Packet models a data packet being traced hop-by-hop through the network
type Packet struct {
	Src     string // source IP  (Router-ID of the originating AS)
	Dst     string // destination prefix being looked up
	SrcASN  int    // source AS number
	Hops    []int  // ordered AS numbers visited during forwarding
	Dropped bool   // true = packet was discarded before delivery
	Reason  string // human-readable drop reason (when Dropped = true)
}

// LogEntry is one timestamped line in the simulation event log
type LogEntry struct {
	Msg   string // message text
	Level int    // 0=INFO  1=WARN  2=ERROR  3=OK  4=BGP
	T     int    // simulation tick when this event occurred
}
