//telemetry state variables

package telemetry

import "sync/atomic"

var RedisHealth atomic.Int64

var CircuitState atomic.Int64

var FallbackMode atomic.Int64

var ActiveIPs atomic.Int64
