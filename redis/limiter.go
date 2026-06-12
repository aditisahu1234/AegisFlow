package redis

import (
	"context"
	"time"		//sliding wind, time stamps'
	"math/rand"
	"log"
	//RETRY+EXPO BACKOFF+JITTER
	"api-gateway/middleware"
	"github.com/cenkalti/backoff/v4"		//failure resistance, CENKALTI
	
	"errors"							//gobreaker errors
	"github.com/sony/gobreaker/v2"
)
type RedisSlidingWindowLimiter struct {
	limit  int
	window time.Duration
}

func NewRedisSlidingWindowLimiter(
	limit int,
	window time.Duration,
) *RedisSlidingWindowLimiter {

	return &RedisSlidingWindowLimiter{
		limit:  limit,
		window: window,
	}
}

//building redis sliding window sorted sets
func (l *RedisSlidingWindowLimiter) Allow(
	ip string,
) bool {
	middleware.GlobalMetrics.TotalRequests++
	key := "rate_limit:" + ip	//create redis key
	now := time.Now().UnixMilli()		//precision---milliseconds
	
	windowStart := now - l.window.Milliseconds()	//calc windowstart


	var result int
	operation := func() error {

		//every redis oper. needs a context
		ctx, cancel := context.WithTimeout(		//adding dependency timeout
			context.Background(),		//every retry gets 100ms fresh timeout instead of sharing one
			100*time.Millisecond,
		)
		defer cancel()

		_, err := RedisBreaker.Execute(
			func() (any, error) {
	
				r, err := Client.Eval(
					ctx,					//lua cript , sony go breaker
					SlidingWindowScript,		//failure recorded by go breaker
					[]string{key},
					now,				//circuit ready to trip-->OPEN
					windowStart,
					l.limit,
				).Int()
		
				if err != nil {

					middleware.GlobalCircuitMetrics.Failures.Add(1)
					middleware.GlobalCircuitMetrics.Requests.Add(1)
				
					if errors.Is(
						err,
						context.DeadlineExceeded,
					) {
				
						middleware.
							GlobalReliabilityMetrics.
							TimeoutCount.
							Add(1)
				
						log.Println(
							"Dependency timeout recorded",
						)
					}
				
					log.Println(
						"Breaker saw error:",
						err,
					)
				
					return nil, err
				}
			
				middleware.GlobalCircuitMetrics.Successes.Add(1)
				middleware.GlobalCircuitMetrics.Requests.Add(1)
				result = r
		
				return nil, nil
			},
		)
		
		//Handle Open Circuit, After retries fail, Gobreaker may start returning:
		if err != nil {

			if errors.Is(
				err,
				gobreaker.ErrOpenState,
			) {
			
				middleware.
					GlobalReliabilityMetrics.
					CircuitRejectedCount.
					Add(1)
			
				log.Println(
					"Circuit OPEN - failing fast",
				)
			
				return err
			}
		
			if errors.Is(
				err,
				gobreaker.ErrTooManyRequests,
			) {
				return err
			}
		
			return err
		}
		return nil
	}

	//using the cenkalti lib, and adding aws exponential backoff
	b := backoff.NewExponentialBackOff()

	b.InitialInterval = 50 * time.Millisecond
	b.Multiplier = 2
	b.MaxInterval = 1 * time.Second
	b.MaxElapsedTime = 3 * time.Second
	b.RandomizationFactor = 0	//disable library jitter, use AWS Full Jitter manually

	b.Reset()

	//RETRY LOGIC
	for attempt := 0; attempt < 3; attempt++ {

		if attempt > 0 {

			middleware.
				GlobalReliabilityMetrics.
				RetryCount.
				Add(1)
		}
		err := operation()
		
		if err == nil {

			if result == 1 {
		
				middleware.GlobalMetrics.
					AllowedRequests++
		
			} else {
		
				middleware.GlobalMetrics.
					BlockedRequests++
			}
		
			return result == 1
		}
		//CIRCUIT REJECTION
		if errors.Is(
			err,
			gobreaker.ErrOpenState,
		) {
		
			log.Println(
				"Skipping retries because breaker is OPEN",
			)
		
			return false
		}

		expoDelay := b.NextBackOff()		//using cenkalti

		jitterDelay :=
			time.Duration(
				rand.Int63n(		//using full jitter AWS, implemented on my own
					int64(expoDelay),
				),
			)

		time.Sleep(jitterDelay)
	}

	/* sleep =		
	random(
		0,
		min(cap, base * 2^attempt)
	)*/
	return false

}