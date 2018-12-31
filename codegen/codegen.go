package codegen

import (
	"github.com/lonegunmanb/syrinx/ast"
	"io"
	"text/template"
)

type Codegen interface {
	GenerateCode() error
}

type codegen struct {
	writer  io.Writer
	genTask GenTask
}

func NewCodegen(writer io.Writer, genTask GenTask) Codegen {
	return &codegen{writer: writer, genTask: genTask}
}

const pkgDecl = `package {{.GetPkgName}}`

func (c *codegen) genPkgDecl() error {
	return c.gen("pkg", pkgDecl)
}

const importDecl = `
import (
    "github.com/lonegunmanb/syrinx/ioc"
{{with .GetDepPkgPaths}}{{range .}}    "{{.}}"
{{end}}{{end}})`

func (c *codegen) genImportDecls() error {
	return c.gen("imports", importDecl)
}

func (c *codegen) GenerateCode() error {
	return Call(func() error {
		return c.genPkgDecl()
	}).Call(func() error {
		return c.genImportDecls()
	}).CallEach(c.genTask.GetTypeInfos(), func(t interface{}) error {
		return NewProductCodegen(t.(ast.TypeInfo), c.writer).GenerateCode()
	}).Err
}

func (c *codegen) gen(templateName string, text string) (err error) {
	t := template.New(templateName)
	t, err = t.Parse(text)
	if err != nil {
		return
	}
	err = t.Execute(c.writer, c.genTask)
	return
}
