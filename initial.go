package initial

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrInvalidValue         = errors.New("initial: invalid value")
	ErrNotStruct            = errors.New("initial: the value must be a struct")
	ErrNotAddressableStruct = errors.New("initial: the struct value must be addressable")
)

func Call(fn, parent reflect.Value) error {
	result := fn.Call([]reflect.Value{parent})
	if len(result) == 0 {
		return nil
	}
	v := result[0].Interface()
	err, ok := v.(error)
	if ok {
		return err
	}
	return fmt.Errorf("initial: %#v is not an error", v)
}

func InitialStruct(v reflect.Value) error {
	if v.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	if !v.CanAddr() {
		return ErrNotAddressableStruct
	}
	parent := v.Addr()
	for _, value := range Load(v.Type()) {
		elem := v.Field(value.Index)
		if elem.Kind() == reflect.Pointer {
			if elem.IsNil() {
				elem.Set(reflect.New(elem.Type().Elem()))
			}
			elem = elem.Elem()
		}
		if value.Default.IsValid() && elem.IsZero() {
			elem.Set(value.Default)
		}
		if value.Before.IsValid() {
			err := Call(value.Before, parent)
			if err != nil {
				return err
			}
		}
		if value.Initial {
			err := InitialValue(elem)
			if err != nil {
				return err
			}
		}
		if value.After.IsValid() {
			err := Call(value.After, parent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func InitialValue(v reflect.Value) error {
	if !v.IsValid() {
		return ErrInvalidValue
	}

	var i any
	if v.Kind() != reflect.Pointer && v.CanAddr() {
		i = v.Addr().Interface()
	} else {
		i = v.Interface()
	}

	if before, ok := i.(BeforeInitial1); ok {
		before.BeforeInitial()
	} else if before, ok := i.(BeforeInitial2); ok {
		err := before.BeforeInitial()
		if err != nil {
			return err
		}
	}

	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.Struct:
		err := InitialStruct(v)
		if err != nil {
			return err
		}
	case reflect.Array, reflect.Slice:
		var err error
		for i := 0; i < v.Len(); i++ {
			err = InitialValue(v.Index(i))
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		var err error
		iter := v.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			err = InitialValue(k)
			if err != nil {
				return err
			}
			err = InitialValue(v)
			if err != nil {
				return err
			}
		}
	}

	if after, ok := i.(AfterInitial1); ok {
		after.AfterInitial()
	} else if after, ok := i.(AfterInitial2); ok {
		err := after.AfterInitial()
		if err != nil {
			return err
		}
	}
	return nil
}

func Initial(v any) error {
	if v == nil {
		return nil
	}
	return InitialValue(reflect.ValueOf(v))
}

func New[T any]() (*T, error) {
	t := new(T)
	return t, Initial(t)
}

func MustNew[T any]() *T {
	t := new(T)
	if Initial(t) != nil {
		return nil
	}
	return t
}
