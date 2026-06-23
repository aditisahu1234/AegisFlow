package runtime

import "fmt"

func (m *Manager) validateDependencies() error {

	for component, deps := range m.dependencies {

		for _, dep := range deps {

			if _, exists := m.components[dep.Name]; !exists {

				if dep.Type == Required {

					return fmt.Errorf(
						"%s depends on unknown component %s",
						component,
						dep.Name,
					)
				}
			}
		}
	}

	return nil
}
