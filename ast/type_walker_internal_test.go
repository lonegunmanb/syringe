package ast

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go/types"
	"reflect"
	"testing"
)

func TestFuncDecl(t *testing.T) {
	sourceCode := `
package ast
type Struct struct {
}
func Test(input int) int{
	i := input
	print(input)
	print(i)
	return 1
}
`
	typeWalker := parseCode(t, sourceCode)
	struct1 := typeWalker.Types()[0]
	assert.Equal(t, 0, len(struct1.Fields))
}

func TestStructWithInterface(t *testing.T) {
	sourceCode := `
package ast
type TestStruct struct {
}
type TestInterface interface {
	Hello()
}
`
	typeWalker := parseCode(t, sourceCode)
	assert.Equal(t, 2, len(typeWalker.Types()))
	testStruct := typeWalker.Types()[0]
	assert.Equal(t, 0, len(testStruct.Fields))
	assert.Equal(t, reflect.Struct, testStruct.Kind)
	testInterface := typeWalker.Types()[1]
	assert.Equal(t, 1, len(testInterface.Fields))
	assert.Equal(t, reflect.Interface, testInterface.Kind)
}

type TestStruct struct {
	Field1 int
	Field2 TestInterface
}
type TestInterface interface {
	Hello()
}

func TestAstTypeShouldEqualToReflectedType(t *testing.T) {
	sourceCode := `
package ast
type TestStruct struct {
	Field1 int
	Field2 TestInterface
}
type TestInterface interface {
	Hello()
}
`
	typeWalker := parseCode(t, sourceCode)
	testStruct := typeWalker.Types()[0]
	testStructInstance := TestStruct{}
	for i := 0; i < 2; i++ {
		assertFieldTypeNameEqual(t, testStruct, testStructInstance, i)
	}
}

func assertFieldTypeNameEqual(t *testing.T, testStruct *TypeInfo, testStructInstance TestStruct, fieldIndex int) {
	astFieldTypeName := testStruct.Fields[fieldIndex].Type.String()
	reflectedFieldType := reflect.TypeOf(testStructInstance).Field(fieldIndex).Type
	pkgPath := reflectedFieldType.PkgPath()
	reflectedTypeName := reflectedFieldType.Name()
	if notEmpty(pkgPath) {
		reflectedTypeName = fmt.Sprintf("%s.%s", pkgPath, reflectedTypeName)
	}
	assert.Equal(t, reflectedTypeName, astFieldTypeName)
}

func notEmpty(pkgPath string) bool {
	return pkgPath != ""
}

func TestFieldTag(t *testing.T) {
	souceCode := `
package ast
type Struct struct {
	Field2 int ` + "`" + "inject:\"Field2\"`" + `
}
`
	typeWalker := parseCode(t, souceCode)
	struct1 := typeWalker.Types()[0]
	field := struct1.Fields[0]
	tag := field.Tag
	assert.Equal(t, "`inject:\"Field2\"`", tag)
}

func TestNewTypeDefinition(t *testing.T) {
	sourceCode := `
package ast
type newint int
type Struct struct {
	Field newint
	Field2 int
}
`
	typeWalker := parseCode(t, sourceCode)
	struct1 := typeWalker.Types()[0]
	field1 := struct1.Fields[0]
	field2 := struct1.Fields[1]
	namedType, ok := field1.Type.(*types.Named)
	assert.True(t, ok)
	assert.Equal(t, "newint", namedType.Obj().Name())
	assert.Equal(t, field2.Type, namedType.Underlying())
}

func TestTypeAliasIsIdenticalToType(t *testing.T) {
	sourceCode := `
package ast
type newint = int
type Struct struct {
	Field1 newint
	Field2 int
}
`
	typeWalker := parseCode(t, sourceCode)
	struct1 := typeWalker.Types()[0]
	field1 := struct1.Fields[0]
	field2 := struct1.Fields[1]
	assert.Equal(t, field1.Type, field2.Type)
}

func TestTypeFromImportWithDot(t *testing.T) {
	sourceCode := `
package ast
import . "go/ast"
type Struct1 struct {
	Field *Decl
}
`
	typeWalker := parseCode(t, sourceCode)
	structs := typeWalker.Types()
	struct1 := structs[0]
	field := struct1.Fields[0]
	assert.Equal(t, "*go/ast.Decl", field.Type.String())
}

func TestWalkFieldInfos(t *testing.T) {
	sourceCode := `
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
	typeWalker := parseCode(t, sourceCode)
	structs := typeWalker.Types()
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

	typeWalker := parseCode(t, structDefine1)
	structs := typeWalker.Types()
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

	typeWalker := parseCode(t, structDefine)
	walkedTypes := typeWalker.Types()
	assert.Len(t, walkedTypes, 4)
	structInfo := walkedTypes[0]
	assert.Equal(t, "NestedStruct", structInfo.Name)
	structInfo = walkedTypes[1]
	assert.Equal(t, "struct{MiddleStructField struct{Field1 int; Field2 string; BottomStructField struct{Field1 int; Field2 string}}}", structInfo.Name)
	structInfo = walkedTypes[2]
	assert.Equal(t, "struct{Field1 int; Field2 string; BottomStructField struct{Field1 int; Field2 string}}", structInfo.Name)
	structInfo = walkedTypes[3]
	assert.Equal(t, "struct{Field1 int; Field2 string}", structInfo.Name)
}

func TestWalkStructWithNestedInterface(t *testing.T) {
	const structDefine = `
package test
type Struct struct {
	Field interface {
		Name() string
	}
}
`
	typeWalker := parseCode(t, structDefine)
	walkedTypes := typeWalker.Types()
	assert.Equal(t, 2, len(walkedTypes))
	nestedType := walkedTypes[1]
	assert.Equal(t, "interface{Name() string}", nestedType.Name)
}

func TestWalkStructInheritingAnotherStruct(t *testing.T) {
	const structDefine = `
package test
type Interface interface {
	Hello()
}
type Struct1 struct {
	FullName string
}
type Struct2 struct {
	Age int
}
type Struct3 struct {
	Interface
	Struct1
	*Struct2
}
`
	typeWalker := parseCode(t, structDefine)
	interface1 := typeWalker.Types()[0]
	struct1 := typeWalker.Types()[1]
	struct2 := typeWalker.Types()[2]
	struct3 := typeWalker.Types()[3]
	assert.Equal(t, 3, len(struct3.EmbeddedTypes))
	assertType(t, interface1, interface1.FullName(), EmbeddedByInterface, struct3.EmbeddedTypes[0])
	assertType(t, struct1, struct1.FullName(), EmbeddedByStruct, struct3.EmbeddedTypes[1])
	assertType(t, struct2, "*"+struct2.FullName(), EmbeddedByPointer, struct3.EmbeddedTypes[2])
}

func assertType(t *testing.T, typeInfo *TypeInfo, fullName string, expectedKind EmbeddedKind, embeddedType *EmbeddedType) {
	assert.Equal(t, fullName, embeddedType.FullName)
	assert.Equal(t, typeInfo.PkgPath, embeddedType.PkgPath)
	assert.Equal(t, expectedKind, embeddedType.Kind)
	assert.Equal(t, "", embeddedType.Tag)
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
