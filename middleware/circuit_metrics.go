package middleware

import "sync/atomic"

type CircuitMetrics struct {
	OpenCount       atomic.Uint64
	HalfOpenCount   atomic.Uint64
	CloseCount      atomic.Uint64

	Requests        atomic.Uint64
	Failures        atomic.Uint64
	Successes       atomic.Uint64
}

var GlobalCircuitMetrics CircuitMetrics