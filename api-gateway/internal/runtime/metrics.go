package runtime

import "sync/atomic"

type RuntimeMetrics struct {
	HealthyComponents atomic.Int64

	UnhealthyComponents atomic.Int64

	StartupFailures atomic.Int64

	Recoveries atomic.Int64
}
