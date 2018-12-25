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
//func Create_car(container ioc.Container) *Car {
//	r := new(Car)
//	r.Engine = container.Resolve("github.com/lonegunmanb/syrinx/test_code/engine.Engine").(engine.Engine)
//	return r
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
{{with .GetDepPkgPaths}}import (
{{range .}}"{{.}}"
{{end}}){{end}}`

func (c *codegen) genImportDecls() (err error) {
	return c.gen("imports", importDecl)
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
