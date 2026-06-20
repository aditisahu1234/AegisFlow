package middleware

import "sync/atomic"
//adding all global matrics here
type ReliabilityMetrics struct {

	TimeoutCount atomic.Uint64

	RetryCount atomic.Uint64

	CircuitRejectedCount atomic.Uint64

	RedisFailureCount atomic.Uint64

	RedisRecoveryCount atomic.Uint64

	FallbackActivations atomic.Uint64

	FallbackRequests atomic.Uint64

	FallbackBlocks atomic.Uint64
}

var GlobalReliabilityMetrics ReliabilityMetrics