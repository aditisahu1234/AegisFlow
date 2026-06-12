package middleware

import "sync/atomic"

type ReliabilityMetrics struct {

	TimeoutCount atomic.Uint64

	RetryCount atomic.Uint64

	CircuitRejectedCount atomic.Uint64

	RedisFailureCount atomic.Uint64

	RedisRecoveryCount atomic.Uint64
}

var GlobalReliabilityMetrics ReliabilityMetrics