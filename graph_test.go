package initial_test

import (
	"fmt"
	"testing"

	"github.com/Drelf2018/initial"
)

func TestGraph(t *testing.T) {
	g := initial.Graph[int](0)
	g.Add(3, 7)
	g.Add(3, 6)
	g.Add(2, 5)
	g.Add(2, 4)
	g.Add(1, 3)
	g.Add(1, 2)
	g.BFS(func(from, to, _ int) {
		fmt.Printf("from: %v to: %v\n", from, to)
	})
}
