/*
 * lib/colors.go
 * ANSI terminal colour and style escape codes.
 * Imported by every package that writes coloured output.
 */
package lib

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
