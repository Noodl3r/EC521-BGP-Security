/*
 * lib/constants.go
 * ANSI terminal colour styles
 */
package lib

// Color codes
const (
	RST  = "\033[0m"  // reset all attributes
	BOLD = "\033[1m"  // bold text
	DIM  = "\033[2m"  // faint / dim text
	RED  = "\033[91m" // bright red
	GRN  = "\033[92m" // bright green
	YLW  = "\033[93m" // bright yellow
	BLU  = "\033[94m" // bright blue
	MAG  = "\033[95m" // bright magenta
	CYN  = "\033[96m" // bright cyan
	WHT  = "\033[97m" // bright white
)

// Limits
const (
	MaxPath = 10  // maximum AS hops before a packet is dropped (TTL)
	LogCap  = 200 // event log capacity; oldest entry discarded when full
)
