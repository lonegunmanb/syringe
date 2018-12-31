package codegen_test

import (
	"bytes"
	"github.com/lonegunmanb/syrinx/ast"
	"github.com/lonegunmanb/syrinx/codegen"
	"github.com/stretchr/testify/assert"
	"testing"
)

const flyCarCode = `
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
}`

const expected = `package fly_car
import (
    "github.com/lonegunmanb/syrinx/ioc"
    "github.com/lonegunmanb/syrinx/test_code/car"
    "github.com/lonegunmanb/syrinx/test_code/flyer"
)
func Create_fly_car_FlyCar(container ioc.Container) *FlyCar {
	product := new(FlyCar)
	Assemble_FlyCar(product, container)
	return product
}
func Assemble_fly_car_FlyCar(product *FlyCar, container ioc.Container) {
	product.Car = container.Resolve("github.com/lonegunmanb/syrinx/test_code/car.Car").(*car.Car)
	product.Plane = *container.Resolve("github.com/lonegunmanb/syrinx/test_code/flyer.Plane").(*flyer.Plane)
	product.Decoration = container.Resolve("github.com/lonegunmanb/syrinx/test_code/fly_car.Decoration").(Decoration)
}
func Register(container ioc.Container) {
	container.RegisterFactory((*FlyCar)(nil), func(ioc ioc.Container) interface{} {
		return Create_FlyCar(ioc)
	})
}`

func TestGenerateCreateProductCode(t *testing.T) {
	walker := ast.NewTypeWalker()
	err := walker.Parse("github.com/lonegunmanb/syrinx/test_code/fly_car", flyCarCode)
	assert.Nil(t, err)
	flyCar := walker.GetTypes()[0]
	writer := &bytes.Buffer{}
	gen := codegen.NewCodegen(writer, codegen.NewCodegenTask("fly_car", []ast.TypeInfo{flyCar}))
	err = gen.GenerateCode()
	assert.Nil(t, err)
	code := writer.String()
	assert.Equal(t, expected, code)
}
