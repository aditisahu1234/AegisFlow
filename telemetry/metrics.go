//we create all instruments here

package telemetry
import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/metric"
)

var RequestCounter metric.Int64Counter

var RetryCounter metric.Int64Counter

var TimeoutCounter metric.Int64Counter

var FallbackCounter metric.Int64Counter

//for histogram
var RequestDurationHistogram metric.Float64Histogram	

var RedisLatencyHistogram metric.Float64Histogram

var RetryDelayHistogram metric.Float64Histogram

//up down counters
var ActiveRequestsCounter metric.Int64UpDownCounter

//observable counter to get goroutine count
var GoroutineCount metric.Int64ObservableUpDownCounter

//memory usage observab counter and gauge
var MemoryAlloc metric.Int64ObservableGauge

var GCCycles metric.Int64ObservableCounter

var RedisHealthGauge metric.Int64ObservableGauge

var CircuitStateGauge metric.Int64ObservableGauge

var FallbackModeGauge metric.Int64ObservableGauge

var ActiveIPsGauge metric.Int64ObservableGauge

//adding metrics from docs
func InitMetrics() error {

	var err error

	RequestCounter, err =
		Meter.Int64Counter(
			"aegisflow_requests_total",

			metric.WithDescription(
				"Total API requests processed by AegisFlow",
			),

			metric.WithUnit(
				"{request}",
			),
		)
	
	if err != nil {
		return err
	}
	
	RequestDurationHistogram, err =
		Meter.Float64Histogram(
			"aegisflow_request_duration_seconds",

			metric.WithDescription(
				"Distribution of API request durations",
			),

			metric.WithUnit(
				"s",
			),
		)

	if err != nil {
		return err
	}

	RetryDelayHistogram, err =
	Meter.Float64Histogram(
		"aegisflow_retry_delay_seconds",

		metric.WithDescription(
			"Distribution of retry delays",
		),

		metric.WithUnit("s"),
	)

	if err != nil {
		return err
	}
	RedisLatencyHistogram, err =
		Meter.Float64Histogram(
			"aegisflow_redis_latency_seconds",

			metric.WithDescription(
				"Distribution of Redis operation latency",
			),

			metric.WithUnit(
				"s",
			),
	)

	if err != nil {
		return err
	}

	RetryCounter, err =
		Meter.Int64Counter(
			"aegisflow_retries_total",

			metric.WithDescription(
				"Total retry attempts triggered by resilience layer",
			),

			metric.WithUnit(
				"{retry}",
			),
		)

	if err != nil {
		return err
	}

	TimeoutCounter, err =
		Meter.Int64Counter(
			"aegisflow_timeouts_total",

			metric.WithDescription(
				"Total dependency timeout events",
			),

			metric.WithUnit(
				"{timeout}",
			),
		)

	if err != nil {
		return err
	}

	FallbackCounter, err =
		Meter.Int64Counter(
			"aegisflow_fallback_activations_total",

			metric.WithDescription(
				"Total fallback limiter activations",
			),

			metric.WithUnit(
				"{activation}",
			),
		)

	if err != nil {
		return err

	}
	//up down counter
	ActiveRequestsCounter, err =
	Meter.Int64UpDownCounter(
		"aegisflow_active_requests",

		metric.WithDescription(
			"Current number of active requests",
		),

		metric.WithUnit(
			"{request}",
		),
	)

	if err != nil {
		return err
	}

	GoroutineCount, err =
	Meter.Int64ObservableUpDownCounter(
		"aegisflow_goroutines",

		metric.WithDescription(
			"Current number of goroutines",
		),

		metric.WithUnit(
			"{goroutine}",
		),
	)

	if err != nil {
		return err
	}

	_, err =
	Meter.RegisterCallback(		//SDK asks for current value, callback updates

		func(
			ctx context.Context,
			o metric.Observer,
		) error {

			o.ObserveInt64(
				GoroutineCount,

				int64(
					runtime.NumGoroutine(),
				),
			)

			return nil
		},

		GoroutineCount,
	)

	if err != nil {
		return err
	}
	//memory gauge
	MemoryAlloc, err =
	Meter.Int64ObservableGauge(
		"aegisflow_memory_alloc_bytes",

		metric.WithDescription(
			"Currently allocated memory",
		),

		metric.WithUnit(
			"bytes",
		),
	)

	if err != nil {
		return err
	}
	//gc counter
	GCCycles, err =
	Meter.Int64ObservableCounter(
		"aegisflow_gc_cycles_total",

		metric.WithDescription(
			"Total GC cycles",
		),
	)

	if err != nil {
		return err
	}
// ---------------------------------------------------
// SYSTEM STATE GAUGES
// ----------------------------------------------------
	RedisHealthGauge, err =		//redis health gauge
		Meter.Int64ObservableGauge(
			"redis.health",
			metric.WithDescription(
				"Redis health status",
			),
		)
	if err != nil {
		return err
	}
	CircuitStateGauge, err =	//circuit state gauge
		Meter.Int64ObservableGauge(
			"circuit.state",
			metric.WithDescription(
				"Circuit breaker state",
			),
		)
	if err != nil {
		return err
	}
		
	FallbackModeGauge, err =	//fallback mode gauge
		Meter.Int64ObservableGauge(
			"fallback.mode",
			metric.WithDescription(
				"Fallback limiter mode",
			),
		)
	if err != nil {
		return err
	}
	
	ActiveIPsGauge, err =		//active ips gauge
		Meter.Int64ObservableGauge(
			"fallback.active_ips",
			metric.WithDescription(
				"Active IPs in fallback limiter",
			),
		)
	if err != nil {
		return err
	}
		

	//register callback

	_, err = Meter.RegisterCallback(
		func(
			ctx context.Context,
			o metric.Observer,
		) error {

			o.ObserveInt64(
				RedisHealthGauge,
				RedisHealth.Load(),
			)
	
			o.ObserveInt64(
				CircuitStateGauge,
				CircuitState.Load(),
			)
	
			o.ObserveInt64(
				FallbackModeGauge,
				FallbackMode.Load(),
			)
	
			o.ObserveInt64(
				ActiveIPsGauge,
				ActiveIPs.Load(),
			)
	
			return nil
		},
		RedisHealthGauge,
		CircuitStateGauge,
		FallbackModeGauge,
		ActiveIPsGauge,
	)
		
	if err != nil {
		return err
	}

	_, err =
	Meter.RegisterCallback(

		func(
			ctx context.Context,
			o metric.Observer,
		) error {

			var mem runtime.MemStats

			runtime.ReadMemStats(
				&mem,
			)

			o.ObserveInt64(
				MemoryAlloc,
				int64(mem.Alloc),
			)

			o.ObserveInt64(
				GCCycles,
				int64(mem.NumGC),
			)

			return nil
		},

		MemoryAlloc,
		GCCycles,
	)

	if err != nil {
		return err
	}

	return nil
}