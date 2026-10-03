/*
 * util/input.go
 * Terminal I/O helpers stuff
 */
package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// Shared buffered reader for all interactive input.
var stdinReader = bufio.NewReader(os.Stdin)

// Reads one line from stdin and parses it as an integer
func ReadInt() (int, error) {
	line, _ := stdinReader.ReadString('\n')
	return strconv.Atoi(strings.TrimSpace(line))
}

// Reads and returns one trimmed line from stdin
func ReadString() string {
	line, _ := stdinReader.ReadString('\n')
	return strings.TrimSpace(line)
}

// Prints a prompt and blocks until the user presses Enter
func PauseKey() {
	fmt.Print("\033[2m\n[Press ENTER to continue...]\033[0m")
	stdinReader.ReadString('\n') //nolint:errcheck
}

// Clears the terminal screen
func Clear() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
