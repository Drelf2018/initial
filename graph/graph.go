package graph

import (
	"github.com/Drelf2018/TypeGo/Pool"
	"github.com/Drelf2018/TypeGo/Queue"
)

type Node[V, E any] struct {
	Index  int
	Value  V
	Parent *Node[V, E]
	Edge   *Edge[V, E]
}

func (n *Node[V, E]) IsRoot() bool {
	return n.Parent == nil
}

func (n *Node[V, E]) Set(value V) *Node[V, E] {
	n.Value = value
	return n
}

func (n *Node[V, E]) Get() V {
	return n.Value
}

type Edge[V, E any] struct {
	Value E
	Node  *Node[V, E]
	Next  *Edge[V, E]
}

func (e *Edge[V, E]) New() {}

func (e *Edge[V, E]) Reset() {
	var v E
	e.Value, e.Node, e.Next = v, nil, nil
}

func (e *Edge[V, E]) Set(x ...any) {
	e.Value, e.Node, e.Next = x[0].(E), x[1].(*Node[V, E]), x[2].(*Edge[V, E])
}

func (e *Edge[V, E]) SetValue(value E) *Edge[V, E] {
	e.Value = value
	return e
}

func (e *Edge[V, E]) Get() E {
	return e.Value
}

type Graph[V, E any] struct {
	Map  map[int]*Node[V, E]
	Pool Pool.TypePool[*Edge[V, E]]
}

func Make[V, E any]() *Graph[V, E] {
	return &Graph[V, E]{make(map[int]*Node[V, E]), Pool.New(&Edge[V, E]{})}
}

func (g *Graph[V, E]) Node(index int, parent ...*Node[V, E]) *Node[V, E] {
	node, ok := g.Map[index]
	if !ok {
		node = &Node[V, E]{Index: index}
		g.Map[index] = node
	}
	if len(parent) != 0 {
		node.Parent = parent[0]
	}
	return node
}

func (g *Graph[V, E]) Add(from, to int, edge E) (first *Node[V, E]) {
	node := g.Node(from)
	node.Edge = g.Pool.Get(edge, g.Node(to, node), node.Edge)
	return node
}

type Walker[V, E any] func(from, to *Node[V, E], edge *Edge[V, E])

func (g *Graph[V, E]) Walk(start []*Node[V, E], do Walker[V, E]) {
	q := Queue.New(start...)
	for !q.IsEmpty() {
		nodes := make([]*Node[V, E], 0, 32)
		for edge := q.MustPop().Edge; edge != nil; edge = edge.Next {
			do(edge.Node.Parent, edge.Node, edge)
			nodes = append(nodes, edge.Node)
		}
		q.Push(nodes...)
	}
}

func (g *Graph[V, E]) BFS(do Walker[V, E]) {
	roots := make([]*Node[V, E], 0, len(g.Map))
	for _, n := range g.Map {
		if n.IsRoot() {
			roots = append(roots, n)
		}
	}
	g.Walk(roots, do)
}

func (g *Graph[V, E]) Delete(node *Node[V, E]) {
	for edge := node.Edge; edge != nil; edge = edge.Next {
		defer g.Pool.Put(edge)
	}
	delete(g.Map, node.Index)
}

func (g *Graph[V, E]) Clear() {
	for _, n := range g.Map {
		g.Delete(n)
	}
}
