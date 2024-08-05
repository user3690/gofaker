package faker

import (
	"maps"
	"reflect"
	"slices"
	"testing"
	"time"
)

type EasyStruct struct {
	Id      uint64
	Name    string
	Price   int64
	Average float64
	Valid   bool
}

type ComplexStruct struct {
	Id             uint64
	Name           string
	SubStruct      ComplexSubStruct
	MapStruct      map[uint8]EasyStruct
	MapString      map[uint8]string
	SliceStruct    []EasyStruct
	SlicePrimitive []string
}

func (c ComplexStruct) DeepCopy() ComplexStruct {
	subStruct := c.SubStruct.DeepCopy()
	mapStruct := maps.Clone(c.MapStruct)
	mapString := maps.Clone(c.MapString)
	sliceStruct := slices.Clone(c.SliceStruct)
	slicePrimitive := slices.Clone(c.SlicePrimitive)

	return ComplexStruct{
		Id:             c.Id,
		Name:           c.Name,
		SubStruct:      subStruct,
		MapStruct:      mapStruct,
		MapString:      mapString,
		SliceStruct:    sliceStruct,
		SlicePrimitive: slicePrimitive,
	}
}

type ComplexSubStruct struct {
	Id            int
	SubName       string
	Dict          map[int]string
	Clock         time.Time
	ComplexSlice  []*ComplexSliceStruct
	PointerStruct *PointerStruct
}

func (cs ComplexSubStruct) DeepCopy() ComplexSubStruct {
	// gut für einfache Maps mit primitiven Typen
	newDict := maps.Clone(cs.Dict)

	newComplexSlice := make([]*ComplexSliceStruct, len(cs.ComplexSlice))
	for i := range cs.ComplexSlice {
		entry := cs.ComplexSlice[i]
		newComplexSlice[i] = entry.DeepCopy()
	}

	return ComplexSubStruct{
		Id:            cs.Id,
		SubName:       cs.SubName,
		Dict:          newDict,
		Clock:         cs.Clock,
		ComplexSlice:  newComplexSlice,
		PointerStruct: cs.PointerStruct.DeepCopy(),
	}
}

type ComplexSliceStruct struct {
	Value        string
	IsValid      bool
	ComplexSlice []*ComplexSliceStruct
}

func (cs ComplexSliceStruct) DeepCopy() *ComplexSliceStruct {
	var newComplexSlice []*ComplexSliceStruct
	if cs.ComplexSlice != nil {
		newComplexSlice = make([]*ComplexSliceStruct, len(cs.ComplexSlice))
		for i := range cs.ComplexSlice {
			entry := cs.ComplexSlice[i]
			newComplexSlice[i] = entry.DeepCopy()
		}
	}

	return &ComplexSliceStruct{
		Value:        cs.Value,
		IsValid:      cs.IsValid,
		ComplexSlice: newComplexSlice,
	}
}

type PointerStruct struct {
	Value         string
	PointerStruct *PointerStruct
}

func (p PointerStruct) DeepCopy() *PointerStruct {
	var newPointerStruct *PointerStruct
	if p.PointerStruct != nil {
		newPointerStruct = p.PointerStruct.DeepCopy()
	}

	return &PointerStruct{
		Value:         p.Value,
		PointerStruct: newPointerStruct,
	}
}

func TestFakeStruct(t *testing.T) {
	var complexStruct ComplexStruct

	FakeStruct(&complexStruct, 5)
	complexStructCopy := complexStruct.DeepCopy()
	if !reflect.DeepEqual(complexStruct, complexStructCopy) {
		t.Errorf("Something wrong while creating random struct")
	}
}
