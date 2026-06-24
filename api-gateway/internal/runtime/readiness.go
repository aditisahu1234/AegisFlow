/*
The Health Engine only tells us whether each
component is healthy. The Readiness Engine combines
those health states with the dependency graph to
decide whether the application should receive
traffic. This separation of health and readiness
is a core cloud-native design principle.
*/
package runtime

import "context"

func (m *Manager) Ready(ctx context.Context) bool {

	for _, component := range m.components {

		deps := component.Dependencies()

		for _, dep := range deps {

			if dep.Type != Required {
				continue
			}

			if m.state.Get(dep.Name) != StateHealthy {
				return false
			}
		}
	}

	return true
}
func (m *Manager) RequiredHealthy() []string {

	var healthy []string

	for _, component := range m.components {

		for _, dep := range component.Dependencies() {

			if dep.Type != Required {
				continue
			}

			if m.state.Get(dep.Name) == StateHealthy {
				healthy = append(
					healthy,
					dep.Name,
				)
			}
		}
	}

	return healthy
}
