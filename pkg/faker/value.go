package faker

import (
	"crypto/rand"
	"fmt"
	mathRand "math/rand"
	"reflect"
	"strings"
)

var possibleStrings = []string{
	"Alpha",
	"Beta",
	"Gamma",
	"Delta",
	"Test",
}

func CreateRandomValue(
	t reflect.Value,
) reflect.Value {
	var err error

	v := make([]byte, 10)
	_, err = rand.Read(v)
	if err != nil {
		panic(err)
	}

	kind := t.Kind().String()
	switch {
	case strings.HasPrefix(kind, "int"):
		valueNumber := int64(0)
		for _, b := range v {
			valueNumber += int64(b)
		}
		res := reflect.New(t.Type())
		res.Elem().SetInt(valueNumber)

		return res.Elem()

	case strings.HasPrefix(kind, "uint"):
		valueNumber := uint64(0)
		for _, b := range v {
			valueNumber += uint64(b)
		}
		res := reflect.New(t.Type())
		res.Elem().SetUint(valueNumber)

		return res.Elem()

	case strings.HasPrefix(kind, "float"):
		valueNumber := float64(0)
		for _, b := range v {
			valueNumber += float64(b)
		}
		res := reflect.New(t.Type())
		res.Elem().SetFloat(valueNumber)

		return res.Elem()

	case strings.HasPrefix(kind, "string"):
		randomString := possibleStrings[mathRand.Intn(len(possibleStrings))]

		res := reflect.New(t.Type())
		res.Elem().SetString(randomString)

		return res.Elem()

	case strings.HasPrefix(kind, "bool"):
		res := reflect.New(t.Type())
		res.Elem().SetBool(true)

		return res.Elem()

	case kind == "slice":
		res := reflect.New(t.Type()).Elem()
		member := reflect.New(t.Type().Elem()).Elem()

		res = reflect.Append(res, reflect.ValueOf(CreateRandomValue(member)))

		return res

	case kind == "map":
		res := reflect.MakeMapWithSize(t.Type(), 0)
		key := reflect.New(t.Type().Key()).Elem()
		member := reflect.New(t.Type().Elem()).Elem()

		res.SetMapIndex(reflect.ValueOf(CreateRandomValue(key)), reflect.ValueOf(CreateRandomValue(member)))

		return res

	default:
		var f interface{}
		if kind == "struct" {
			f = reflect.New(t.Type()).Interface()
		} else {
			f = reflect.New(t.Type().Elem()).Interface()
		}

		panic(fmt.Errorf("can not set type %T", f))
	}
}
