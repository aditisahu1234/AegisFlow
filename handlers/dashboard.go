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

	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
<title>AegisFlow Dashboard</title>

<style>

body {
	font-family: Arial, sans-serif;
	background: #FFE6F0;
	padding: 40px;
	color: #333;
}

h1 {
	text-align: center;
	font-size: 38px;
	font-weight: bold;
	color: #C2185B;
}

.card {
	background: white;
	padding: 20px;
	margin-bottom: 20px;
	border-radius: 20px;
	box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.section-title {
	font-size: 24px;
	font-weight: bold;
	color: #D81B60;
	margin-bottom: 15px;
}

.metric {
	padding: 6px 0;
	font-size: 18px;
}

.value {
	font-weight: bold;
	color: #AD1457;
}

</style>

</head>

<body>

<h1>🌸 AegisFlow Reliability Dashboard 🌸</h1>

<div class="card">
	<div class="section-title">Traffic</div>

	<div class="metric">Total Requests:
		<span class="value">%d</span>
	</div>

	<div class="metric">Allowed Requests:
		<span class="value">%d</span>
	</div>

	<div class="metric">Blocked Requests:
		<span class="value">%d</span>
	</div>
</div>

<div class="card">
	<div class="section-title">Circuit Breaker</div>

	<div class="metric">Current State:
		<span class="value">%s</span>
	</div>

	<div class="metric">Open Count:
		<span class="value">%d</span>
	</div>

	<div class="metric">Half Open Count:
		<span class="value">%d</span>
	</div>

	<div class="metric">Close Count:
		<span class="value">%d</span>
	</div>

	<div class="metric">Requests:
		<span class="value">%d</span>
	</div>

	<div class="metric">Failures:
		<span class="value">%d</span>
	</div>

	<div class="metric">Successes:
		<span class="value">%d</span>
	</div>
</div>

<div class="card">
	<div class="section-title">Reliability</div>

	<div class="metric">Timeouts:
		<span class="value">%d</span>
	</div>

	<div class="metric">Retries:
		<span class="value">%d</span>
	</div>

	<div class="metric">Circuit Rejections:
		<span class="value">%d</span>
	</div>
</div>

<div class="card">
	<div class="section-title">Redis Health</div>

	<div class="metric">Failures:
		<span class="value">%d</span>
	</div>

	<div class="metric">Recoveries:
		<span class="value">%d</span>
	</div>
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
	)
}