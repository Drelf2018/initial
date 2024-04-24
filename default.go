package initial

import (
	"errors"
	"reflect"
)

var (
	ErrNilValue         = errors.New("initial: the value is nil")
	ErrInvalidValue     = errors.New("initial: invalid value")
	ErrNotStruct        = errors.New("initial: the value must be a struct")
	ErrCannotAddrStruct = errors.New("initial: the struct value must be addressable")

	BoolBreak = true
	ErrBreak  = errors.New("initial: the error is not a real error, it is only used to interrupt")
)

func DefaultStruct(v reflect.Value) error {
	if v.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	if !v.CanAddr() {
		return ErrCannotAddrStruct
	}

	parent := v.Addr()

	values, _ := valuesMap.GetType(v.Type())
	for _, value := range values {
		elem := v.Field(value.Index)

		if elem.Kind() == reflect.Ptr {
			if elem.IsNil() {
				elem.Set(reflect.New(elem.Type().Elem()))
			}
			elem = elem.Elem()
		}

		if value.Value.IsValid() && elem.IsZero() {
			elem.Set(value.Value)
		}

		if value.Before != nil {
			err := value.Before.Call(elem.Addr(), parent)
			if err == ErrBreak {
				continue
			}
			if err != nil {
				return err
			}
		}

		if value.Initial {
			err := DefaultValue(elem)
			if err != nil {
				return err
			}
		}

		if value.After != nil {
			err := value.After.Call(elem.Addr(), parent)
			if err != nil && err != ErrBreak {
				return err
			}
		}
	}
	return nil
}

type BeforeDefault interface {
	BeforeDefault() error
}

type AfterDefault interface {
	AfterDefault() error
}

func DefaultValue(v reflect.Value) (err error) {
	if !v.IsValid() {
		return ErrInvalidValue
	}

	if vb, ok := v.Interface().(BeforeDefault); ok {
		err = vb.BeforeDefault()
		if err != nil {
			return
		}
	}

	elem := reflect.Indirect(v)

	switch elem.Kind() {
	case reflect.Struct:
		err = DefaultStruct(elem)
		if err != nil {
			return
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < elem.Len(); i++ {
			err = DefaultValue(elem.Index(i))
			if err != nil {
				return
			}
		}
	case reflect.Map:
		iter := elem.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			err = DefaultValue(k)
			if err != nil {
				return
			}
			err = DefaultValue(v)
			if err != nil {
				return
			}
		}
	}

	if va, ok := v.Interface().(AfterDefault); ok {
		err = va.AfterDefault()
	}
	return
}

func Default(v any) error {
	if v == nil {
		return ErrNilValue
	}
	return DefaultValue(reflect.ValueOf(v))
}

func New[T any]() (*T, error) {
	t := new(T)
	return t, Default(t)
}

func MustNew[T any]() *T {
	t := new(T)
	if Default(t) != nil {
		return nil
	}
	return t
}
