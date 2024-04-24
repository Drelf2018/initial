package initial

import (
	"reflect"
)

type Method struct {
	numIn int
	fn    reflect.Value
}

func (m *Method) Call(ptr, parent reflect.Value) error {
	var out []reflect.Value
	switch m.numIn {
	case 1:
		out = m.fn.Call([]reflect.Value{ptr})
	case 2:
		out = m.fn.Call([]reflect.Value{ptr, parent})
	}

	if len(out) == 0 {
		return nil
	}
	return ParseOutValue(out[len(out)-1])
}

func ParseOutValue(out reflect.Value) error {
	switch v := out.Interface().(type) {
	case error:
		return v
	case bool:
		if v == BoolBreak {
			return ErrBreak
		}
		return nil
	default:
		return nil
	}
}

func NewMethod(m reflect.Method, ok bool) *Method {
	if ok {
		return &Method{
			numIn: m.Type.NumIn(),
			fn:    m.Func,
		}
	}
	return nil
}
