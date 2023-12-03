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
		// set value
		if value.Value.IsValid() {
			if vi.CanSet() && vi.IsZero() {
				vi.Set(value.Value)
			}
		} else {
			execMethods(vi, vv, value.methods)
		}
	}
	return v
}

func execMethods(v, parent reflect.Value, methods *tag.Sentence) {
	CallMethod(v, "BeforeInitial", parent.Addr())
	defer CallMethod(v, "AfterInitial", parent.Addr())

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
			iter := v.MapRange()
			for iter.Next() {
				execMethods(iter.Key(), parent, method.Body[0])
				execMethods(iter.Value(), parent, method.Body[1])
			}
		default:
			if execMethod(v, parent, method.Value) {
				return
			}
		}
	}
}

func execMethod(v, parent reflect.Value, method string) bool {
	if method == "" {
		return false
	}
	if fn, ok := functions[method]; ok {
		return fn.CallWithArgs(v.Addr(), parent.Addr())
	}
	fn, ok := FindMethod(v, method)
	if !ok {
		panic(ErrMethod)
	}
	return fn.Call(parent.Addr())
}
