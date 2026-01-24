package controlflow

import "github.com/frederik-jatzkowski/havel/pkg/hvil/pass/optimization/statistics"

func discoverBlocks(entry Node) (order []Node) {
	visited := make(map[statistics.BlockID]bool)

	var visit func(Node)
	visit = func(n Node) {
		if n == nil || visited[n.ID()] {
			return
		}

		visited[n.ID()] = true

		// Post-order discovery (for faster backward analysis)
		for _, succ := range n.Successors() {
			visit(succ)
		}

		order = append(order, n)
	}

	visit(entry)

	return order
}
