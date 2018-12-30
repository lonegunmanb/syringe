package codegen

//go:generate mockgen -package=codegen -destination=./mock_type_info.go github.com/lonegunmanb/syrinx/ast TypeInfo
//go:generate mockgen -package=codegen -destination=./mock_field_info.go github.com/lonegunmanb/syrinx/ast FieldInfo
//go:generate mockgen -package=codegen -destination=./mock_embedded_type.go github.com/lonegunmanb/syrinx/ast EmbeddedType
//go:generate mockgen -package=codegen -destination=./mock_assembler.go github.com/lonegunmanb/syrinx/codegen Assembler
//go:generate mockgen -package=codegen -destination=./mock_type_codegen.go github.com/lonegunmanb/syrinx/codegen TypeCodegen
import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/lonegunmanb/syrinx/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenImportsDecl(t *testing.T) {
	testProductGen(t, func(typeInfo *MockTypeCodegen) {
		depImports := []string{
			"go/ast",
			"go/token",
			"go/types",
		}
		setupMockToGenImports(typeInfo, depImports)
	}, func(gen *productCodegen) error {
		return gen.genImportDecls()
	}, `
import (
    "github.com/lonegunmanb/syrinx/ioc"
    "go/ast"
    "go/token"
    "go/types"
)`)
}

func TestShouldNotGenExtraImportsIfDepPathsEmpty(t *testing.T) {
	testProductGen(t, func(typeInfo *MockTypeCodegen) {
		var depImports []string
		setupMockToGenImports(typeInfo, depImports)
	}, func(gen *productCodegen) error {
		return gen.genImportDecls()
	}, `
import (
    "github.com/lonegunmanb/syrinx/ioc"
)`)
}

const expectedFlyCarCreateCode = `
func Create_fly_car_FlyCar(container ioc.Container) *FlyCar {
	product := new(FlyCar)
	Assemble_FlyCar(product, container)
	return product
}`

func TestGenCreateFuncDecl(t *testing.T) {
	testProductGen(t, func(typeInfo *MockTypeCodegen) {
		typeInfo.EXPECT().GetName().Times(4).Return("FlyCar")
		typeInfo.EXPECT().GetPkgName().Times(1).Return("fly_car")
	}, func(gen *productCodegen) error {
		r := gen.genCreateFuncDecl()
		return r
	}, expectedFlyCarCreateCode)
}

const expectedFlyCarAssembleCode = `
func Assemble_fly_car_FlyCar(product *FlyCar, container ioc.Container) {
	product.Car = container.Resolve("github.com/lonegunmanb/syrinx/test_code/car.Car").(*car.Car)
	product.Plane = *container.Resolve("github.com/lonegunmanb/syrinx/test_code/flyer.Plane").(*flyer.Plane)
	product.Decoration = container.Resolve("github.com/lonegunmanb/syrinx/test_code/fly_car.Decoration").(Decoration)
}`

func TestGenAssembleFuncDecl(t *testing.T) {
	testProductGen(t, func(typeInfo *MockTypeCodegen) {
		embeddedCarMock := NewMockAssembler(typeInfo.ctrl)
		embeddedCarMock.EXPECT().AssembleCode().Times(1).Return(`product.Car = container.Resolve("github.com/lonegunmanb/syrinx/test_code/car.Car").(*car.Car)`)
		embeddedPlaneMock := NewMockAssembler(typeInfo.ctrl)
		embeddedPlaneMock.EXPECT().AssembleCode().Times(1).Return(`product.Plane = *container.Resolve("github.com/lonegunmanb/syrinx/test_code/flyer.Plane").(*flyer.Plane)`)
		typeInfo.EXPECT().GetName().Times(2).Return("FlyCar")
		typeInfo.EXPECT().GetEmbeddedTypeAssigns().Times(1).Return([]Assembler{embeddedCarMock, embeddedPlaneMock})
		typeInfo.EXPECT().GetPkgName().Times(1).Return("fly_car")
		decorationMock := NewMockAssembler(typeInfo.ctrl)
		decorationMock.EXPECT().AssembleCode().Times(1).Return(`product.Decoration = container.Resolve("github.com/lonegunmanb/syrinx/test_code/fly_car.Decoration").(Decoration)`)
		typeInfo.EXPECT().GetFieldAssigns().Times(1).Return([]Assembler{decorationMock})
	}, func(gen *productCodegen) error {
		r := gen.genAssembleFuncDecl()
		return r
	}, expectedFlyCarAssembleCode)
}

const actualFlyCarCode = `
package fly_car

import (
	"github.com/lonegunmanb/syrinx/test_code/car"
	"github.com/lonegunmanb/syrinx/test_code/flyer"
)

type FlyCar struct {
	*car.Car
	flyer.Plane
	Decoration  Decoration
}

type Decoration interface {
	LookAndFeel() string
}
`

func TestActualAssembleFuncDecl(t *testing.T) {
	walker := ast.NewTypeWalker()
	err := walker.Parse("github.com/lonegunmanb/syrinx/test_code/fly_car", actualFlyCarCode)
	assert.Nil(t, err)
	flyCar := walker.GetTypes()[0]
	writer := &bytes.Buffer{}
	codegen := newProductCodegen(flyCar, writer)
	err = codegen.genAssembleFuncDecl()
	assert.Nil(t, err)
	code := writer.String()
	assert.Equal(t, expectedFlyCarAssembleCode, code)
}

const expectedRegisterFuncCode = `
func Register(container ioc.Container) {
	container.RegisterFactory((*FlyCar)(nil), func(ioc ioc.Container) interface{} {
		return Create_FlyCar(ioc)
	})
}`

func TestRegisterFuncDecl(t *testing.T) {
	testProductGen(t, func(typeInfo *MockTypeCodegen) {
		typeInfo.EXPECT().GetName().Times(2).Return("FlyCar")
	}, func(gen *productCodegen) error {
		r := gen.genRegisterFuncDecl()
		return r
	}, expectedRegisterFuncCode)
}

func testProductGen(t *testing.T, setupMockFunc func(info *MockTypeCodegen),
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

func setupMockToGenImports(typeInfo *MockTypeCodegen, depImports []string) {
	typeInfo.EXPECT().GetDepPkgPaths().Times(1).Return(depImports)
}

func prepareTypeCodegenMock(t *testing.T) (*gomock.Controller, *MockTypeCodegen) {
	ctrl := gomock.NewController(t)
	typeCodegen := NewMockTypeCodegen(ctrl)
	return ctrl, typeCodegen
}
