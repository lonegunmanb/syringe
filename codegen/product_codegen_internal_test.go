package codegen

//go:generate mockgen -package=codegen -destination=./mock_type_info.go github.com/lonegunmanb/syringe/ast TypeInfo
//go:generate mockgen -package=codegen -destination=./mock_field_info.go github.com/lonegunmanb/syringe/ast FieldInfo
//go:generate mockgen -package=codegen -destination=./mock_embedded_type.go github.com/lonegunmanb/syringe/ast EmbeddedType
//go:generate mockgen -package=codegen -destination=./mock_assembler.go github.com/lonegunmanb/syringe/codegen Assembler
//go:generate mockgen -package=codegen -destination=./mock_type_codegen.go github.com/lonegunmanb/syringe/codegen TypeInfoWrap
import (
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/lonegunmanb/syringe/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

const expectedPackageDecl = "package ast"

func TestGenPackageDecl(t *testing.T) {
	testProductGen(t, func(mockTypeCodegen *MockTypeInfoWrap) {
		mockTypeCodegen.EXPECT().GetPkgName().Times(1).Return("ast")
	}, func(gen *productCodegen) error {
		return gen.genPkgDecl()
	}, expectedPackageDecl)
}

const expectedImportDecl = `
import (
    "github.com/lonegunmanb/syringe/ioc"
    "go/ast"
    "go/token"
    "go/types"
)`

func TestGenImportsDecl(t *testing.T) {
	testProductGen(t, func(mockTypeCodegen *MockTypeInfoWrap) {
		depImports := []string{
			`"go/ast"`,
			`"go/token"`,
			`"go/types"`,
		}
		setupMockToGenImports(mockTypeCodegen, depImports)
	}, func(gen *productCodegen) error {
		return gen.genImportDecls()
	}, expectedImportDecl)
}

func TestShouldNotGenExtraImportsIfDepPathsEmpty(t *testing.T) {
	testProductGen(t, func(mockTypeCodegen *MockTypeInfoWrap) {
		var depImports []string
		setupMockToGenImports(mockTypeCodegen, depImports)
	}, func(gen *productCodegen) error {
		return gen.genImportDecls()
	}, `
import (
    "github.com/lonegunmanb/syringe/ioc"
)`)
}

func TestGenImportsDeclForDuplicatePkgName(t *testing.T) {
	testProductGen(t, func(mockTypeCodegen *MockTypeInfoWrap) {
		setupMockToGenImports(mockTypeCodegen, []string{
			`p0 "a/pkg"`,
			`p1 "b/pkg"`,
		})
	}, func(gen *productCodegen) error {
		return gen.genImportDecls()
	}, `
import (
    "github.com/lonegunmanb/syringe/ioc"
    p0 "a/pkg"
    p1 "b/pkg"
)`)
}

func TestGenCreateFuncDecl(t *testing.T) {
	const expectedFlyCarCreateCode = `
func Create_FlyCar(container ioc.Container) *FlyCar {
	product := new(FlyCar)
	Assemble_FlyCar(product, container)
	return product
}`
	testProductGen(t, func(typeInfo *MockTypeInfoWrap) {
		typeInfo.EXPECT().GetName().Times(4).Return("FlyCar")
	}, func(gen *productCodegen) error {
		r := gen.genCreateFuncDecl()
		return r
	}, expectedFlyCarCreateCode)
}

func TestGenCreateFuncDeclForCustomIdent(t *testing.T) {
	const expectedFlyCarCreateCode = `
func Create_FlyCar(c ioc.Container) *FlyCar {
	product := new(FlyCar)
	Assemble_FlyCar(product, c)
	return product
}`
	originIdent := ContainerIdentName
	ContainerIdentName = "c"
	defer func() {
		ContainerIdentName = originIdent
	}()
	testProductGen(t, func(typeInfo *MockTypeInfoWrap) {
		typeInfo.EXPECT().GetName().Times(4).Return("FlyCar")
	}, func(gen *productCodegen) error {
		r := gen.genCreateFuncDecl()
		return r
	}, expectedFlyCarCreateCode)
}

const expectedFlyCarAssembleCode = `
func Assemble_FlyCar(product *FlyCar, container ioc.Container) {
	product.Car = container.Resolve("github.com/lonegunmanb/syringe/test_code/car.Car").(*car.Car)
	product.Plane = *container.Resolve("github.com/lonegunmanb/syringe/test_code/flyer.Plane").(*flyer.Plane)
	product.Decoration = container.Resolve("github.com/lonegunmanb/syringe/test_code/fly_car.Decoration").(Decoration)
}`

func TestGenAssembleFuncDecl(t *testing.T) {
	testProductGen(t, func(typeInfo *MockTypeInfoWrap) {
		embeddedCarMock := NewMockAssembler(typeInfo.ctrl)
		embeddedCarMock.EXPECT().AssembleCode().Times(1).Return(`product.Car = container.Resolve("github.com/lonegunmanb/syringe/test_code/car.Car").(*car.Car)`)
		embeddedPlaneMock := NewMockAssembler(typeInfo.ctrl)
		embeddedPlaneMock.EXPECT().AssembleCode().Times(1).Return(`product.Plane = *container.Resolve("github.com/lonegunmanb/syringe/test_code/flyer.Plane").(*flyer.Plane)`)
		typeInfo.EXPECT().GetName().Times(2).Return("FlyCar")
		typeInfo.EXPECT().GetEmbeddedTypeAssigns().Times(1).Return([]Assembler{embeddedCarMock, embeddedPlaneMock})
		decorationMock := NewMockAssembler(typeInfo.ctrl)
		decorationMock.EXPECT().AssembleCode().Times(1).Return(`product.Decoration = container.Resolve("github.com/lonegunmanb/syringe/test_code/fly_car.Decoration").(Decoration)`)
		typeInfo.EXPECT().GetFieldAssigns().Times(1).Return([]Assembler{decorationMock})
	}, func(gen *productCodegen) error {
		r := gen.genAssembleFuncDecl()
		return r
	}, expectedFlyCarAssembleCode)
}

func TestGenAssembleFuncDeclWithCustomIdent(t *testing.T) {
	expectedFlyCarAssembleCode := `
func Assemble_FlyCar(product *FlyCar, c ioc.Container) {
	product.Decoration = c.Resolve("github.com/lonegunmanb/syringe/test_code/fly_car.Decoration").(Decoration)
}`
	originIdent := ContainerIdentName
	ContainerIdentName = "c"
	defer func() {
		ContainerIdentName = originIdent
	}()
	testProductGen(t, func(typeInfo *MockTypeInfoWrap) {
		typeInfo.EXPECT().GetName().Times(2).Return("FlyCar")
		typeInfo.EXPECT().GetEmbeddedTypeAssigns().Times(1).Return([]Assembler{})
		decorationMock := NewMockAssembler(typeInfo.ctrl)
		decorationMock.EXPECT().AssembleCode().Times(1).Return(`product.Decoration = c.Resolve("github.com/lonegunmanb/syringe/test_code/fly_car.Decoration").(Decoration)`)
		typeInfo.EXPECT().GetFieldAssigns().Times(1).Return([]Assembler{decorationMock})
	}, func(gen *productCodegen) error {
		r := gen.genAssembleFuncDecl()
		return r
	}, expectedFlyCarAssembleCode)
}

const actualFlyCarCode = `
package fly_car

import (
	"github.com/lonegunmanb/syringe/test_code/car"
	"github.com/lonegunmanb/syringe/test_code/flyer"
)

type FlyCar struct {
	*car.Car
	flyer.Plane %s
	Decoration  Decoration %s
}

type Decoration interface {
	LookAndFeel() string
}
`
const injectTag = "`inject:\"\"`"

func TestActualAssembleFuncDecl(t *testing.T) {
	walker := ast.NewTypeWalker()

	err := walker.Parse("github.com/lonegunmanb/syringe/test_code/fly_car", fmt.Sprintf(actualFlyCarCode, injectTag, injectTag))
	assert.Nil(t, err)
	flyCar := walker.GetTypes()[0]
	writer := &bytes.Buffer{}
	productCodegen := newProductCodegen(flyCar, writer)
	err = productCodegen.genAssembleFuncDecl()
	assert.Nil(t, err)
	code := writer.String()
	assert.Equal(t, expectedFlyCarAssembleCode, code)
}

const expectedRegisterFuncCode = `
func Register_FlyCar(container ioc.Container) {
	container.RegisterFactory((*FlyCar)(nil), func(container1 ioc.Container) interface{} {
		return Create_FlyCar(container1)
	})
}`

func TestRegisterFuncDecl(t *testing.T) {
	testProductGen(t, func(typeInfo *MockTypeInfoWrap) {
		typeInfo.EXPECT().GetName().Times(3).Return("FlyCar")
	}, func(gen *productCodegen) error {
		r := gen.genRegisterFuncDecl()
		return r
	}, expectedRegisterFuncCode)
}

func TestRegisterFuncDeclWithCustomIdent(t *testing.T) {
	expectedRegisterFuncCode := `
func Register_FlyCar(c ioc.Container) {
	c.RegisterFactory((*FlyCar)(nil), func(c1 ioc.Container) interface{} {
		return Create_FlyCar(c1)
	})
}`
	originIdent := ContainerIdentName
	ContainerIdentName = "c"
	defer func() {
		ContainerIdentName = originIdent
	}()
	testProductGen(t, func(typeInfo *MockTypeInfoWrap) {
		typeInfo.EXPECT().GetName().Times(3).Return("FlyCar")
	}, func(gen *productCodegen) error {
		r := gen.genRegisterFuncDecl()
		return r
	}, expectedRegisterFuncCode)
}

func testProductGen(t *testing.T, setupMockFunc func(info *MockTypeInfoWrap),
	testMethod func(gen *productCodegen) error, expected string) {
	writer := &bytes.Buffer{}
	ctrl, typeInfo := prepareTypeCodegenMock(t)
	defer ctrl.Finish()
	//
	setupMockFunc(typeInfo)
	codegen := &productCodegen{writer: writer, typeInfo: typeInfo}
	err := testMethod(codegen)
	assert.Nil(t, err)
	code := writer.String()
	assert.Equal(t, expected, code)
}

func prepareTypeCodegenMock(t *testing.T) (*gomock.Controller, *MockTypeInfoWrap) {
	ctrl := gomock.NewController(t)
	typeCodegen := NewMockTypeInfoWrap(ctrl)
	return ctrl, typeCodegen
}

func setupMockToGenImports(typeCodegen *MockTypeInfoWrap, expectedImports []string) {
	typeCodegen.EXPECT().GenImportDecls().Times(1).Return(expectedImports)
}
