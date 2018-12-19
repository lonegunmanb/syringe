package syrinx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go/types"
	"testing"
)

func TestWalkFieldInfos(t *testing.T) {
	var sourceCode = `
package test
import "go/ast"
type Struct1 struct {
	Field1 int
	Field2 Struct2
	Field3 *int
	Field4, Field5 float64
	Field6 *ast.Decl
}

type Struct2 struct {
	
}
`
	structWalker := NewStructWalker()
	err := structWalker.Parse(sourceCode)
	assert.Nil(t, err)
	structs := structWalker.Structs()
	struct1 := structs[0]
	assert.Equal(t, 6, len(struct1.Fields))
	field1 := struct1.Fields[0]
	assert.Equal(t, "Field1", field1.Name)
	assert.Equal(t, "int", field1.Type.String())
	field2 := struct1.Fields[1]
	assert.Equal(t, "Field2", field2.Name)
	namedType, ok := field2.Type.(*types.Named)
	assert.True(t, ok)
	assert.Equal(t, "Struct2", namedType.Obj().Name())
	struct2Type, ok := namedType.Underlying().(*types.Struct)
	assert.True(t, ok)
	assert.Equal(t, structs[1].Type, struct2Type)
	field3 := struct1.Fields[2]
	pointer, ok := field3.Type.(*types.Pointer)
	assert.True(t, ok)
	assert.Equal(t, "int", pointer.Elem().String())
	field4 := struct1.Fields[3]
	field5 := struct1.Fields[4]
	assert.Equal(t, field4.Type, field5.Type)
	float64Type, ok := field4.Type.(*types.Basic)
	assert.True(t, ok)
	assert.Equal(t, "float64", float64Type.String())
	field6 := struct1.Fields[5]
	assert.Equal(t, "*go/ast.Decl", field6.Type.String())
}

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
