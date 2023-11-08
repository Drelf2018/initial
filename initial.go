package initial

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Drelf2018/initial/graph"
)

var (
	ErrMethod   = errors.New("tag method name not found")
	ErrNil      = errors.New("pass in value is nil")
	ErrPtrNil   = errors.New("pass in pointer's value is nil")
	ErrTypeStr  = errors.New("value v must be a pointer to a struct")
	ErrTagAbs   = errors.New("tag abs must be a field name")
	ErrTagRange = errors.New("can't range value v")
	ErrBreak    = errors.New("isn't a real error, just used to break loop")
)

type fn struct {
	fn   reflect.Value
	args []string
}

var functions = make(map[string]fn)

func Register(name string, f any, args ...string) {
	functions[name] = fn{reflect.ValueOf(f), args}
}

func call(name string, self, parent reflect.Value) bool {
	if v, ok := functions[name]; ok {
		if v.fn.IsZero() {
			return false
		}
		var in []reflect.Value
		for _, arg := range v.args {
			switch arg {
			case "self":
				in = append(in, self)
			case "parent":
				in = append(in, parent)
			}
		}
		v.fn.Call(in)
		return true
	}
	return false
}

func init() {
	Register("-", func() {})
	Register("initial.Abs", abs, "self", "parent")
	Register("initial.Default", default_, "self")
}

type BeforeDefault interface {
	BeforeDefault()
}

type AfterDefault interface {
	AfterDefault()
}

func Default[T any](v *T) *T {
	if b, ok := any(v).(BeforeDefault); ok {
		b.BeforeDefault()
	}
	r := default_(v)
	if a, ok := r.(AfterDefault); ok {
		a.AfterDefault()
	}
	return r.(*T)
}

func findMethod(value reflect.Value, name string, in []reflect.Value) bool {
	fn := value.Addr().MethodByName(name)
	if !fn.IsValid() {
		panic(ErrMethod)
	}
	vs := fn.Call(in)
	if len(vs) == 0 {
		return false
	}
	// return true means break
	v := vs[0].Interface()
	switch v := v.(type) {
	case bool:
		return v
	case error:
		return v != nil
	default:
		return false
	}
}

func splitMethod(method string) (prefix, name string) {
	s := strings.SplitN(method, ".", 2)
	if len(s) == 1 {
		return "", s[0]
	}
	return s[0], s[1]
}

func execMethod(v, parent reflect.Value, methods ...string) {
	for _, method := range methods {
		if method == "" || call(method, v.Addr(), parent) {
			continue
		}
		prefix, name := splitMethod(method)
		switch prefix {
		case "range":
			if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
				panic(ErrTagRange)
			}
			for j, k := 0, v.Len(); j < k; j++ {
				execMethod(v.Index(j), parent, name)
			}
		default:
			if findMethod(v, name, []reflect.Value{parent}) {
				return
			}
		}
	}
}

func default_(v any) any {
	if v == nil {
		panic(ErrNil)
	}
	vv := reflect.ValueOf(v).Elem()
	if !vv.IsValid() {
		panic(ErrPtrNil)
	}
	for idx, value := range ref.Get(v) {
		if !value.IsValid() {
			continue
		}
		vi := vv.Field(idx)
		// create ptr value
		if vi.Kind() == reflect.Ptr {
			if vi.IsNil() {
				vi.Set(reflect.New(vi.Type().Elem()))
			}
			vi = vi.Elem()
		}
		// ckeck
		if !vi.CanSet() {
			continue
		}
		// set value
		if vi.IsZero() && value.Value.IsValid() {
			vi.Set(value.Value)
			continue
		}
		execMethod(vi, vv.Addr(), value.methods...)
	}
	return v
}

func Abs[T any](dst, src *T) *T {
	return abs(dst, src).(*T)
}

func abs(dst, src any) any {
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
