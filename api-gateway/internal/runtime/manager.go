/*internal is nothing but the RUNTIME MANQGER.
The Runtime Manager owns every infrastructure
component in the process. Initially it just
registers components, but later it will also
manage retries, readiness, metrics, events,
graceful shutdown, and dependency ordering.*/

package runtime

import (
	"context"
	"sync"
)

type Manager struct {
	mu sync.RWMutex

	components map[string]Component

	dependencies map[string][]Dependency

	state *StateStore

	events *EventBus
}

func NewManager() *Manager {
	return &Manager{
		components:   make(map[string]Component),
		dependencies: make(map[string][]Dependency),
		state:        NewStateStore(),
		events:       NewEventBus(),
	}
}

func (m *Manager) Register(c Component) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.components[c.Name()] = c

	m.dependencies[c.Name()] = c.Dependencies()
}

func (m *Manager) Start(ctx context.Context) error {
	//bad runtime configurations die immediately
	if err := m.validateDependencies(); err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, component := range m.components {

		wg.Add(1)

		go func(c Component) {
			defer wg.Done()

			m.startComponent(
				ctx,
				c,
			)

		}(component)
	}

	wg.Wait()
	return nil
}
