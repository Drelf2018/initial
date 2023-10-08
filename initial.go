package initial

import (
	"errors"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/Drelf2018/initial/graph"
)

var (
	ErrMethod   = errors.New("tag method name not found")
	ErrPtrNil   = errors.New("pointer parament is nil")
	ErrTypeStr  = errors.New("value v must be a pointer to a struct")
	ErrTagAbs   = errors.New("tag abs must be a field name")
	ErrTagType  = errors.New("tag default's type didn't match")
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

func findMethod(field reflect.Value, name string, in []reflect.Value) bool {
	fn := field.Addr().MethodByName(name)
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

func default_(v any) any {
	if v == nil {
		panic(ErrPtrNil)
	}
	vt := reflect.TypeOf(v)
	if vt.Kind() != reflect.Ptr {
		panic(ErrTypeStr)
	}
	vt = vt.Elem()
	if vt.Kind() != reflect.Struct {
		panic(ErrTypeStr)
	}
	vv := reflect.ValueOf(v).Elem()
	in := []reflect.Value{vv.Addr()}
	for i, l := 0, vt.NumField(); i < l; i++ {
		// check tag
		tag, ok := vt.Field(i).Tag.Lookup("default")
		if !ok {
			continue
		}
		// get real value
		field := vv.Field(i)
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				field.Set(reflect.New(vt.Field(i).Type.Elem()))
			}
			field = field.Elem()
		}
		// check field
		if !field.CanSet() {
			continue
		}
		// func
		switch field.Kind() {
		case reflect.Array, reflect.Chan, reflect.Func, reflect.Map, reflect.Slice, reflect.Struct:
			for _, method := range strings.Split(tag, ";") {
				if method == "" || call(method, field.Addr(), vv.Addr()) {
					continue
				}
				var prefix, name string
				s := strings.Split(method, ".")
				if len(s) == 1 {
					name = s[0]
				} else {
					prefix, name = s[0], s[1]
				}
				if prefix == "range" {
					switch field.Kind() {
					case reflect.Array, reflect.Slice:
					default:
						panic(ErrTagRange)
					}
					for j, k := 0, field.Len(); j < k; j++ {
						if findMethod(field.Index(j), name, in) {
							break
						}
					}
				} else {
					if findMethod(field, name, in) {
						break
					}
				}
			}
			continue
		}
		if !field.IsZero() {
			continue
		}
		// parse tag
		switch field.Kind() {
		case reflect.String:
			field.SetString(tag)
		case reflect.Bool:
			if tag == "true" {
				field.SetBool(true)
			} else if tag == "false" {
				field.SetBool(false)
			} else {
				panic(ErrTagType)
			}
		case reflect.Float64:
			f, err := strconv.ParseFloat(tag, 64)
			if err != nil {
				panic(err)
			}
			field.SetFloat(f)
		case reflect.Int64:
			j, err := strconv.ParseInt(tag, 10, 64)
			if err != nil {
				panic(err)
			}
			field.SetInt(j)
		}
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
