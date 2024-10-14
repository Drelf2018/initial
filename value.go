package initial

import (
	"reflect"
	"strings"
)

type Value struct {
	Index   int
	Initial bool
	Default reflect.Value
	Before  reflect.Value
	After   reflect.Value
}

func ParseValues(parent reflect.Type) (values []Value) {
	if parent.Kind() == reflect.Pointer {
		parent = parent.Elem()
	}
	if parent.Kind() != reflect.Struct {
		return nil
	}

	var (
		parentPtrType  = reflect.PointerTo(parent)
		parentPtrValue = reflect.New(parent)
	)

	for idx := 0; idx < parent.NumField(); idx++ {
		field := parent.Field(idx)
		if !field.IsExported() {
			continue
		}

		var (
			defaultTag = field.Tag.Get("default")
			initialTag = field.Tag.Get("initial")
			fieldElem  = Indirect(field.Type)
			value      = Value{Index: idx, Initial: initialTag != "-" && !(field.Type.Kind() == reflect.Pointer && IsRecursiveType(field.Type))}
		)

		if strings.HasPrefix(defaultTag, "$") {
			value.Default = defaultValues[defaultTag]
		} else if IsOrdinaryValue(fieldElem) {
			value.Default, _ = ParseOrdinaryValue(fieldElem, defaultTag)
		}

		if initialTag == "" {
			initialTag = field.Name
		}
		before, ok := parentPtrType.MethodByName("Before" + initialTag)
		if ok {
			switch parentPtrValue.Method(before.Index).Interface().(type) {
			case func(), func() error:
				value.Before = before.Func
			}
		}
		after, ok := parentPtrType.MethodByName("After" + initialTag)
		if ok {
			switch parentPtrValue.Method(after.Index).Interface().(type) {
			case func(), func() error:
				value.After = after.Func
			}
		}
		values = append(values, value)
	}
	return
}
