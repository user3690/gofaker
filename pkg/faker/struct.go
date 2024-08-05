package faker

import (
	"fmt"
	"reflect"
	"time"
)

func FakeStruct(ptrStruct any, depth uint8) {
	ptrStructType := reflect.TypeOf(ptrStruct)
	if ptrStructType.Kind() != reflect.Pointer {
		panic("ptrStruct needs to be a pointer to a struct")
	}

	// dereference underlying value
	refStruct := ptrStructType.Elem()
	if refStruct.Kind() != reflect.Struct {
		panic("refStruct needs to be a struct")
	}

	fields := reflect.VisibleFields(refStruct)
	values := reflect.ValueOf(ptrStruct).Elem()

	for _, field := range fields {
		if !field.IsExported() {
			continue
		}

		if field.Type.Kind() == reflect.Pointer {
			fakePtrField(field, ptrStruct, values, depth)

			continue
		}

		if field.Type.Kind() == reflect.Struct {
			fakeStructField(field, ptrStruct, values, depth)

			continue
		}

		if field.Type.Kind() == reflect.Slice {
			fakeSliceField(field, ptrStruct, values, depth)

			continue
		}

		if field.Type.Kind() == reflect.Map {
			fakeMapField(field, ptrStruct, values, depth)

			continue
		}

		structFieldValue := values.FieldByIndex(field.Index)
		value := CreateRandomValue(structFieldValue)

		if !structFieldValue.IsValid() {
			panic(fmt.Sprintf("no such field: %s in obj", field.Name))
		}

		// If obj field value is not settable an error is thrown
		if !structFieldValue.CanSet() {
			panic(fmt.Sprintf("cannot set %s field value", field.Name))
		}

		structFieldType := structFieldValue.Type()
		if structFieldType != value.Type() {
			panic(fmt.Sprintf("struct field type mismatch: %s vs %s", structFieldType, value.Type()))
		}

		structFieldValue.Set(value)
	}
}

func fakeStructRecursive(
	pointerStruct any,
	parentType any,
	depth uint8,
) {
	if reflect.TypeOf(parentType).String() == reflect.TypeOf(pointerStruct).String() {
		return
	}

	depth--
	if depth == 0 {
		return
	}

	FakeStruct(pointerStruct, depth)
}

func fakePtrField(
	field reflect.StructField,
	currentPtrStruct any,
	currentPtrStructValues reflect.Value,
	depth uint8,
) {
	// dereference underlying value
	refValue := field.Type.Elem()
	if refValue.Kind() != reflect.Struct {
		return
	}

	ptrVal := reflect.New(refValue)
	fakeStructRecursive(ptrVal.Interface(), currentPtrStruct, depth)
	currentPtrStructValues.FieldByIndex(field.Index).Set(ptrVal)
}

func fakeStructField(
	field reflect.StructField,
	currentPtrStruct any,
	currentPtrStructValues reflect.Value,
	depth uint8,
) {
	// time.Time struct special case cause of unexported fields
	if field.Type.Name() == "Time" {
		newTime := time.Now()

		currentPtrStructValues.FieldByIndex(field.Index).Set(reflect.ValueOf(newTime))

		return
	}

	// create pointer of field type
	ptrVal := reflect.New(field.Type)
	fakeStructRecursive(ptrVal.Interface(), currentPtrStruct, depth)
	currentPtrStructValues.FieldByIndex(field.Index).Set(ptrVal.Elem())
}

func fakeSliceField(
	field reflect.StructField,
	currentPtrStruct any,
	currentPtrStructValues reflect.Value,
	depth uint8,
) {
	sliceField := currentPtrStructValues.FieldByIndex(field.Index)
	if !sliceField.IsNil() {
		return
	}

	// create slice at 1 length and capacity
	newSlice := reflect.MakeSlice(field.Type, 1, 1)

	// get the first value
	sliceValue := newSlice.Index(0)

	// if value of slice is a pointer, process it again
	if sliceValue.Kind() == reflect.Ptr {
		ptrVal := reflect.New(sliceValue.Type().Elem())
		fakeStructRecursive(ptrVal.Interface(), currentPtrStruct, depth)
		sliceValue.Set(ptrVal)
	} else if sliceValue.Kind() == reflect.Struct {
		ptrVal := reflect.New(sliceValue.Type())
		fakeStructRecursive(ptrVal.Interface(), currentPtrStruct, depth)
		sliceValue.Set(ptrVal.Elem())
	} else {
		value := CreateRandomValue(sliceValue)
		sliceValue.Set(value)
	}

	sliceField.Set(newSlice)
}

func fakeMapField(
	field reflect.StructField,
	currentPtrStruct any,
	currentPtrStructValues reflect.Value,
	depth uint8,
) {
	mapField := currentPtrStructValues.FieldByIndex(field.Index)
	if !mapField.IsNil() {
		return
	}

	newMap := reflect.MakeMap(mapField.Type())

	key := reflect.New(mapField.Type().Key()).Elem()
	if key.Kind() == reflect.Pointer {
		ptrKey := reflect.New(key.Type())
		fakeStructRecursive(ptrKey.Interface(), currentPtrStruct, depth)
		key.Set(ptrKey)
	} else if key.Kind() == reflect.Struct {
		ptrKey := reflect.New(key.Type())
		fakeStructRecursive(ptrKey.Interface(), currentPtrStruct, depth)
		key.Set(ptrKey.Elem())
	} else {
		key = CreateRandomValue(key)
	}

	value := reflect.New(mapField.Type().Elem()).Elem()
	if value.Kind() == reflect.Pointer {
		ptrKey := reflect.New(value.Type().Elem())
		fakeStructRecursive(ptrKey.Interface(), currentPtrStruct, depth)
		value.Set(ptrKey)
	} else if value.Kind() == reflect.Struct {
		ptrValue := reflect.New(value.Type())
		fakeStructRecursive(ptrValue.Interface(), currentPtrStruct, depth)
		value.Set(ptrValue.Elem())
	} else {
		value = CreateRandomValue(value)
	}

	newMap.SetMapIndex(key, value)

	mapField.Set(newMap)
}
