/*
 * lib/state.go
 * Package level globals
 */
package lib

var (
	Network  []AutonomousSystem // live set of Autonomous Systems
	EventLog []LogEntry         // chronological simulation event log
	Tick     int                // monotonically-increasing simulation clock
)
