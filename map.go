package initial

import (
	"reflect"
	"sync"
	"unsafe"
)

var cache sync.Map

func Load(in reflect.Type) []Value {
	ptr := ValuePtr(in)
	if v, ok := cache.Load(ptr); ok {
		return v.([]Value)
	}
	actual, _ := cache.LoadOrStore(ptr, ParseValues(in))
	return actual.([]Value)
}

type Any struct {
	Type  unsafe.Pointer
	Value unsafe.Pointer
}

func TypePtr(in any) uintptr {
	return uintptr((*Any)(unsafe.Pointer(&in)).Type)
}

// ValuePtr can obtain the uintptr of a type from its reflect.Type
//
//	TypePtr(something{}) is equal to ValuePtr(reflect.TypeOf(something{}))
func ValuePtr(in any) uintptr {
	return uintptr((*Any)(unsafe.Pointer(&in)).Value)
}
