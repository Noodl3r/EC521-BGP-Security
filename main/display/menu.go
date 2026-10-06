/*
 * menu/runner.go
 *
 */
package display

import (
	"bgpsim/lib"
	"bgpsim/util"
	"fmt"
)

func ShowMenu() {
	fmt.Printf(lib.BOLD + lib.CYN +
		"\n+=============================================================+\n" +
		"                 BGP SIMULATION      MAIN MENU                  \n" +
		"+=============================================================+\n" + lib.RST)

	fmt.Printf(lib.RST + lib.BOLD + "VIEW\n" + lib.RST)
	fmt.Printf(lib.RST + " 1) Topology\n" + lib.RST)
	fmt.Printf(lib.RST + " 2) BGP Status\n" + lib.RST)
	fmt.Printf(lib.RST + " 3) Routing Tables\n" + lib.RST)
	fmt.Printf(lib.RST + " 4) Event Log\n" + lib.RST)

	fmt.Printf(lib.RST + lib.BOLD + "BGP OPERATIONS\n" + lib.RST)
	fmt.Printf(lib.RST + " 5) Establish BGP Sessions                          \n" + lib.RST)
	fmt.Printf(lib.RST + " 6) Run BGP Route Advertisement\n" + lib.RST)

	fmt.Printf(lib.RST + lib.BOLD + "SYSTEM\n" + lib.RST)
	fmt.Printf(lib.RST + "10) Restore\n" + lib.RST)
	fmt.Printf(lib.RST + "11) Exit\n" + lib.RST)
	fmt.Printf(lib.CYN +
		"+=============================================================+\n" + lib.RST)
	fmt.Printf(lib.BOLD + "\nEnter choice: " + lib.RST)
}

func Run() {
	for {
		ShowMenu()

		choice, err := util.ReadInt()
		if err != nil {
			fmt.Printf(lib.RED + "  Invalid input\n" + lib.RST)
			util.PauseKey()
			continue
		}
		util.Clear()

		switch choice {
		case 1:
			fmt.Printf(lib.BOLD + "TOPOLOGY DIAGRAM\n\n" + lib.RST)
			util.PauseKey()
		case 2:
			fmt.Printf(lib.BOLD + "BGP STATUS\n\n" + lib.RST)
			util.PauseKey()
		case 3:
			fmt.Printf(lib.BOLD + "ROUTING TABLES\n\n" + lib.RST)
			util.PauseKey()
		case 4:
			// fmt.Printf(lib.BOLD + "EVENT LOGS\n\n" + lib.RST)
			util.PauseKey()
		case 5:
			fmt.Printf(lib.BOLD + "ESTABLISH BGP SESSIONS\n\n" + lib.RST)
			util.PauseKey()
		case 6:
			fmt.Printf(lib.BOLD + "ADVERTISING ROUTES\n\n" + lib.RST)
			util.PauseKey()
		case 10:
			// fmt.Printf(lib.BOLD + "RESTORE\n\n" + lib.RST)
			util.PauseKey()
		case 11:
			fmt.Printf(lib.CYN + "\n[+] Goodbye! Simulation terminated\n\n" + lib.RST)
			return
		default:
			fmt.Printf(lib.RED + "[-] Unknown option\n" + lib.RST)
			util.PauseKey()
		}
	}
}
