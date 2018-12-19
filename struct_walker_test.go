package syrinx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWalkStructNames(t *testing.T) {
	const structName1 = "Struct"
	const structName2 = "Struct2"
	const field1Name = "MaleOne"
	const field2Name = "Field2"
	var structDefine1 = fmt.Sprintf(`
package test
type %s struct {
	%s int
	%s string
}
type %s struct{
		%s string
		%s int
}
`, structName1,
		field1Name,
		field2Name,

		structName2,
		field1Name,
		field2Name)

	structWalker := NewStructWalker()
	err := structWalker.Parse(structDefine1)
	assert.Nil(t, err)
	structs := structWalker.Structs()
	assert.Len(t, structs, 2)
	structInfo := structs[0]
	assert.Equal(t, structName1, structInfo.Name)
	structInfo = structs[1]
	assert.Equal(t, structName2, structInfo.Name)
}

func TestWalkStructWithNestedStruct(t *testing.T) {
	const structDefine = `
package test
type NestedStruct struct {
	TopStructField struct {
		MiddleStructField struct {
			Field1 int
			Field2 string
			BottomStructField struct {
				Field1 int
				Field2 string
			}
		}
	}
}
`
	structWalker := NewStructWalker()
	err := structWalker.Parse(structDefine)
	assert.Nil(t, err)
	structs := structWalker.Structs()
	assert.Len(t, structs, 4)
	structInfo := structs[0]
	assert.Equal(t, "NestedStruct", structInfo.Name)
	structInfo = structs[1]
	assert.Equal(t, "struct{MiddleStructField struct{Field1 int; Field2 string; BottomStructField struct{Field1 int; Field2 string}}}", structInfo.Name)
	structInfo = structs[2]
	assert.Equal(t, "struct{Field1 int; Field2 string; BottomStructField struct{Field1 int; Field2 string}}", structInfo.Name)
	structInfo = structs[3]
	assert.Equal(t, "struct{Field1 int; Field2 string}", structInfo.Name)
}

type Struct struct {
	Field struct {
		Field1 string
		Field2 int
		Field3 struct {
			Field1 int
			Field2 string
		}
	}
}

type Struct2 struct {
	Field struct {
		Field1 string
		Field2 int
		Field3 struct {
			Field1 int
			Field2 string
		}
	}
}

func TestDifferentNamedStructCannotConvert(t *testing.T) {
	assert.Panics(t, func() {
		s1 := Struct{
			Field: struct {
				Field1 string
				Field2 int
				Field3 struct {
					Field1 int
					Field2 string
				}
			}{
				Field1: "1",
				Field2: 1,
				Field3: struct {
					Field1 int
					Field2 string
				}{Field1: 1, Field2: "1"},
			},
		}
		returnField1(s1)
	})
}

func TestSubStructsWithSameStructureAreIdentical(t *testing.T) {
	s1 := Struct{
		Field: struct {
			Field1 string
			Field2 int
			Field3 struct {
				Field1 int
				Field2 string
			}
		}{
			Field1: "1",
			Field2: 1,
			Field3: struct {
				Field1 int
				Field2 string
			}{Field1: 1, Field2: "1"},
		},
	}
	s2 := Struct2{}
	s2.Field = s1.Field
	assert.Equal(t, 1, s2.Field.Field3.Field1)
	assert.Equal(t, "1", s2.Field.Field3.Field2)
}

func returnField1(input interface{}) int {
	return input.(Struct2).Field.Field2
}
