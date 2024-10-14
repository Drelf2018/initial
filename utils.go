package initial

import (
	"reflect"
	"strconv"
)

func Indirect(t reflect.Type) reflect.Type {
	if t.Kind() != reflect.Pointer {
		return t
	}
	return t.Elem()
}

func IsOrdinaryValue(fieldType reflect.Type) bool {
	switch fieldType.Kind() {
	case
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}

func IsRecursiveType(typ reflect.Type) bool {
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return false
	}
	types := make([]reflect.Type, 0)
	fields := make([]reflect.StructField, 0, typ.NumField())
	for idx := 0; idx < typ.NumField(); idx++ {
		fields = append(fields, typ.Field(idx))
	}
outer:
	for i := 0; i < len(fields); i++ {
		field := fields[i]
		if !field.IsExported() {
			continue
		}
		fieldType := field.Type
		if fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() != reflect.Struct {
			continue
		}
		if fieldType == typ {
			return true
		}
		for _, t := range types {
			if t == fieldType {
				continue outer
			}
		}
		types = append(types, fieldType)
		for j := 0; j < fieldType.NumField(); j++ {
			fields = append(fields, fieldType.Field(j))
		}
	}
	return false
}

// ParseOrdinaryValue convert val to reflect.Value with type fieldType.
//
// If the usual Go conversion rules do not allow this conversion, ParseOrdinaryValue returns error.
func ParseOrdinaryValue(fieldType reflect.Type, val string) (v reflect.Value, err error) {
	if val == "" {
		return
	}
	var x any
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
		return
	}
	return reflect.ValueOf(x).Convert(fieldType), nil
}

var defaultValues map[string]reflect.Value

func init() {
	defaultValues = make(map[string]reflect.Value)
}

func SetDefaultValue(name string, value any) {
	defaultValues[name] = reflect.ValueOf(value)
}
