package initial

import "github.com/Drelf2018/TypeGo/Queue"

type edge struct {
	Target int
	Next   int
}

type graph[V any] struct {
	Roots  map[int]bool
	Nodes  map[int]int
	Values map[int]V
	Edges  []edge
}

func (g *graph[V]) Add(from, to int) {
	if _, ok := g.Roots[from]; !ok {
		g.Roots[from] = true
	}
	g.Roots[to] = false
	g.Edges = append(g.Edges, edge{
		Target: to,
		Next:   g.Nodes[from],
	})
	g.Nodes[from] = len(g.Edges) - 1
}

func (g *graph[V]) BFS(do func(from, to int, value V)) {
	roots := make([]int, 0, len(g.Roots))
	for r, ok := range g.Roots {
		if ok {
			roots = append(roots, r)
		}
	}
	q := Queue.New(roots...)
	for i := range q.Chan() {
		nodes := make([]int, 0, 2<<7)
		value := g.Values[i]
		for j := g.Nodes[i]; j != 0; j = g.Edges[j].Next {
			do(i, g.Edges[j].Target, value)
			nodes = append(nodes, g.Edges[j].Target)
		}
		q.Next(nodes...)
	}
}

func Graph[V any](cap int) graph[V] {
	if cap < 1 {
		cap = 1
	}
	return graph[V]{
		Roots:  make(map[int]bool),
		Nodes:  make(map[int]int),
		Values: make(map[int]V),
		Edges:  make([]edge, 1, cap),
	}
}
