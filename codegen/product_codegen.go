//WARNING!!!!
//DO NOT REFORMAT THIS CODE
//CODE TEMPLATES AND UNIT TESTS ARE FRAGILE!!!!!
package codegen

import (
	"github.com/ahmetb/go-linq"
	"github.com/lonegunmanb/syrinx/ast"
	"html/template"
	"io"
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

type productTypeInfoWrap struct {
	ast.TypeInfo
}

type productFieldInfoWrap struct {
	ast.FieldInfo
}

func (t *productTypeInfoWrap) GetFields() []ast.FieldInfo {
	fields := t.TypeInfo.GetFields()
	results := make([]ast.FieldInfo, 0, len(fields))
	linq.From(fields).Select(func(fieldInfo interface{}) interface{} {
		return &productFieldInfoWrap{FieldInfo: fieldInfo.(ast.FieldInfo)}
	}).ToSlice(&results)
	return results
}

type productCodegen struct {
	codegen
	typeInfo *productTypeInfoWrap
}

func newProductCodegen(t ast.TypeInfo, writer io.Writer) *productCodegen {
	return &productCodegen{codegen: codegen{writer}, typeInfo: &productTypeInfoWrap{TypeInfo: t}}
}

const pkgDecl = `package {{.GetPkgName}}`

func (c *productCodegen) genPkgDecl() (err error) {
	return c.gen("pkg", pkgDecl)
}

const importDecl = `
import (
"github.com/lonegunmanb/syrinx/ioc"
{{with .GetDepPkgPaths}}{{range .}}"{{.}}"
{{end}}{{end}})`

func (c *productCodegen) genImportDecls() (err error) {
	return c.gen("imports", importDecl)
}

const createFuncDecl = `
func Create_{{.GetName}}(container ioc.Container) *{{.GetName}} {
	product := new({{.GetName}})
	Assemble_{{.GetName}}(product, container)
	return product
}`

func (c *productCodegen) genCreateFuncDecl() (err error) {
	return c.gen("createFunc", createFuncDecl)
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
