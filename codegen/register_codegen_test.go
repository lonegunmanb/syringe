package codegen_test

import (
	"bytes"
	"github.com/ahmetb/go-linq"
	"github.com/lonegunmanb/syringe/codegen"
	"github.com/lonegunmanb/varys/ast"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

const modelACode = `
package model

type Request struct {
	Name string
}
`

const modelBCode = modelACode

const engineCode = `
package engine

type HorsePower = int

type Engine interface {
	OutputPower() HorsePower
}

type pistonEngine struct {
	horsePower HorsePower
}

func (e *pistonEngine) OutputPower() HorsePower {
	return e.horsePower
}

func NewEngine(horsePower HorsePower) Engine {
	return &pistonEngine{horsePower: horsePower}
}
`

const carCode = `
package car

import "github.com/lonegunmanb/syringe/test_code/engine"

type Car struct {
	Engine engine.Engine
}`

const flyerCode = `
package flyer

type Wing interface {
	Material() string
}

type AluminumWing struct {
}

func (*AluminumWing) Material() string {
	return "aluminum"
}

type Plane struct {
	Wing Wing
}`

const fly_carCode = `
package fly_car

import (
	"github.com/lonegunmanb/syringe/test_code/car"
	model1 "github.com/lonegunmanb/syringe/test_code/check_package_name_duplicate_a/model"
	model2 "github.com/lonegunmanb/syringe/test_code/check_package_name_duplicate_b/model"
	"github.com/lonegunmanb/syringe/test_code/flyer"
)

type FlyCar struct {
	*car.Car    
	flyer.Plane 
	Decoration  Decoration 
	R1          *model1.Request
	R2          *model2.Request
}

type Decoration interface {
	LookAndFeel() string
}

type FancyDecoration struct {
}

func (f *FancyDecoration) LookAndFeel() string {
	return "Fancy"
}
`

const expectedRegisterCode = `package test_code
import (
    "github.com/lonegunmanb/syringe/ioc"
    p0 "github.com/lonegunmanb/syringe/test_code/check_package_name_duplicate_a/model"
    p1 "github.com/lonegunmanb/syringe/test_code/check_package_name_duplicate_b/model"
    "github.com/lonegunmanb/syringe/test_code/car"
    "github.com/lonegunmanb/syringe/test_code/engine"
    "github.com/lonegunmanb/syringe/test_code/flyer"
    "github.com/lonegunmanb/syringe/test_code/fly_car"
)
func CreateIoc() ioc.Container {
    container := ioc.NewContainer()
    RegisterCodeWriter(container)
    return container
}
func RegisterCodeWriter(container ioc.Container) {
    p0.Register_Request(container)
    p1.Register_Request(container)
    car.Register_Car(container)
    engine.Register_pistonEngine(container)
    flyer.Register_AluminumWing(container)
    flyer.Register_Plane(container)
    fly_car.Register_FlyCar(container)
    fly_car.Register_FancyDecoration(container)
}`

func TestGenRegisterCode(t *testing.T) {
	typeInfos := parseTypeInfos(t, "github.com/lonegunmanb/syringe/test_code/check_package_name_duplicate_a/model", modelACode)
	typeInfos = append(typeInfos, parseTypeInfos(t, "github.com/lonegunmanb/syringe/test_code/check_package_name_duplicate_b/model", modelBCode)...)
	typeInfos = append(typeInfos, parseTypeInfos(t, "github.com/lonegunmanb/syringe/test_code/car", carCode)...)
	typeInfos = append(typeInfos, parseTypeInfos(t, "github.com/lonegunmanb/syringe/test_code/engine", engineCode)...)
	typeInfos = append(typeInfos, parseTypeInfos(t, "github.com/lonegunmanb/syringe/test_code/flyer", flyerCode)...)
	typeInfos = append(typeInfos, parseTypeInfos(t, "github.com/lonegunmanb/syringe/test_code/fly_car", fly_carCode)...)
	var structTypes []ast.TypeInfo
	linq.From(typeInfos).Where(func(t interface{}) bool {
		return t.(ast.TypeInfo).GetKind() == reflect.Struct
	}).ToSlice(&structTypes)
	writer := &bytes.Buffer{}
	sut := codegen.NewRegisterCodegen(writer, structTypes, "test_code", "github.com/lonegunmanb/syringe/test_code")
	err := sut.GenerateCode()
	assert.Nil(t, err)
	actual := writer.String()
	assert.Equal(t, expectedRegisterCode, actual)
}

func parseTypeInfos(t *testing.T, pkgPath string, code string) []ast.TypeInfo {
	typeWalker := createTypeWalker(t, pkgPath)
	err := typeWalker.Parse(pkgPath, code)
	assert.Nil(t, err)
	return typeWalker.GetTypes()
}
