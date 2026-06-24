/*
Right now, components can tell us if they're
healthy, but nothing periodically checks them.
The Health Engine is a background service that
continuously polls every component and updates the
runtime state. Kubernetes readiness and liveness
probes will later read this information instead
of checking components directly.
*/

package runtime

import (
	"context"
	"log"
	"time"
)

const healthCheckInterval = 10 * time.Second

func (m *Manager) StartHealthEngine(ctx context.Context) {

	ticker := time.NewTicker(healthCheckInterval)

	go func() {

		defer ticker.Stop()

		for {

			select {

			case <-ctx.Done():
				return

			case <-ticker.C:
				m.runHealthChecks(ctx)
			}
		}
	}()
}

func (m *Manager) runHealthChecks(ctx context.Context) {

	for _, name := range m.startupOrder {

		component := m.components[name]

		err := component.Health(ctx)

		if err != nil {

			log.Printf(
				"[Health] %s unhealthy: %v",
				component.Name(),
				err,
			)

			if m.state.Get(component.Name()) != StateUnhealthy {

				m.metrics.UnhealthyComponents.Add(1)
			}

			m.state.Set(
				component.Name(),
				StateUnhealthy,
			)

			m.events.Publish(Event{
				Time:      time.Now(),
				Component: component.Name(),
				Type:      EventUnhealthy,
				Err:       err,
			})

			continue
		}

		if m.state.Get(component.Name()) != StateHealthy {

			m.metrics.HealthyComponents.Add(1)
		}

		m.state.Set(
			component.Name(),
			StateHealthy,
		)

		m.events.Publish(Event{
			Time:      time.Now(),
			Component: component.Name(),
			Type:      EventHealthy,
		})
	}
}
