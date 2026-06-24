/*
When Kubernetes terminates a Pod, it sends a
SIGTERM and waits (usually 30 seconds) before
forcefully killing it. A production application
should stop accepting new work, gracefully shut
down all components in reverse dependency order,
and then exit cleanly. This prevents data loss,
broken requests, and inconsistent state
*/
package runtime

import (
	"context"
	"log"
	"time"
)

func (m *Manager) Stop(ctx context.Context) error {

	log.Println("[Runtime] beginning graceful shutdown")

	for i := len(m.startupOrder) - 1; i >= 0; i-- {

		name := m.startupOrder[i]

		component := m.components[name]

		m.events.Publish(Event{
			Time:      time.Now(),
			Component: name,
			Type:      EventStopping,
		})

		log.Printf(
			"[Runtime] stopping %s",
			name,
		)

		if err := component.Stop(ctx); err != nil {

			log.Printf(
				"[Runtime] failed stopping %s: %v",
				name,
				err,
			)

			continue
		}

		m.events.Publish(Event{
			Time:      time.Now(),
			Component: name,
			Type:      EventStopped,
		})

		m.state.Set(
			name,
			StateStopped,
		)
	}

	log.Println("[Runtime] graceful shutdown complete")

	return nil
}
