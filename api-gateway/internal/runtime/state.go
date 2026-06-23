/*
This file defines the lifecycle states a runtime
component can be in. Instead of using booleans
like connected or healthy, we model a real state
machine, which makes monitoring, retries, and
debugging much cleaner.*/

package runtime

type State string

const (
	StateStopped   State = "stopped"
	StateStarting  State = "starting"
	StateHealthy   State = "healthy"
	StateDegraded  State = "degraded"
	StateUnhealthy State = "unhealthy"
	StateStopping  State = "stopping"
)
