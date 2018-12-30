//WARNING!!!!
//DO NOT REFORMAT THIS CODE
//CODE TEMPLATES AND UNIT TESTS ARE FRAGILE!!!!!
package codegen

import (
	"github.com/lonegunmanb/syrinx/ast"
	"io"
	"text/template"
)

//`
//package {{.PkgName}}
//
//import (
//	"github.com/lonegunmanb/syrinx/ioc"
//	"github.com/lonegunmanb/syrinx/test_code/engine"
//)
//
//func Create_FlyCar(container ioc.Container) *FlyCar {
//	product := new(FlyCar)
//	Assemble_FlyCar(product, container)
//	return product
//}
//+++++++++++++++++++++++++++++++++++++++++++++++++++++
//func Assemble_FlyCar(product *FlyCar, container ioc.Container) {
//	product.Car = container.Resolve("github.com/lonegunmanb/syrinx/test_code/car.Car").(*car.Car)
//	product.Plane = *container.Resolve("github.com/lonegunmanb/syrinx/test_code/flyer.Plane").(*flyer.Plane)
//	product.Decoration = container.Resolve("github.com/lonegunmanb/syrinx/test_code/fly_car.Decoration").(Decoration)
//}
//
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
	codegen
	typeInfo TypeCodegen
}

func newProductCodegen(t ast.TypeInfo, writer io.Writer) *productCodegen {
	return &productCodegen{codegen: codegen{writer}, typeInfo: &productTypeInfoWrap{TypeInfo: t}}
}

func NewProductCodegen(t ast.TypeInfo, writer io.Writer) ProductCodegen {
	return newProductCodegen(t, writer)
}

func (c *productCodegen) Writer() io.Writer {
	return c.writer
}

func (c *productCodegen) GenerateCode() error {
	return Call(func() error {
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

const pkgDecl = `package {{.GetPkgName}}`

func (c *productCodegen) genPkgDecl() error {
	return c.gen("pkg", pkgDecl)
}

const importDecl = `
import (
    "github.com/lonegunmanb/syrinx/ioc"
{{with .GetDepPkgPaths}}{{range .}}    "{{.}}"
{{end}}{{end}})`

func (c *productCodegen) genImportDecls() error {
	return c.gen("imports", importDecl)
}

const createFuncDecl = `
func Create_{{.GetPkgName}}_{{.GetName}}(container ioc.Container) *{{.GetName}} {
	product := new({{.GetName}})
	Assemble_{{.GetName}}(product, container)
	return product
}`

func (c *productCodegen) genCreateFuncDecl() error {
	return c.gen("createFunc", createFuncDecl)
}

const assembleFuncDecl = `
func Assemble_{{.GetPkgName}}_{{.GetName}}(product *{{.GetName}}, container ioc.Container) {
{{with .GetEmbeddedTypeAssigns}}{{range .}}	{{.AssembleCode}}
{{end}}{{end}}{{with .GetFieldAssigns}}{{range .}}	{{.AssembleCode}}
{{end}}{{end}}}`

func (c *productCodegen) genAssembleFuncDecl() error {
	return c.gen("assembleFunc", assembleFuncDecl)
}

const registerFuncDecl = `
func Register(container ioc.Container) {
	container.RegisterFactory((*{{.GetName}})(nil), func(ioc ioc.Container) interface{} {
		return Create_{{.GetName}}(ioc)
	})
}`

func (c *productCodegen) genRegisterFuncDecl() (err error) {
	return c.gen("registerFunc", registerFuncDecl)
}

func (c *productCodegen) gen(templateName string, text string) (err error) {
	t := template.New(templateName)
	t, err = t.Parse(text)
	if err != nil {
		return
	}
	err = t.Execute(c.writer, c.typeInfo)
	return
}
