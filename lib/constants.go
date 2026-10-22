/*
 * lib/constants.go
 * Simulation-wide numeric limits.
 */
package lib

const (
	MaxPath = 10  // maximum AS hops before a packet is dropped (TTL)
	LogCap  = 200 // event log capacity; oldest entry discarded when full
)
