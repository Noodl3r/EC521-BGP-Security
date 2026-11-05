/*
 * lib/state.go
 * Package-level globals shared across the entire simulation.
 * All mutable state lives here so the read/write surface is explicit.
 */
package lib

var (
	Network  []AutonomousSystem // live set of Autonomous Systems
	EventLog []LogEntry         // chronological simulation event log
	Tick     int                // monotonically-increasing simulation clock
)
