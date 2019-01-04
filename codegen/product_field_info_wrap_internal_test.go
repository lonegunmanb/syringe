package codegen

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/lonegunmanb/syrinx/ast"
	"github.com/stretchr/testify/assert"
	"go-funk"
	"go/types"
	"testing"
)

const fieldAssignTestTemplate = `
package test
%s
type Struct struct {
Field %s
}
type Struct2 struct {
}
`

func TestGetBasicTypeDecl(t *testing.T) {
	testTypeDecl(t, "", "int")
}

func TestAssembleCodeForBasicType(t *testing.T) {
	field := getField(t, "", "int")
	sut := &productFieldInfoWrap{
		FieldInfo: field,
		typeInfo:  NewTypeInfoWrap(field.GetReferenceFrom()),
	}
	code := sut.AssembleCode()
	assert.Equal(t, "product.Field = container.Resolve(\"int\").(int)", code)
}

func TestGetNamedTypeDecl(t *testing.T) {
	testTypeDecls(t, []funk.Tuple{
		{"import \"go/ast\"", "ast.Node"},
		{"", "Struct2"},
	})
	testTypeDecl(t, "import \"go/ast\"", "ast.Node")
}

func TestGetPointerTypeDecl(t *testing.T) {
	testTypeDecls(t, []funk.Tuple{
		{"", "*Struct2"},
		{"import \"go/ast\"", "*ast.Node"},
	})
}

func TestGetSliceTypeDecl(t *testing.T) {
	testTypeDecls(t, []funk.Tuple{
		{"", "[]Struct2"},
		{"", "[]*Struct2"},
		{"import \"go/ast\"", "[]ast.Node"},
		{"import \"go/ast\"", "[]*ast.Node"},
	})
}

func TestGetArrayTypeDecl(t *testing.T) {
	testTypeDecls(t, []funk.Tuple{
		{"", "[1]Struct2"},
		{"", "[1]*Struct2"},
		{"import \"go/ast\"", "[1]ast.Node"},
		{"import \"go/ast\"", "[1]*ast.Node"},
	})
}

func TestGetMapTypeDecl(t *testing.T) {
	testTypeDecls(t, []funk.Tuple{
		{"", "map[int]Struct2"},
		{"", "map[Struct2]int"},
		{"", "map[*Struct2]int"},
		{"", "map[int]*Struct2"},
		{"import \"go/ast\"", "map[ast.Node]ast.Node"},
		{"import \"go/ast\"", "map[ast.Node]Struct2"},
		{"import \"go/ast\"", "map[*ast.Node]Struct2"},
		{"import \"go/ast\"", "map[*ast.Node]*Struct2"},
		{"import \"go/ast\"", "map[ast.Node]*Struct2"},
		{"import \"go/ast\"", "map[ast.Node]*ast.Node"},
		{"import \"go/ast\"", "map[*ast.Node]*ast.Node"},
		{"import \"go/ast\"", "map[*ast.Node]ast.Node"},
		{"import \"go/ast\"", "map[Struct2]ast.Node"},
		{"import \"go/ast\"", "map[*Struct2]ast.Node"},
		{"import \"go/ast\"", "map[*Struct2]*ast.Node"},
		{"import \"go/ast\"", "map[Struct2]*ast.Node"},
		{"import (\"go/ast\"\n\"go/types\")", "map[types.Type]ast.Node"},
		{"import (\"go/ast\"\n\"go/types\")", "map[*types.Type]ast.Node"},
		{"import (\"go/ast\"\n\"go/types\")", "map[types.Type]*ast.Node"},
		{"import (\"go/ast\"\n\"go/types\")", "map[*types.Type]*ast.Node"},
	})
}

func TestGetChanTypeDecl(t *testing.T) {

	testTypeDecls(t, []funk.Tuple{
		{"", "chan Struct2"},
		{"", "chan *Struct2"},
		{"", "chan<- Struct2"},
		{"", "chan<- *Struct2"},
		{"", "<-chan Struct2"},
		{"", "<-chan *Struct2"},
		{"import \"go/ast\"", "chan ast.Node"},
		{"import \"go/ast\"", "chan *ast.Node"},
		{"import \"go/ast\"", "chan<- ast.Node"},
		{"import \"go/ast\"", "chan<- *ast.Node"},
		{"import \"go/ast\"", "<-chan ast.Node"},
		{"import \"go/ast\"", "<-chan *ast.Node"},
	})
}

