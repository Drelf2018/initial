package initial

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/Drelf2018/TypeGo/Reflect"
)

type value struct {
	reflect.Value
	methods []string
}

func (v *value) set(i any, err error) {
	if err != nil {
		panic(err)
	}
	v.Value = reflect.ValueOf(i)
}

func (v *value) IsValid() bool {
	if v.Value.IsValid() {
		return true
	}
	return len(v.methods) != 0
}

var ref = Reflect.New(func(self *Reflect.Reflect[value], field reflect.StructField, elem reflect.Type) (v value) {
	self.GetType(field.Type, nil)
	val, ok := field.Tag.Lookup("default")
	if !ok {
		return
	}
	switch field.Type.Kind() {
	case reflect.String:
		v.set(val, nil)
	case reflect.Bool:
		v.set(strconv.ParseBool(val))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.set(strconv.ParseInt(val, 10, 64))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.set(strconv.ParseUint(val, 10, 64))
	case reflect.Float32, reflect.Float64:
		v.set(strconv.ParseFloat(val, 64))
	case reflect.Complex64, reflect.Complex128:
		v.set(strconv.ParseComplex(val, 128))
	default:
		v.methods = strings.Split(val, ";")
	}
	return
})
