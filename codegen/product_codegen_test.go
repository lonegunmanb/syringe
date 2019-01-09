package codegen_test

import (
	"bytes"
	"fmt"
	"github.com/lonegunmanb/syringe/codegen"
	"github.com/lonegunmanb/varys/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

const flyCarCode = `
package fly_car

import (
	"github.com/lonegunmanb/syringe/test_code/car"
	"github.com/lonegunmanb/syringe/test_code/flyer"
)

type FlyCar struct {
	*car.Car %s
	flyer.Plane %s
	Decoration  Decoration %s
}

type Decoration interface {
	LookAndFeel() string
}`

const expected = `package fly_car
import (
    "github.com/lonegunmanb/syringe/ioc"
    "github.com/lonegunmanb/syringe/test_code/car"
    "github.com/lonegunmanb/syringe/test_code/flyer"
)
func Create_FlyCar(container ioc.Container) *FlyCar {
	product := new(FlyCar)
	Assemble_FlyCar(product, container)
	return product
}
func Assemble_FlyCar(product *FlyCar, container ioc.Container) {
	product.Car = container.Resolve("github.com/lonegunmanb/syringe/test_code/car.Car").(*car.Car)
	product.Plane = *container.Resolve("github.com/lonegunmanb/syringe/test_code/flyer.Plane").(*flyer.Plane)
	product.Decoration = container.Resolve("github.com/lonegunmanb/syringe/test_code/fly_car.Decoration").(Decoration)
}
func Register_FlyCar(container ioc.Container) {
	container.RegisterFactory((*FlyCar)(nil), func(container1 ioc.Container) interface{} {
		return Create_FlyCar(container1)
	})
}`
const injectTag = "`inject:\"\"`"

func TestGenerateCreateProductCode(t *testing.T) {
	walker := ast.NewTypeWalker()
	err := walker.Parse("github.com/lonegunmanb/syringe/test_code/fly_car", fmt.Sprintf(flyCarCode, injectTag, injectTag, injectTag))
	assert.Nil(t, err)
	flyCar := walker.GetTypes()[0]
	writer := &bytes.Buffer{}
	gen := codegen.NewProductCodegen(flyCar, writer)
	err = gen.GenerateCode()
	assert.Nil(t, err)
	code := writer.String()
	assert.Equal(t, expected, code)
}
