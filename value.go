package initial

import (
	"reflect"
	"strconv"

	"github.com/Drelf2018/reflectMap"
)

func Indirect(t reflect.Type) reflect.Type {
	if t.Kind() != reflect.Pointer {
		return t
	}
	return t.Elem()
}

func IsOrdinaryValue(fieldType reflect.Type) bool {
	switch fieldType.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		return true
	default:
		return false
	}
}

// ParseOrdinaryValue convert val to reflect.Value with type fieldType.
//
// If the usual Go conversion rules do not allow this conversion, ParseOrdinaryValue panics.
func ParseOrdinaryValue(fieldType reflect.Type, val string) reflect.Value {
	if val == "" {
		return reflect.Value{}
	}

	var (
		x   any
		err error
	)

	kind := fieldType.Kind()
	switch kind {
	case reflect.String:
		x = val
	case reflect.Bool:
		x, err = strconv.ParseBool(val)
	case reflect.Int:
		x, err = strconv.ParseInt(val, 10, 0)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err = strconv.ParseInt(val, 10, 1<<kind)
	case reflect.Uint:
		x, err = strconv.ParseUint(val, 10, 0)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, err = strconv.ParseUint(val, 10, 1<<(kind-5))
	case reflect.Float32, reflect.Float64:
		x, err = strconv.ParseFloat(val, 1<<(kind-8))
	case reflect.Complex64, reflect.Complex128:
		x, err = strconv.ParseComplex(val, 1<<(kind-9))
	}

	if err != nil {
		panic(err)
	}
	return reflect.ValueOf(x).Convert(fieldType)
}

type Value struct {
	reflect.Value
	Index   int
	Before  *Method
	Initial bool
	After   *Method
}

var defaultValues map[string]reflect.Value

func SetDefaultValue(name string, value any) {
	defaultValues[name] = reflect.ValueOf(value)
}

var valuesMap = reflectMap.New(func(m *reflectMap.Map[[]Value], elem reflect.Type) (values []Value) {
	parent := elem.Name()

	for idx, field := range reflectMap.FieldsOf(elem) {
		if !field.IsExported() {
			continue
		}

		var (
			defaultTag = field.Tag.Get("default")
			initialTag = field.Tag.Get("initial")

			fieldElem = Indirect(field.Type)
			fieldPtr  = reflect.PointerTo(fieldElem)
		)

		if initialTag != "-" {
			if initialTag == "" {
				initialTag = parent + field.Name
			}
		}

		value := Value{
			Index:   idx,
			Before:  NewMethod(fieldPtr.MethodByName("Before" + initialTag)),
			Initial: initialTag != "-",
			After:   NewMethod(fieldPtr.MethodByName("After" + initialTag)),
		}

		if IsOrdinaryValue(fieldElem) {
			value.Value = ParseOrdinaryValue(fieldElem, defaultTag)
		} else if defaultTag != "" {
			value.Value = defaultValues[defaultTag]
		}

		if value.Initial || value.Value.IsValid() {
			values = append(values, value)
		}
	}
	return
})

func Get(in any) []Value {
	return valuesMap.Get(in)
}

func GetType(in reflect.Type) (values []Value, ok bool) {
	return valuesMap.GetType(in)
}
