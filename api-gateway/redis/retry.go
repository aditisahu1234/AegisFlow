// will contain all startup retry logic for edis
package redis

import (
	"context"
	"log"
	"math"
	"math/rand"
	"time"

	"api-gateway/config"
)

const (
	baseDelay = 100 * time.Millisecond //exponential growth prevents hammering Redis during outages.
	maxDelay  = 30 * time.Second
)

// make it configurable
type RetryConfig struct {
	BaseDelay  time.Duration
	MaxDelay   time.Duration
	MaxRetries int // -1 means forever
}

// The reconnect attempts are spread across time, dramatically reducing load on Redis.
func fullJitter(attempt int) time.Duration {

	maxSleep := float64(baseDelay) * math.Pow(2, float64(attempt))

	if maxSleep > float64(maxDelay) {
		maxSleep = float64(maxDelay)
	}

	return time.Duration(
		rand.Int63n(int64(maxSleep)),
	)
}
func ConnectWithRetry(
	ctx context.Context,
	cfg config.Config,
) {

	attempt := 0

	for {

		err := ConnectRedis(cfg)

		if err == nil {
			log.Printf(
				"Connected to Redis after %d attempt(s)",
				attempt+1,
			)

			return
		}

		log.Printf(
			"[Redis Retry %d] connection failed: %v",
			attempt+1,
			err,
		)

		delay := fullJitter(attempt)

		log.Printf(
			"[Redis Retry %d] sleeping for %v",
			attempt+1,
			delay,
		)

		//because production services must be able to stop immediately.
		select {

		case <-ctx.Done():
			return

		case <-time.After(delay):
		}

		if attempt < 10 {
			attempt++
		}

	}
}
