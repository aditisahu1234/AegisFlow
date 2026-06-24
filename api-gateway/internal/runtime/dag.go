// dependency DAG, level-order DAG execution.
package runtime

import "fmt"

type DAG struct {
	adjacency map[string][]string

	inDegree map[string]int
}

func NewDAG() *DAG {

	return &DAG{

		adjacency: make(map[string][]string),

		inDegree: make(map[string]int),
	}
}

func (m *Manager) BuildDAG() (*DAG, error) {

	dag := NewDAG()

	// initialize every component
	for name := range m.components {

		dag.inDegree[name] = 0
	}

	for component, deps := range m.dependencies {

		for _, dep := range deps {

			if dep.Type != Required {
				continue
			}

			if _, ok := m.components[dep.Name]; !ok {

				return nil, fmt.Errorf(
					"%s depends on unknown component %s",
					component,
					dep.Name,
				)
			}

			// dep → component
			dag.adjacency[dep.Name] =
				append(
					dag.adjacency[dep.Name],
					component,
				)

			dag.inDegree[component]++
		}
	}

	return dag, nil
}