func TestGetNestedStructTypeDecl(t *testing.T) {
	testTypeDecls(t, []funk.Tuple{
		{"", "struct{Field2 struct{Name string}}"},
		{"", "*struct{Name string}"},
		{"", "*struct{Field2 *struct{Name string}}"},
		{"import \"go/ast\"", "struct{Field2 struct{Node ast.Node}}"},
		{"import \"go/ast\"", "struct{Node *ast.Node}"},
	})
}

func TestGetNestedInterfaceTypeDeclShouldPanic(t *testing.T) {
	assert.Panics(t, func() {
		testTypeDecl(t, "", "interface{Name() string}")
	})
}

func TestGetFuncTypeDecl(t *testing.T) {
	testTypeDecl(t, "import (\"go/ast\"\n\"go/types\")",
		"func(node ast.Node) types.Package")
}

const flyCarCode = `
package fly_car

import (
	"github.com/lonegunmanb/syrinx/test_code/car"
	"github.com/lonegunmanb/syrinx/test_code/flyer"
)

type FlyCar struct {
	*car.Car
	flyer.Plane
	D Decoration
	C car.Car
	P *flyer.Plane
}

type Decoration interface {
	LookAndFeel() string
}
`

func TestGenerateAssembleCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTypeCodegen := NewMockTypeInfoWrap(ctrl)
	mockTypeCodegen.EXPECT().GetPkgNameFromPkgPath("github.com/lonegunmanb/syrinx/test_code/fly_car").Return("fly_car")
	mockTypeCodegen.EXPECT().GetPkgNameFromPkgPath("github.com/lonegunmanb/syrinx/test_code/car").Return("car")
	mockTypeCodegen.EXPECT().GetPkgNameFromPkgPath("github.com/lonegunmanb/syrinx/test_code/flyer").Return("flyer")
	walker := ast.NewTypeWalker()
	err := walker.Parse("github.com/lonegunmanb/syrinx/test_code/fly_car", flyCarCode)
	assert.Nil(t, err)
	flyCar := walker.GetTypes()[0]
	decorationField := &productFieldInfoWrap{flyCar.GetFields()[0], mockTypeCodegen}
	assert.Equal(t, `product.D = container.Resolve("github.com/lonegunmanb/syrinx/test_code/fly_car.Decoration").(Decoration)`,
		decorationField.AssembleCode())
	carField := &productFieldInfoWrap{flyCar.GetFields()[1], mockTypeCodegen}
	assert.Equal(t, `product.C = *container.Resolve("github.com/lonegunmanb/syrinx/test_code/car.Car").(*car.Car)`,
		carField.AssembleCode())
	planeField := &productFieldInfoWrap{flyCar.GetFields()[2], mockTypeCodegen}
	assert.Equal(t, `product.P = container.Resolve("github.com/lonegunmanb/syrinx/test_code/flyer.Plane").(*flyer.Plane)`,
		planeField.AssembleCode())
}

func TestGetInjectKeyFromTag(t *testing.T) {
	assert.Equal(t, "expected", getKeyFromTag("`json:\"name\" xml:\"name\" inject:\"expected\"`"))
}

func testTypeDecls(t *testing.T, args []funk.Tuple) {
	for _, tuple := range args {
		testTypeDecl(t, tuple.Element1.(string), tuple.Element2.(string))
	}
}

func testTypeDecl(t *testing.T, importString string, typeString string) {
	field := getField(t, importString, typeString)
	fieldType := field.GetType()
	pkgPath := field.GetReferenceFrom().GetPkgPath()
	assert.Equal(t, typeString, getDeclType(pkgPath, fieldType, func(p *types.Package) string {
		return p.Name()
	}))
}

func getField(t *testing.T, importString string, typeString string) ast.FieldInfo {
	code := fmt.Sprintf(fieldAssignTestTemplate, importString, typeString)
	walker := parseCode(t, code)
	struct1 := walker.GetTypes()[0]
	return struct1.GetFields()[0]
}
