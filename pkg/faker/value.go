package faker

import (
	"fmt"
	"math"
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
	kind := t.Kind().String()
	switch {
	case strings.HasPrefix(kind, "int"):
		valueNumber := mathRand.Int63n(10000)
		res := reflect.New(t.Type())
		res.Elem().SetInt(valueNumber)

		return res.Elem()
	case strings.HasPrefix(kind, "uint"):
		valueNumber := mathRand.Int63n(10000)
		res := reflect.New(t.Type())
		res.Elem().SetUint(uint64(math.Abs(float64(valueNumber))))

		return res.Elem()
	case strings.HasPrefix(kind, "float"):
		valueNumber := mathRand.NormFloat64()
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
