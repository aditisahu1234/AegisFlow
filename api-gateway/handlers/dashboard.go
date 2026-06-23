package handlers

import (
	"fmt"
	"net/http"

	"api-gateway/middleware"
	"api-gateway/redis"
)

func DashboardHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set(
		"Content-Type",
		"text/html",
	)

	currentMode := "STANDBY"

	if !redis.RedisHealthy.Load() {
		currentMode = "ACTIVE"
	}

	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>

<head>

<title>AegisFlow Dashboard</title>

<style>

body {
	font-family: "Segoe UI", Arial, sans-serif;

	background: linear-gradient(
		135deg,
		#ffe4ec,
		#ffd6ea,
		#fff5fa
	);

	padding: 40px;
	color: #222;
}

h1 {
	text-align: center;
	font-size: 42px;
	font-weight: 800;
	color: #c2185b;
	margin-bottom: 10px;
}

.subtitle {
	text-align: center;
	font-size: 18px;
	color: #7a4b63;
	margin-bottom: 35px;
}

.card {

	background: rgba(
		255,
		255,
		255,
		0.96
	);

	padding: 24px;

	margin-bottom: 22px;

	border-radius: 22px;

	border-left: 8px solid #ec407a;

	box-shadow:
		0 8px 20px rgba(
			0,
			0,
			0,
			0.08
		);

	transition: all 0.2s ease;
}

.card:hover {

	transform: translateY(-4px);

	box-shadow:
		0 14px 28px rgba(
			0,
			0,
			0,
			0.12
		);
}

.section-title {

	font-size: 24px;
	font-weight: 700;

	color: #d81b60;

	margin-bottom: 18px;

	border-bottom: 2px solid #f8bbd0;

	padding-bottom: 8px;
}

.metric {

	display: flex;

	justify-content: space-between;

	padding: 10px 0;

	font-size: 18px;

	border-bottom: 1px solid #f4f4f4;
}

.metric:last-child {
	border-bottom: none;
}

.value {

	font-weight: 800;

	color: #3f7d20;

	font-size: 19px;
}

.footer {

	text-align: center;

	margin-top: 30px;

	color: #666;

	font-size: 14px;
}

</style>

</head>

<body>

<h1>AegisFlow Reliability Dashboard</h1>

<p class="subtitle">
Production Resilience & Distributed API Protection Platform
</p>

<div class="card">

	<div class="section-title">
		Traffic Metrics
	</div>

	<div class="metric">
		<span>Total Requests</span>
		<span class="value">%d</span>
	</div>

	<div class="metric">
		<span>Allowed Requests</span>
		<span class="value">%d</span>
	</div>

	<div class="metric">
		<span>Blocked Requests</span>
		<span class="value">%d</span>
	</div>

</div>

<div class="card">

	<div class="section-title">
		Circuit Breaker
	</div>

	<div class="metric">
		<span>Current State</span>
		<span class="value">%s</span>
	</div>

	<div class="metric">
		<span>Open Count</span>
		<span class="value">%d</span>
	</div>

	<div class="metric">
		<span>Half Open Count</span>
		<span class="value">%d</span>
	</div>

	<div class="metric">
		<span>Close Count</span>
		<span class="value">%d</span>
	</div>

	<div class="metric">
		<span>Total Requests</span>
		<span class="value">%d</span>
	</div>

	<div class="metric">
		<span>Failures</span>
		<span class="value">%d</span>
	</div>

	<div class="metric">
		<span>Successes</span>
		<span class="value">%d</span>
	</div>

</div>

<div class="card">

	<div class="section-title">
		Reliability Metrics
	</div>

	<div class="metric">
		<span>Timeout Events</span>
		<span class="value">%d</span>
	</div>

	<div class="metric">
		<span>Retries</span>
		<span class="value">%d</span>
	</div>

	<div class="metric">
		<span>Circuit Rejections</span>
		<span class="value">%d</span>
	</div>

</div>

<div class="card">

	<div class="section-title">
		Redis Health
	</div>

	<div class="metric">
		<span>Redis Failures</span>
		<span class="value">%d</span>
	</div>

	<div class="metric">
		<span>Redis Recoveries</span>
		<span class="value">%d</span>
	</div>

</div>

<div class="card">

	<div class="section-title">
		Fallback Limiter
	</div>

	<div class="metric">
		<span>Fallback Activations</span>
		<span class="value">%d</span>
	</div>

	<div class="metric">
		<span>Current Mode</span>
		<span class="value">%s</span>
	</div>

	<div class="metric">
		<span>Fallback Requests</span>
		<span class="value">%d</span>
	</div>

	<div class="metric">
		<span>Fallback Blocks</span>
		<span class="value">%d</span>
	</div>

</div>

<div class="footer">
	AegisFlow • Distributed Rate Limiting • Circuit Breaking • Reliability Engineering
</div>

</body>
</html>
`,
		middleware.GlobalMetrics.TotalRequests,
		middleware.GlobalMetrics.AllowedRequests,
		middleware.GlobalMetrics.BlockedRequests,

		redis.RedisBreaker.State().String(),

		middleware.GlobalCircuitMetrics.OpenCount.Load(),
		middleware.GlobalCircuitMetrics.HalfOpenCount.Load(),
		middleware.GlobalCircuitMetrics.CloseCount.Load(),

		middleware.GlobalCircuitMetrics.Requests.Load(),
		middleware.GlobalCircuitMetrics.Failures.Load(),
		middleware.GlobalCircuitMetrics.Successes.Load(),

		middleware.GlobalReliabilityMetrics.TimeoutCount.Load(),
		middleware.GlobalReliabilityMetrics.RetryCount.Load(),
		middleware.GlobalReliabilityMetrics.CircuitRejectedCount.Load(),

		middleware.GlobalReliabilityMetrics.RedisFailureCount.Load(),
		middleware.GlobalReliabilityMetrics.RedisRecoveryCount.Load(),

		middleware.GlobalReliabilityMetrics.FallbackActivations.Load(),
		currentMode,
		middleware.GlobalReliabilityMetrics.FallbackRequests.Load(),

		middleware.GlobalReliabilityMetrics.FallbackBlocks.Load(),
	)
}
