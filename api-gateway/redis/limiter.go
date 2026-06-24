package redis

import (
	"context"
	"log"
	"math/rand"
	"time" //sliding wind, time stamps'

	//RETRY+EXPO BACKOFF+JITTER
	"api-gateway/middleware"

	"github.com/cenkalti/backoff/v4" //failure resistance, CENKALTI

	"api-gateway/telemetry"
	"errors" //gobreaker errors

	"github.com/sony/gobreaker/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

// building redis sliding window sorted sets
func (l *RedisSlidingWindowLimiter) Allow(
	ctx context.Context,
	ip string,
) bool {

	ctx, span := //creaing child span for redis rate limiter
		telemetry.Tracer.Start(
			ctx,
			"redis_rate_limit_check",
		)

	defer span.End()

	telemetry.ActiveRequestsCounter.Add(
		ctx,
		1,
	)

	span.SetAttributes(
		attribute.String(
			"rate_limit.ip",
			ip,
		),

		attribute.Int(
			"rate_limit.limit",
			l.limit,
		),

		attribute.String(
			"limiter.type",
			"redis_sorted_set",
		),
	)

	//if redis works then:

	middleware.GlobalMetrics.TotalRequests++

	span.AddEvent(
		"redis_check_started",
	)

	key := "rate_limit:" + ip     //create redis key
	now := time.Now().UnixMilli() //precision---milliseconds

	windowStart := now - l.window.Milliseconds() //calc windowstart

	//CIRCUIT BREAKER

	var result int
	operation := func() error {

		//every redis oper. needs a context
		ctx, cancel := context.WithTimeout( //adding dependency timeout
			ctx, //every retry gets 100ms fresh timeout instead of sharing one
			100*time.Millisecond,
		)
		defer cancel()

		ctx, breakerSpan :=
			telemetry.Tracer.Start(
				ctx,
				"circuit_breaker",
			)

		defer breakerSpan.End()

		breakerSpan.SetAttributes(
			attribute.String(
				"breaker.name",
				"redis-breaker",
			),
		)

		breakerSpan.AddEvent(
			"breaker_execute_started",
		)

		if RedisBreaker == nil {
			log.Println("FATAL: RedisBreaker is nil")
			return errors.New("redis breaker is nil")
		}

		if SlidingWindowScript == "" {
			log.Println("FATAL: Lua script is empty")
			return errors.New("lua script not loaded")
		}
		log.Println(
			"RedisBreaker nil?",
			RedisBreaker == nil,
		)

		log.Println(
			"Redis Client nil?",
			Client == nil,
		)

		log.Println(
			"Lua Empty?",
			SlidingWindowScript == "",
		)
		_, err := RedisBreaker.Execute(
			func() (any, error) {

				redisStart := time.Now() //start timer

				r, err := Client.Eval(
					ctx,                 //lua cript , sony go breaker
					SlidingWindowScript, //failure recorded by go breaker
					[]string{key},
					now, //circuit ready to trip-->OPEN
					windowStart,
					l.limit,
				).Int()

				if err != nil { //redis fails

					//both redis rate limit check and circuit breaker show failure
					span.RecordError(err) //record error
					span.SetStatus(
						codes.Error,
						err.Error(),
					)

					breakerSpan.RecordError(err)

					breakerSpan.SetStatus(
						codes.Error,
						err.Error(),
					)

					span.AddEvent(
						"redis_operation_failed",
					)

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

						telemetry.TimeoutCounter.Add(
							ctx,
							1,
						)

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

				//REDIS SUCCEEDS
				middleware.GlobalCircuitMetrics.Successes.Add(1)
				middleware.GlobalCircuitMetrics.Requests.Add(1)

				breakerSpan.AddEvent(
					"breaker_closed",
				)
				telemetry.FallbackMode.Store(0)

				breakerSpan.SetAttributes(
					attribute.String(
						"breaker.state",
						"closed",
					),
				)

				redisDuration :=
					time.Since(redisStart) //end timer, record time

				telemetry.
					RedisLatencyHistogram.
					Record(
						ctx,
						redisDuration.Seconds(),
					)
				span.AddEvent(
					"redis_check_completed",
				)
				result = r

				span.SetAttributes(
					attribute.Int(
						"rate_limit.result",
						result,
					),
				)

				if result == 1 {

					span.AddEvent(
						"request_allowed",
					)

				} else {

					span.AddEvent(
						"request_blocked",
					)
				}

				return nil, nil
			},
		)

		//Handle Open Circuit, After retries fail, Gobreaker may start returning:
		if err != nil {

			if errors.Is(
				err,
				gobreaker.ErrOpenState,
			) {

				breakerSpan.AddEvent(
					"breaker_open",
				)

				breakerSpan.SetAttributes(
					attribute.String(
						"breaker.state",
						"open",
					),
				)

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
	b.RandomizationFactor = 0 //disable library jitter, use AWS Full Jitter manually

	b.Reset()

	//RETRY LOGIC
	for attempt := 0; attempt < 3; attempt++ {

		retryCtx, retrySpan := //retry span
			telemetry.Tracer.Start(
				ctx,
				"retry_attempt",
			)

		retrySpan.SetAttributes(
			attribute.Int(
				"retry.attempt",
				attempt,
			),
		)

		if attempt > 0 {

			middleware.
				GlobalReliabilityMetrics.
				RetryCount.
				Add(1)

			telemetry.RetryCounter.Add( //updating global telemetry
				retryCtx,
				1,
			)
		}

		err := operation()

		if err != nil { //RETRY FAILURE

			retrySpan.RecordError(err)

			retrySpan.SetAttributes(
				attribute.Bool(
					"retry.success",
					false,
				),
			)

			retrySpan.SetStatus(
				codes.Error,
				err.Error(),
			)

		}

		if err == nil { //RETRY SUCCEEDS
			retrySpan.AddEvent(
				"retry_succeeded",
			)
			retrySpan.SetAttributes(
				attribute.Bool(
					"retry.success",
					true,
				),
			)

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
			//acivate fallbakck limiter
			log.Println(
				"FALLBACK LIMITER ACTIVE",
			)

			telemetry.FallbackMode.Store(1)

			middleware.
				GlobalReliabilityMetrics.
				CircuitRejectedCount.
				Add(1)

			//store result
			span.AddEvent(
				"fallback_limiter_activated",
			)

			span.SetAttributes(
				attribute.String(
					"limiter.mode",
					"fallback",
				),
			)
			_, fallbackSpan :=
				telemetry.Tracer.Start(
					ctx,
					"fallback_limiter",
				)

			fallbackSpan.SetAttributes(
				attribute.String(
					"fallback.type",
					"in_memory_sliding_window",
				),

				attribute.String(
					"fallback.reason",
					"circuit_breaker_open",
				),
			)
			allowed := FallbackLimiter.Allow(
				ctx,
				ip,
			)

			fallbackSpan.AddEvent(
				"fallback_completed",
			)

			if !allowed {

				fallbackSpan.AddEvent(
					"fallback_request_blocked",
				)

				middleware.
					GlobalReliabilityMetrics.
					FallbackBlocks.
					Add(1)

			} else {

				fallbackSpan.AddEvent(
					"fallback_request_allowed",
				)
			}

			fallbackSpan.SetAttributes(
				attribute.Bool(
					"fallback.allowed",
					allowed,
				),
			)

			middleware.
				GlobalReliabilityMetrics.
				FallbackRequests.
				Add(1)

			retrySpan.End()
			fallbackSpan.End()

			return allowed
		}

		retrySpan.AddEvent(
			"backoff_started",
		)
		//BACKOFF STARTED
		expoDelay := b.NextBackOff() //using cenkalti

		jitterDelay :=
			time.Duration(
				rand.Int63n( //using full jitter AWS, implemented on my own
					int64(expoDelay),
				),
			)

		telemetry. //delay histogram
				RetryDelayHistogram.
				Record(
				retryCtx, //using trace context
				jitterDelay.Seconds(),
			)

		time.Sleep(jitterDelay) //BACKOFF FINISHED
		retrySpan.AddEvent(
			"backoff_finished",
		)
		retrySpan.End()

	}

	/* sleep =
	random(
		0,
		min(cap, base * 2^attempt)
	)*/
	return false

}
