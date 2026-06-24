package runtime

import (
	"context"
	"errors"
)

func (m *Manager) StartComponent(
	ctx context.Context,
	name string,
	visited map[string]bool,
) error {

	if visited[name] {
		return nil
	}

	component, ok := m.components[name]
	if !ok {
		return errors.New("component not registered: " + name)
	}

	visited[name] = true

	for _, dep := range m.dependencies[name] {

		if dep.Type != Required {
			continue
		}

		if err := m.StartComponent(
			ctx,
			dep.Name,
			visited,
		); err != nil {

			return err
		}
	}

	m.startComponent(
		ctx,
		component,
	)

	return nil
}
