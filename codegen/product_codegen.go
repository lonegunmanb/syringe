//WARNING!!!!
//DO NOT REFORMAT THIS CODE
//CODE TEMPLATES AND UNIT TESTS ARE FRAGILE!!!!!
package codegen

import (
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
func Create_{{.GetName}}(container ioc.Container) *{{.GetName}} {
	product := new({{.GetName}})
	Assemble_{{.GetName}}(product, container)
	return product
}`

func (c *productCodegen) genCreateFuncDecl() error {
	return c.gen("createFunc", createFuncDecl)
}

const assembleFuncDecl = `
func Assemble_{{.GetName}}(product *{{.GetName}}, container ioc.Container) {
{{with .GetEmbeddedTypeAssigns}}{{range .}}	{{.AssembleCode}}
{{end}}{{end}}{{with .GetFieldAssigns}}{{range .}}	{{.AssembleCode}}
{{end}}{{end}}}`

func (c *productCodegen) genAssembleFuncDecl() error {
	return c.gen("assembleFunc", assembleFuncDecl)
}

const registerFuncDecl = `
func Register_{{.GetName}}(container ioc.Container) {
	container.RegisterFactory((*{{.GetName}})(nil), func(ioc ioc.Container) interface{} {
		return Create_{{.GetName}}(ioc)
	})
}`

func (c *productCodegen) genRegisterFuncDecl() (err error) {
	return c.gen("registerFunc", registerFuncDecl)
}

func (c *productCodegen) gen(templateName string, text string) (err error) {
	return gen(templateName, text, c.writer, c.typeInfo)
}
