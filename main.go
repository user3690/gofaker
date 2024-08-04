package main

import (
	"fmt"
	"gofaker/pkg/faker"
	"maps"
	"time"
)

type EasyStruct struct {
	Id int
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

	return ComplexStruct{
		Id:        c.Id,
		Name:      c.Name,
		SubStruct: subStruct,
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
	return &ComplexSliceStruct{
		Value:   cs.Value,
		IsValid: cs.IsValid,
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

func main() {
	test := ComplexStruct{
		Id:   1,
		Name: "Test",
		SubStruct: ComplexSubStruct{
			Id:      1,
			SubName: "SubTest",
			Dict: map[int]string{
				0: "Zero",
				1: "One",
				2: "Two",
				3: "Three",
			},
			Clock: time.Now(),
			ComplexSlice: []*ComplexSliceStruct{
				{
					Value:   "Value1",
					IsValid: false,
				},
				{
					Value:   "Value2",
					IsValid: false,
				},
			},
			PointerStruct: &PointerStruct{
				Value: "Value1",
				PointerStruct: &PointerStruct{
					Value:         "Value2",
					PointerStruct: nil,
				},
			},
		},
	}

	testCopy := test.DeepCopy()

	time.Sleep(1 * time.Second)

	testCopy.SubStruct.Dict[0] = "ZeroChanged"
	testCopy.SubStruct.Clock = time.Now()
	testCopy.SubStruct.ComplexSlice[0].IsValid = true
	testCopy.SubStruct.PointerStruct.Value = "PSValue1Changed"
	testCopy.SubStruct.PointerStruct.PointerStruct.Value = "PSValue2Changed"

	fmt.Println(test.SubStruct.Dict[0])
	fmt.Println(testCopy.SubStruct.Dict[0])
	fmt.Println(test.SubStruct.Clock)
	fmt.Println(testCopy.SubStruct.Clock)
	fmt.Println(test.SubStruct.ComplexSlice[0].IsValid)
	fmt.Println(testCopy.SubStruct.ComplexSlice[0].IsValid)
	fmt.Println(test.SubStruct.PointerStruct.Value)
	fmt.Println(testCopy.SubStruct.PointerStruct.Value)
	fmt.Println(test.SubStruct.PointerStruct.PointerStruct.Value)
	fmt.Println(testCopy.SubStruct.PointerStruct.PointerStruct.Value)

	fmt.Println("")
	fmt.Println("########################")
	fmt.Println("")

	testEasy := &EasyStruct{}
	faker.FakeStruct(testEasy, 5)
	fmt.Println(*testEasy)

	fmt.Println("")
	fmt.Println("########################")
	fmt.Println("")

	testComplex := &ComplexStruct{}
	faker.FakeStruct(testComplex, 5)
	fmt.Println(testComplex)
}
