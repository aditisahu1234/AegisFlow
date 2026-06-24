/*internal is nothing but the RUNTIME MANQGER.
The Runtime Manager owns every infrastructure
component in the process. Initially it just
registers components, but later it will also
manage retries, readiness, metrics, events,
graceful shutdown, and dependency ordering.*/

package runtime

import (
	"context"
	"log"
	"sync"
)

type Manager struct {
	mu sync.RWMutex

	components map[string]Component

	dependencies map[string][]Dependency

	state *StateStore

	events *EventBus

	dag *DAG

	startupOrder []string

	metrics *RuntimeMetrics
}

func NewManager() *Manager {
	return &Manager{
		components:   make(map[string]Component),
		dependencies: make(map[string][]Dependency),
		state:        NewStateStore(),
		events:       NewEventBus(),
		metrics:      &RuntimeMetrics{},
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

	if err := m.BuildRuntime(); err != nil {
		return err
	}
	log.Println("Runtime startup order:")

	for i, name := range m.startupOrder {

		log.Printf(
			"%d. %s",
			i+1,
			name,
		)
	}

	//swquential startup: deterministic, simple, debuggable
	for _, name := range m.startupOrder {

		component := m.components[name]

		m.startComponent(
			ctx,
			component,
		)
	}
	return nil
}

// we compute everything once
func (m *Manager) BuildRuntime() error {

	dag, err := m.BuildDAG()
	if err != nil {
		return err
	}

	order, err := dag.TopologicalSort()
	if err != nil {
		return err
	}

	m.dag = dag

	m.startupOrder = order

	return nil

}
