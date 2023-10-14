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
	ErrPtrNil   = errors.New("pointer parament is nil")
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

func Default[P any, T *P](v T) T {
	return default_(v).(T)
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

func execMethod(v, parent reflect.Value, val value) {
	for _, method := range val.methods {
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
				execMethod(v.Index(j), parent, value{methods: []string{name}})
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
		panic(ErrPtrNil)
	}
	vv := reflect.ValueOf(v).Elem()
	for i, val := range ref.Get(v) {
		value := vv.Field(i)
		// create ptr value
		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				value.Set(reflect.New(value.Type().Elem()))
			}
			value = value.Elem()
		}
		// ckeck
		if !value.CanSet() {
			continue
		}
		// set value
		if value.IsZero() && val.v.IsValid() {
			value.Set(val.v)
			continue
		}
		execMethod(value, vv.Addr(), val)
	}
	return v
}

func Abs[P any, T *P](dst, src T) T {
	return abs(dst, src).(T)
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

	g := make(graph.Graph[string, int])
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
