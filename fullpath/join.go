package fullpath

import (
	"errors"
	"path/filepath"
	"reflect"

	"github.com/Drelf2018/reflectMap"
	"golang.org/x/exp/slices"
)

type Node struct {
	Target *Node
	depth  int
}

func (n *Node) Depth() int {
	if n.depth == 0 {
		n.depth = n.Target.Depth() + 1
	}
	return n.depth
}

var (
	ErrInvalidName = errors.New("fullpath: invalid join tag value")
	ErrNotPointer  = errors.New("fullpath: value 'in' must be a pointer to struct")
)

var valuesMap = reflectMap.New(func(m *reflectMap.Map[[][2]int], elem reflect.Type) (edges [][2]int) {
	names := make(map[string]int)

	for idx, field := range reflectMap.FieldsOf(elem) {
		names[field.Name] = idx

		if !field.IsExported() {
			continue
		}

		if field.Type.Kind() != reflect.String {
			continue
		}

		name := field.Tag.Get("join")

		if name == "" {
			edges = append(edges, [2]int{idx, -1})
		} else {
			i, ok := names[name]
			if !ok {
				f, exists := elem.FieldByName(name)
				if !exists {
					panic(ErrInvalidName)
				}
				i = f.Index[0]
				names[name] = i
			}
			edges = append(edges, [2]int{idx, i})
		}

	}

	nodes := make([]Node, elem.NumField())
	for _, data := range edges {
		if data[1] == -1 {
			nodes[data[0]].depth = 1
		} else {
			nodes[data[0]].Target = &nodes[data[1]]
		}
	}

	slices.SortFunc(edges, func(a, b [2]int) (i int) {
		i = nodes[a[0]].Depth() - nodes[b[0]].Depth()
		if i == 0 {
			i = a[1] - b[1]
		}
		if i == 0 {
			i = a[0] - b[0]
		}
		return
	})

	return edges
})

func Edges(in any) [][2]int {
	return valuesMap.Get(in)
}

func Join(in any) error {
	v := reflect.ValueOf(in)
	if v.Kind() != reflect.Pointer {
		return ErrNotPointer
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return ErrNotPointer
	}

	values := make([]string, v.NumField())

	for _, edge := range valuesMap.Get(in) {
		index := edge[0]
		target := edge[1]
		if target == -1 {
			values[index] = v.Field(index).String()
		} else {
			field := v.Field(index)
			values[index] = filepath.Join(values[target], field.String())
			field.SetString(values[index])
		}
	}

	return nil
}

func New[T any](in T) (*T, error) {
	return &in, Join(&in)
}

func MustNew[T any](in T) (t *T) {
	Join(&in)
	return &in
}
