package codegen

import (
	"github.com/lonegunmanb/syrinx/ast"
	"html/template"
	"io"
)

const pkgDecl = `package {{.PkgName}}`

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
	typeInfo *ast.TypeInfo
	writer   io.Writer
}

func newCodegen(t *ast.TypeInfo, writer io.Writer) *codegen {
	return &codegen{t, writer}
}

func (c *codegen) genPkgDecl() (err error) {
	t := template.New("pkg")
	t, err = t.Parse(pkgDecl)
	if err != nil {
		return
	}
	err = t.Execute(c.writer, c.typeInfo)
	return
}
