/*
 * main.go
 * Entry point.  Initialises the simulation, shows the startup screen,
 * then hands off to menu.Run() for the interactive loop.
 *
 *  Build:   go build -o bgp_sim .
 *  Run:     ./bgp_sim
 *  Or:      go run .
 *
 *  Recommended first run: option 10 (Full Demo).
 *  Manual flow: options 5 → 6 → 7 in sequence.
 */
package main

import (
	"bgpsim/lib"
	"bgpsim/main/display"
	"bgpsim/util"
	"fmt"
)

func main() {
	fmt.Printf(lib.GRN + "[+] Simulation initialised.\n" + lib.RST)
	fmt.Printf("[*] BGP sessions not yet established.  " + "No routes advertised yet.\n")

	util.PauseKey()
	display.Run()
}
