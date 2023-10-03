package initial

import (
	"errors"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrMethod  = errors.New("tag method name not found")
	ErrPtrNil  = errors.New("pointer parament is nil")
	ErrTypeStr = errors.New("value v must be a pointer to a struct")
	ErrTagAbs  = errors.New("tag abs must be a field name")
	ErrTagType = errors.New("tag default's type didn't match")
)

type fn struct {
	fn   reflect.Value
	args []string
}

var functions map[string]fn

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
	functions = make(map[string]fn)
	Register("-", func() {})
	Register("initial.Abs", abs, "self", "parent")
	Register("initial.Default", default_, "self")
}

func Default[P any, T *P](v T) T {
	return default_(v).(T)
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
		if !field.IsZero() && field.Kind() != reflect.Struct {
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
		case reflect.Struct:
			for _, t := range strings.Split(tag, ";") {
				if !call(t, field.Addr(), vv.Addr()) {
					fn := field.Addr().MethodByName(t)
					if !fn.IsValid() {
						panic(ErrMethod)
					}
					fn.Call([]reflect.Value{vv.Addr()})
				}
			}
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

	g := Graph[string](0)
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
				g.Values[i] = vi.String()
			}
		} else {
			field, ok := vt.FieldByName(abs)
			if !ok {
				panic(ErrTagAbs)
			}
			g.Add(field.Index[0], i)
		}
	}

	g.BFS(func(from, to int, value string) {
		vvt := vv.Field(to)
		if vvt.Kind() == reflect.Ptr {
			vvt = vvt.Elem()
		}
		val := vvt.String()
		if val == "" {
			vo := vs.Field(to)
			if vo.Kind() == reflect.Ptr {
				vo = vo.Elem()
			}
			val = vo.String()
		}
		value = filepath.Join(value, val)
		vvt.SetString(value)
		g.Values[to] = value
	})
	return dst
}
