package codegen

import (
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
//func Assemble_FlyCar(l *FlyCar, container ioc.Container) {
//	l.Car = new(car.Car)
//	car.Assemble_car(l.Car, container)
//	flyer.Assemble_plane(&l.Plane, container)
//	l.Decoration = container.Resolve("github.com/lonegunmanb/syrinx/test_code/fly_car.Decoration").(Decoration)
//}
//`
type codegen struct {
	typeInfo ast.TypeInfo
	writer   io.Writer
}

func newCodegen(t ast.TypeInfo, writer io.Writer) *codegen {
	return &codegen{t, writer}
}

const pkgDecl = `package {{.GetPkgName}}`

func (c *codegen) genPkgDecl() (err error) {
	return c.gen("pkg", pkgDecl)
}

const importDecl = `
import (
"github.com/lonegunmanb/syrinx/ioc"
{{with .GetDepPkgPaths}}{{range .}}"{{.}}"
{{end}}{{end}})`

func (c *codegen) genImportDecls() (err error) {
	return c.gen("imports", importDecl)
}

const createFuncDecl = `
func Create_{{.GetName}}(container ioc.Container) *{{.GetName}} {
	product := new({{.GetName}})
	Assemble_{{.GetName}}(product, container)
	return product
}`

func (c *codegen) genCreateFuncDecl() (err error) {
	return c.gen("createFunc", createFuncDecl)
}

func (c *codegen) gen(templateName string, text string) (err error) {
	t := template.New(templateName)
	t, err = t.Parse(text)
	if err != nil {
		return
	}
	err = t.Execute(c.writer, c.typeInfo)
	return
}
