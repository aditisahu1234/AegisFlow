package handlers

import (
	"encoding/json"
	"net/http"
	"api-gateway/redis"
	"api-gateway/middleware"
)

func MetricsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	middleware.GlobalMetrics.CurrentCircuitState =
	redis.RedisBreaker.
		State().
		String()

	middleware.GlobalMetrics.CircuitOpenCount =
		middleware.GlobalCircuitMetrics.
			OpenCount.Load()

	middleware.GlobalMetrics.CircuitHalfOpenCount =
		middleware.GlobalCircuitMetrics.
			HalfOpenCount.Load()

	middleware.GlobalMetrics.CircuitCloseCount =
		middleware.GlobalCircuitMetrics.
			CloseCount.Load()

	middleware.GlobalMetrics.CircuitRequests =
		middleware.GlobalCircuitMetrics.
			Requests.Load()

	middleware.GlobalMetrics.CircuitFailures =
		middleware.GlobalCircuitMetrics.
			Failures.Load()

	middleware.GlobalMetrics.CircuitSuccesses =
		middleware.GlobalCircuitMetrics.
			Successes.Load()

	middleware.GlobalMetrics.TimeoutCount =
		middleware.GlobalReliabilityMetrics.
			TimeoutCount.Load()

	middleware.GlobalMetrics.RetryCount =
		middleware.GlobalReliabilityMetrics.
			RetryCount.Load()

	middleware.GlobalMetrics.CircuitRejectedCount =
		middleware.GlobalReliabilityMetrics.
			CircuitRejectedCount.Load()

	middleware.GlobalMetrics.RedisFailureCount =
		middleware.GlobalReliabilityMetrics.
			RedisFailureCount.Load()

	middleware.GlobalMetrics.RedisRecoveryCount =
		middleware.GlobalReliabilityMetrics.
			RedisRecoveryCount.Load()

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	encoder.Encode(
		middleware.GlobalMetrics,
	)
}