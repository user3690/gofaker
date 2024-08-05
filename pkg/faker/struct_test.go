package faker

import (
	"maps"
	"reflect"
	"runtime"
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

func (e EasyStruct) DeepCopy() EasyStruct {
	return EasyStruct{
		Id:      e.Id,
		Name:    e.Name,
		Price:   e.Price,
		Average: e.Average,
		Valid:   e.Valid,
	}
}

type ComplexStruct struct {
	Id             uint64
	Name           string
	SubStruct      ComplexSubStruct
	MapPointer     map[uint8]*EasyStruct
	MapStruct      map[uint8]EasyStruct
	MapString      map[uint8]string
	MapStructKey   map[MapKey]string
	SliceStruct    []EasyStruct
	SlicePrimitive []string
}

type MapKey struct {
	Id1 int64
	Id2 int64
}

func (c ComplexStruct) DeepCopy() ComplexStruct {
	subStruct := c.SubStruct.DeepCopy()
	mapStruct := maps.Clone(c.MapStruct)
	mapString := maps.Clone(c.MapString)
	mapStructKey := maps.Clone(c.MapStructKey)
	sliceStruct := slices.Clone(c.SliceStruct)
	slicePrimitive := slices.Clone(c.SlicePrimitive)

	var newMapPointer map[uint8]*EasyStruct
	if c.MapPointer != nil {
		newMapPointer = make(map[uint8]*EasyStruct, len(c.MapPointer))
		for k, v := range c.MapPointer {
			newEasyStruct := v.DeepCopy()
			newMapPointer[k] = &newEasyStruct
		}
	}

	return ComplexStruct{
		Id:             c.Id,
		Name:           c.Name,
		SubStruct:      subStruct,
		MapPointer:     newMapPointer,
		MapStruct:      mapStruct,
		MapString:      mapString,
		MapStructKey:   mapStructKey,
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

func BenchmarkDeepCopy(b *testing.B) {
	var complexStruct, complexStructCopy ComplexStruct
	FakeStruct(&complexStruct, 5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		complexStructCopy = complexStruct.DeepCopy()
	}

	runtime.KeepAlive(complexStructCopy)
}

func BenchmarkFakeStruct(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var complexStruct ComplexStruct
		FakeStruct(&complexStruct, 5)
	}
}
