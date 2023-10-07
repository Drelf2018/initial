package graph_test

import (
	"fmt"
	"testing"

	"github.com/Drelf2018/initial/graph"
)

func TestGraph(t *testing.T) {
	g := make(graph.Graph[int, int])
	g.Add(3, 7, 1)
	g.Add(3, 6, 2)
	g.Add(2, 5, 3)
	g.Add(2, 4, 4)
	g.Add(1, 3, 5)
	g.Add(1, 2, 6)
	g.BFS(func(from, to *graph.Node[int, int], edge *graph.Edge[int, int]) {
		fmt.Printf("from: %v to: %v edge: %v\n", from.Index, to.Index, edge.Value)
		idx := 0
		if from.Parent != nil {
			idx = from.Parent.Index
		}
		fmt.Printf("  from.parent: %v to.parent: %v\n\n", idx, to.Parent.Index)
	})
}
