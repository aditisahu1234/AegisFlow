/*internal is nothing but the RUNTIME MANQGER.
The Runtime Manager owns every infrastructure
component in the process. Initially it just
registers components, but later it will also
manage retries, readiness, metrics, events,
graceful shutdown, and dependency ordering.*/

package runtime

import "sync"

type Manager struct {
	mu sync.RWMutex

	components map[string]Component
}

func NewManager() *Manager {
	return &Manager{
		components: make(map[string]Component),
	}
}

func (m *Manager) Register(c Component) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.components[c.Name()] = c
}
