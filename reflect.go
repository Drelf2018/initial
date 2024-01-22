package initial

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/Drelf2018/TypeGo/Reflect"
	"github.com/Drelf2018/initial/tag"
)

type value struct {
	reflect.Value
	methods *tag.Sentence
}

func (v value) String() string {
	l := 0
	s := ""
	if v.methods != nil {
		l = len(v.methods.Body)
	}
	if l > 1 {
		s = "s"
	}
	return fmt.Sprintf("%v<%d method%s>", v.Value, l, s)
}

func (v *value) IsValid() bool {
	if v.Value.IsValid() {
		return true
	}
	return v.methods != nil
}

func (v *value) zero(typ reflect.Type, x any) {
	v.Value = Zero(typ)
	switch i := x.(type) {
	case int64:
		v.SetInt(i)
	case uint64:
		v.SetUint(i)
	case float64:
		v.SetFloat(i)
	case complex128:
		v.SetComplex(i)
	}
}

// Zero returns a Value representing the zero value for the specified type.
// The result is addressable and settable.
func Zero(typ reflect.Type) reflect.Value {
	return reflect.New(typ).Elem()
}

func parseAny[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}

var ref = Reflect.New(func(self *Reflect.Reflect[value], field reflect.StructField, elem reflect.Type) (v value) {
	val, ok := field.Tag.Lookup("default")
	if val == "-" {
		return
	}
	typ := field.Type
	if !ok && typ.Kind() != reflect.Struct {
		return
	}
	switch k := typ.Kind(); k {
	case reflect.String:
		v.Value = reflect.ValueOf(val)
	case reflect.Bool:
		v.Value = reflect.ValueOf(parseAny(strconv.ParseBool(val)))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bitSize := 1 << k
		if k == reflect.Int {
			bitSize = 0
		}
		v.zero(typ, parseAny(strconv.ParseInt(val, 10, bitSize)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bitSize := 1 << (k - 5)
		if k == reflect.Uint {
			bitSize = 0
		}
		v.zero(typ, parseAny(strconv.ParseUint(val, 10, bitSize)))
	case reflect.Float32, reflect.Float64:
		v.zero(typ, parseAny(strconv.ParseFloat(val, 1<<(k-8))))
	case reflect.Complex64, reflect.Complex128:
		v.zero(typ, parseAny(strconv.ParseComplex(val, 1<<(k-9))))
	default:
		self.GetType(field.Type, nil)
		v.methods = tag.NewParser(val).Sentence
	}
	return
})
