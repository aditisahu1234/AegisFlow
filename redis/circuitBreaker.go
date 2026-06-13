//implementing sony circuitbreaker here
package redis

import (
	"log"
	"time"
	"api-gateway/middleware"
	"github.com/sony/gobreaker/v2"
)

var RedisBreaker *gobreaker.CircuitBreaker[any]

func InitCircuitBreaker() {

	settings := gobreaker.Settings{

		Name: "redis-breaker",

		// Half-open state, allow 3 reqs while half open, check if success or not  
		MaxRequests: 3,

		// Rolling window		look at last 10 seconds
		Interval: 10 * time.Second,

		// 1 second buckets     10 buckets of 1s each
		BucketPeriod: 1 * time.Second,

		// Stay open for 5 seconds
		Timeout: 5 * time.Second,

		ReadyToTrip: func(
			counts gobreaker.Counts,
		) bool {

			// Need enough traffic first
			if counts.Requests < 3 {
				return false
			}
			log.Printf(
				"Requests=%d Failures=%d",
				counts.Requests,
				counts.TotalFailures,
			)

			failureRate :=
				float64(counts.TotalFailures) /
					float64(counts.Requests)

			return failureRate >= 0.5		//error treshold, open circuit after this much error percentage
		},

		OnStateChange: func(
			name string,
			from gobreaker.State,
			to gobreaker.State,
		) {

			log.Printf(
				"[CircuitBreaker] %s: %s -> %s",
				name,
				from.String(),
				to.String(),
			)
			//add metrics, every transition is recorded from circuit_metrics.go
			switch to {

			case gobreaker.StateOpen:
				middleware.GlobalCircuitMetrics.
					OpenCount.Add(1)
				middleware.GlobalReliabilityMetrics.
					FallbackActivations.Add(1)
		
			case gobreaker.StateHalfOpen:
				middleware.GlobalCircuitMetrics.
					HalfOpenCount.Add(1)
		
			case gobreaker.StateClosed:
				middleware.GlobalCircuitMetrics.
					CloseCount.Add(1)
			}
		},
	}
/*creates sony's internal state machine, internally it tracks:
type Counts struct {
    Requests
    TotalSuccesses
    TotalFailures
    ConsecutiveSuccesses
    ConsecutiveFailures
}
*/
	RedisBreaker =
		gobreaker.NewCircuitBreaker[any](		//sony gobreaker used
			settings,
		)
}