package redis

import (
	"context"
	"time"		//sliding wind, time stamps'
	"math/rand"
	//RETRY+EXPO BACKOFF+JITTER
	"github.com/cenkalti/backoff/v4"		//failure resistance, CENKALTI
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

		r, err := Client.Eval(
			ctx,
			SlidingWindowScript,		//lua script
			[]string{key},
			now,
			windowStart,
			l.limit,
		).Int()
	
		if err != nil {
			return err
		}
	
		result = r
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

	for attempt := 0; attempt < 3; attempt++ {

		err := operation()

		if err == nil {
			return result == 1
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