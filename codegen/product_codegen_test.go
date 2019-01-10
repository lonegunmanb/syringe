package codegen_test

import (
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
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
	pkgPath := "github.com/lonegunmanb/syringe/test_code/fly_car"
	walker := createTypeWalker(t, pkgPath)
	err := walker.Parse(pkgPath, fmt.Sprintf(flyCarCode, injectTag, injectTag, injectTag))
	assert.Nil(t, err)
	flyCar := walker.GetTypes()[0]
	writer := &bytes.Buffer{}
	gen := codegen.NewProductCodegen(flyCar, writer)
	err = gen.GenerateCode()
	assert.Nil(t, err)
	code := writer.String()
	assert.Equal(t, expected, code)
}

func createTypeWalker(t *testing.T, pkgPath string) ast.TypeWalker {
	ctrl := gomock.NewController(t)
	mockOsEnv := codegen.NewMockGoPathEnv(ctrl)
	mockOsEnv.EXPECT().GetPkgPath(gomock.Any()).Times(1).Return(pkgPath, nil)
	ast.RegisterType((*ast.GoPathEnv)(nil), func() interface{} {
		return mockOsEnv
	})
	walker := ast.NewTypeWalker()
	return walker
}
