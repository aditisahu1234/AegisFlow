//implementing khan's algorithm, writes a startuporder follows it

package runtime

import "fmt"

func (d *DAG) TopologicalSort() ([]string, error) {

	queue := []string{}

	for node, degree := range d.inDegree {

		if degree == 0 {

			queue = append(
				queue,
				node,
			)
		}
	}

	order := []string{}

	for len(queue) > 0 {

		node := queue[0]

		queue = queue[1:]

		order = append(
			order,
			node,
		)

		for _, neighbor := range d.adjacency[node] {

			d.inDegree[neighbor]--

			if d.inDegree[neighbor] == 0 {

				queue = append(
					queue,
					neighbor,
				)
			}
		}
	}

	if len(order) != len(d.inDegree) {

		return nil, fmt.Errorf(
			"dependency cycle detected",
		)
	}

	return order, nil
}
