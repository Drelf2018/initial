package initial

import (
	"path/filepath"
	"reflect"

	"github.com/Drelf2018/initial/graph"
)

func Abs[T any](dst, src *T) *T {
	return AbsUnsafe(dst, src).(*T)
}

func AbsUnsafe(dst, src any) any {
	if dst == nil {
		panic(ErrPtrNil)
	}
	vt := reflect.TypeOf(dst)
	if vt.Kind() != reflect.Ptr {
		panic(ErrTypeStr)
	}
	vt = vt.Elem()
	if vt.Kind() != reflect.Struct {
		panic(ErrTypeStr)
	}
	vv := reflect.ValueOf(dst).Elem()
	vs := reflect.ValueOf(src).Elem()

	g := graph.Make[string, int]()
	for i, l := 0, vt.NumField(); i < l; i++ {
		abs := vt.Field(i).Tag.Get("abs")
		if abs == "" {
			vi := vv.Field(i)
			vsi := vs.Field(i)
			if vi.Kind() == reflect.Ptr {
				vi = vi.Elem()
				vsi = vsi.Elem()
			}
			if vi.Kind() == reflect.String {
				if vi.IsZero() {
					vi.SetString(vsi.String())
				}
				g.Node(i).Set(vi.String())
			}
		} else {
			field, ok := vt.FieldByName(abs)
			if !ok {
				panic(ErrTagAbs)
			}
			g.Add(field.Index[0], i, 0)
		}
	}

	g.BFS(func(from, to *graph.Node[string, int], edge *graph.Edge[string, int]) {
		vvt := vv.Field(to.Index)
		if vvt.Kind() == reflect.Ptr {
			vvt = vvt.Elem()
		}
		val := vvt.String()
		if val == "" {
			vo := vs.Field(to.Index)
			if vo.Kind() == reflect.Ptr {
				vo = vo.Elem()
			}
			val = vo.String()
		}
		value := filepath.Join(from.Get(), val)
		vvt.SetString(value)
		to.Set(value)
	})
	return dst
}
