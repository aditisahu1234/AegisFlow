/*Instead of every component implementing its own
retry loop, the Runtime Manager will own a single
retry engine. Every component gets production-grade
retries with exponential backoff and jitter automatically,
eliminating duplicated code and ensuring consistent
behavior across the platform.*/

package runtime

import (
	"context"
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	baseDelay = 100 * time.Millisecond
	maxDelay  = 30 * time.Second
)

func fullJitter(attempt int) time.Duration {

	maxSleep := float64(baseDelay) * math.Pow(2, float64(attempt))

	if maxSleep > float64(maxDelay) {
		maxSleep = float64(maxDelay)
	}

	return time.Duration(
		rand.Int63n(int64(maxSleep)),
	)
}

func (m *Manager) startComponent(
	ctx context.Context,
	component Component,
) {

	m.state.Set(
		component.Name(),
		StateStarting,
	)

	m.events.Publish(Event{
		Time:      time.Now(),
		Component: component.Name(),
		Type:      EventStarting,
	})
	attempt := 0

	for {

		err := component.Start(ctx)

		if err == nil {

			log.Printf(
				"[Runtime] %s started successfully",
				component.Name(),
			)

			m.state.Set(
				component.Name(),
				StateHealthy,
			)

			m.events.Publish(Event{
				Time:      time.Now(),
				Component: component.Name(),
				Type:      EventHealthy,
			})

			return
		}

		m.state.Set(
			component.Name(),
			StateUnhealthy,
		)

		delay := fullJitter(attempt)

		log.Printf(
			"[Runtime] %s failed: %v",
			component.Name(),
			err,
		)

		m.events.Publish(Event{
			Time:      time.Now(),
			Component: component.Name(),
			Type:      EventFailed,
			Err:       err,
		})

		log.Printf(
			"[Runtime] retrying %s in %v",
			component.Name(),
			delay,
		)

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
