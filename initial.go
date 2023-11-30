package initial

import (
	"errors"
	"reflect"

	"github.com/Drelf2018/initial/tag"
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

type Arg int

const (
	SELF Arg = iota
	PARENT
)

type fn struct {
	fn   reflect.Value
	args []Arg
}

var functions = make(map[string]fn)

func Register(name string, f any, args ...Arg) {
	functions[name] = fn{reflect.ValueOf(f), args}
}

func needBreak(out []reflect.Value) bool {
	if len(out) == 0 {
		return false
	}
	// return true means break
	o := out[0].Interface()
	switch o := o.(type) {
	case bool:
		return o
	case error:
		return o != nil
	default:
		return false
	}
}

func call(self, parent reflect.Value, fn fn) bool {
	if fn.fn.IsZero() {
		return false
	}
	var in []reflect.Value
	for _, arg := range fn.args {
		switch arg {
		case SELF:
			in = append(in, self)
		case PARENT:
			in = append(in, parent)
		}
	}
	return needBreak(fn.fn.Call(in))
}

func init() {
	Register("initial.Abs", AbsUnsafe, SELF, PARENT)
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
	r := DefaultUnsafe(v)
	if a, ok := r.(AfterDefault); ok {
		a.AfterDefault()
	}
	return r.(*T)
}

func findMethod(v, parent reflect.Value, method string) (bool, error) {
	var fn reflect.Value
	if v.CanAddr() {
		fn = v.Addr().MethodByName(method)
	} else {
		fn = v.MethodByName(method)
	}
	if !fn.IsValid() {
		return true, ErrMethod
	}
	return needBreak(fn.Call([]reflect.Value{parent.Addr()})), nil
}

func execMethod(v, parent reflect.Value, method string) bool {
	if method == "" {
		return false
	}
	if fn, ok := functions[method]; ok {
		return call(v.Addr(), parent.Addr(), fn)
	}
	b, err := findMethod(v, parent, method)
	if err != nil {
		panic(err)
	}
	return b
}

func execMethods(v, parent reflect.Value, methods *tag.Sentence) {
	findMethod(v, parent, "BeforeInitial")
	defer findMethod(v, parent, "AfterInitial")
	if len(methods.Body) == 0 || methods.Body[0].Value != "-" {
		if v.Kind() == reflect.Struct {
			DefaultUnsafe(v.Addr().Interface())
		}
	}
	for _, method := range methods.Body {
		if method.Value == "-" {
			continue
		}
		switch method.Kind {
		case tag.LBRACKET:
			if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
				panic(ErrTagRange)
			}
			for j, k := 0, v.Len(); j < k; j++ {
				execMethods(v.Index(j), parent, method)
			}
		case tag.LBRACE:
			if v.Kind() != reflect.Map {
				panic(ErrTagRange)
			}
			for _, key := range v.MapKeys() {
				execMethods(key, parent, method.Body[0])
				execMethods(v.MapIndex(key), parent, method.Body[1])
			}
		default:
			if execMethod(v, parent, method.Value) {
				return
			}
		}
	}
}

func DefaultUnsafe(v any) any {
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
		if value.Value.IsValid() {
			if vi.IsZero() {
				vi.Set(value.Value)
			}
		} else {
			execMethods(vi, vv, value.methods)
		}
	}
	return v
}
