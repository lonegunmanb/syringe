//WARNING!!!!
//DO NOT REFORMAT THIS CODE
//CODE TEMPLATES AND UNIT TESTS ARE FRAGILE!!!!!
package codegen

import (
	"fmt"
	"github.com/lonegunmanb/syringe/ast"
	"github.com/lonegunmanb/syringe/util"
	"io"
)

//`
//package {{.PkgName}}
//
//import (
//	"github.com/lonegunmanb/syringe/ioc"
//	"github.com/lonegunmanb/syringe/test_code/engine"
//)
//func Create_FlyCar(container ioc.Container) *FlyCar {
//	product := new(FlyCar)
//	Assemble_FlyCar(product, container)
//	return product
//}
//func Assemble_FlyCar(product *FlyCar, container ioc.Container) {
//	product.Car = container.Resolve("github.com/lonegunmanb/syringe/test_code/car.Car").(*car.Car)
//	product.Plane = *container.Resolve("github.com/lonegunmanb/syringe/test_code/flyer.Plane").(*flyer.Plane)
//	product.Decoration = container.Resolve("github.com/lonegunmanb/syringe/test_code/fly_car.Decoration").(Decoration)
//}
//func Register_FlyCar(container ioc.Container) {
//	container.RegisterFactory((*FlyCar)(nil), func(ioc ioc.Container) interface{} {
//		return Create_FlyCar(ioc)
//	})
//}
//`

type ProductCodegen interface {
	GenerateCode() error
	Writer() io.Writer
}

type productCodegen struct {
	writer   io.Writer
	typeInfo TypeInfoWrap
}

func newProductCodegen(t ast.TypeInfo, writer io.Writer) *productCodegen {
	return &productCodegen{writer: writer, typeInfo: NewTypeInfoWrap(t)}
}

func NewProductCodegen(t ast.TypeInfo, writer io.Writer) ProductCodegen {
	return newProductCodegen(t, writer)
}

func (c *productCodegen) Writer() io.Writer {
	return c.writer
}

func (c *productCodegen) GenerateCode() error {
	return util.Call(func() error {
		return c.genPkgDecl()
	}).Call(func() error {
		return c.genImportDecls()
	}).Call(func() error {
		return c.genCreateFuncDecl()
	}).Call(func() error {
		return c.genAssembleFuncDecl()
	}).Call(func() error {
		return c.genRegisterFuncDecl()
	}).Err
}

func (c *productCodegen) genPkgDecl() error {
	return genPkgDecl(c.writer, c.typeInfo)
}

func (c *productCodegen) genImportDecls() error {
	return genImportDecl(c.writer, c.typeInfo)
}

const createFuncDecl = `
func Create_{{.GetName}}(%s ioc.Container) *{{.GetName}} {
	%s := new({{.GetName}})
	Assemble_{{.GetName}}(%s, %s)
	return %s
}`

func (c *productCodegen) genCreateFuncDecl() error {
	return c.gen("createFunc", fmt.Sprintf(createFuncDecl, ContainerIdentName, ProductIdentName, ProductIdentName, ContainerIdentName, ProductIdentName))
}

const assembleFuncDecl = `
func Assemble_{{.GetName}}(%s *{{.GetName}}, %s ioc.Container) {
{{with .GetEmbeddedTypeAssigns}}{{range .}}	{{.AssembleCode}}
{{end}}{{end}}{{with .GetFieldAssigns}}{{range .}}	{{.AssembleCode}}
{{end}}{{end}}}`

func (c *productCodegen) genAssembleFuncDecl() error {
	return c.gen("assembleFunc", fmt.Sprintf(assembleFuncDecl, ProductIdentName, ContainerIdentName))
}

const registerFuncDecl = `
func Register_{{.GetName}}(%s ioc.Container) {
	%s.RegisterFactory((*{{.GetName}})(nil), func(%s1 ioc.Container) interface{} {
		return Create_{{.GetName}}(%s1)
	})
}`

func (c *productCodegen) genRegisterFuncDecl() (err error) {
	return c.gen("registerFunc", fmt.Sprintf(registerFuncDecl, ContainerIdentName, ContainerIdentName, ContainerIdentName, ContainerIdentName))
}

func (c *productCodegen) gen(templateName string, text string) (err error) {
	return gen(templateName, text, c.writer, c.typeInfo)
}
