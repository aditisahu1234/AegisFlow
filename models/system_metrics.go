package models

type SystemMetrics struct {

	TotalRequests   int `json:"total_requests"`
	AllowedRequests int `json:"allowed_requests"`
	BlockedRequests int `json:"blocked_requests"`

	CircuitOpenCount     uint64 `json:"circuit_open_count"`
	CircuitHalfOpenCount uint64 `json:"circuit_half_open_count"`
	CircuitCloseCount    uint64 `json:"circuit_close_count"`

	CircuitRequests uint64 `json:"circuit_requests"`
	CircuitFailures uint64 `json:"circuit_failures"`
	CircuitSuccesses uint64 `json:"circuit_successes"`

	TimeoutCount uint64 `json:"timeout_count"`
	RetryCount uint64 `json:"retry_count"`

	CircuitRejectedCount uint64 `json:"circuit_rejected_count"`

	RedisFailureCount uint64 `json:"redis_failure_count"`
	RedisRecoveryCount uint64 `json:"redis_recovery_count"`
}