package initial

import "reflect"

type Arg int

const (
	SELF Arg = iota
	PARENT
)

type Fn struct {
	reflect.Value
	Args []Arg
}

func (fn *Fn) Call(in ...reflect.Value) bool {
	out := fn.Value.Call(in)
	if len(out) == 0 {
		return false
	}
	// return true or error means break
	switch r := out[0].Interface().(type) {
	case bool:
		return r
	case error:
		return r != nil
	default:
		return false
	}
}

func (fn *Fn) CallWithArgs(self, parent reflect.Value) bool {
	if fn.IsZero() {
		return false
	}
	var in []reflect.Value
	for _, arg := range fn.Args {
		switch arg {
		case SELF:
			in = append(in, self)
		case PARENT:
			in = append(in, parent)
		}
	}
	return fn.Call(in...)
}

func FindMethod(v reflect.Value, method string) (fn Fn, ok bool) {
	if v.CanAddr() {
		fn.Value = v.Addr().MethodByName(method)
	} else {
		fn.Value = v.MethodByName(method)
	}
	return fn, fn.IsValid()
}

func CallMethod(v reflect.Value, method string, in ...reflect.Value) bool {
	fn, ok := FindMethod(v, method)
	if ok {
		fn.Call(in...)
	}
	return ok
}

var functions = make(map[string]Fn)

func Register(name string, f any, args ...Arg) {
	functions[name] = Fn{reflect.ValueOf(f), args}
}

func RegisterFn(name string, fn Fn) {
	functions[name] = fn
}

func init() {
	Register("initial.Abs", AbsUnsafe, SELF, PARENT)
}
