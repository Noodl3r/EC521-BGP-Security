/*
 * main/main.go
 * Entry point.  Initialises the simulation, shows the welcome screen,
 * then hands control to display.Run() for the interactive menu loop.
 *
 *  Build:  go build -o bgp_sim ./main/
 *  Run:    go run ./main/
 *
 *  Recommended first run : option 10  (Full guided demo)
 *  Manual exploration     : options 5 → 6 → 7 in sequence
 */
package main

import (
	"bgpsim/lib"
	"bgpsim/main/display"
	"bgpsim/util"
	"fmt"
)

func main() {
	util.InitSim()

	fmt.Printf(lib.GRN + "Ready\n" + lib.RST)

	util.PauseKey()
	display.Run()
}
