package graph

import "github.com/Drelf2018/TypeGo/Queue"

type Node[V, E any] struct {
	Index  int
	Value  V
	Parent *Node[V, E]
	Edge   *Edge[V, E]
}

func (n *Node[V, E]) IsRoot() bool {
	return n.Parent == nil
}

func (n *Node[V, E]) Set(value V) {
	n.Value = value
}

func (n *Node[V, E]) Get() V {
	return n.Value
}

type Edge[V, E any] struct {
	Value E
	Node  *Node[V, E]
	Next  *Edge[V, E]
}

func (e *Edge[V, E]) Set(value E) {
	e.Value = value
}

func (e *Edge[V, E]) Get() E {
	return e.Value
}

type Graph[V, E any] map[int]*Node[V, E]

func (g Graph[V, E]) Node(index int) *Node[V, E] {
	node, ok := g[index]
	if !ok {
		node = &Node[V, E]{Index: index}
		g[index] = node
	}
	return node
}

func (g Graph[V, E]) To(index int, parent *Node[V, E]) *Node[V, E] {
	node := g.Node(index)
	node.Parent = parent
	return node
}

func (g Graph[V, E]) Add(from, to int, edge E) {
	node := g.Node(from)
	node.Edge = &Edge[V, E]{
		Value: edge,
		Node:  g.To(to, node),
		Next:  node.Edge,
	}
}

func (g Graph[V, E]) BFS(do func(from, to *Node[V, E], edge *Edge[V, E])) {
	roots := make([]*Node[V, E], 0, len(g))
	for _, n := range g {
		if n.IsRoot() {
			roots = append(roots, n)
		}
	}
	q := Queue.New(roots...)
	for i := range q.Chan() {
		nodes := make([]*Node[V, E], 0, 32)
		for edge := i.Edge; edge != nil; edge = edge.Next {
			do(i, edge.Node, edge)
			nodes = append(nodes, edge.Node)
		}
		q.Next(nodes...)
	}
}
