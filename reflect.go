package initial

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Drelf2018/TypeGo/Reflect"
)

type value struct {
	v       reflect.Value
	vs      []value
	methods []string
}

func (v value) String() string {
	return fmt.Sprintf("v(%v%v)", v.v, vals(v.vs))
}

type vals []value

func (t vals) String() string {
	l := len(t)
	if l == 0 {
		return ""
	}
	buf := bytes.NewBufferString(", [")
	for i, f := range t {
		buf.WriteString(f.String())
		if i != l-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString("]")
	return buf.String()
}

func (v *value) set(i any, err error) {
	if err != nil {
		panic(err)
	}
	v.v = reflect.ValueOf(i)
}

var ref = Reflect.New(func(self *Reflect.Reflect[value], field reflect.StructField) (v value) {
	self.GetType(field.Type, &v.vs)
	val := field.Tag.Get("default")
	switch field.Type.Kind() {
	case reflect.String:
		v.v = reflect.ValueOf(val)
	case reflect.Bool:
		v.set(strconv.ParseBool(val))
	case reflect.Int64:
		v.set(strconv.ParseInt(val, 10, 64))
	case reflect.Uint64:
		v.set(strconv.ParseUint(val, 10, 64))
	case reflect.Float64:
		v.set(strconv.ParseFloat(val, 64))
	default:
		v.methods = strings.Split(val, ";")
	}
	return
})
