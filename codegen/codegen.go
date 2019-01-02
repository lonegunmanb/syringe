package codegen

import (
	"io"
	"text/template"
)

type Codegen interface {
	//GenerateCode() error
	//GetPkgNameFromPkgPath(pkgPath string) string
}

type codegen struct {
	writer         io.Writer
	depPkgPathInfo DepPkgPathInfo
}

func (c *codegen) GetPkgNameFromPkgPath(pkgPath string) string {
	return c.depPkgPathInfo.GetPkgNameFromPkgPath(pkgPath)
}

//noinspection GoUnusedExportedFunction
func NewCodegen(writer io.Writer, depPkgPathInfo DepPkgPathInfo) Codegen {
	return &codegen{writer: writer, depPkgPathInfo: depPkgPathInfo}
}

const pkgDecl = `package {{.GetPkgName}}`

func (c *codegen) genPkgDecl() error {
	return c.gen("pkg", pkgDecl)
}

const importDecl = `
import (
    "github.com/lonegunmanb/syrinx/ioc"
{{with .GenImportDecls}}{{range .}}    {{.}}
{{end}}{{end}})`

func (c *codegen) genImportDecls() error {
	return c.gen("imports", importDecl)
}

//func (c *codegen) GenerateCode() error {
//	return Call(func() error {
//		return c.genPkgDecl()
//	}).Call(func() error {
//		return c.genImportDecls()
//	}).CallEach(c.depPkgPathInfo.GetTypeInfos(), func(t interface{}) error {
//		return NewProductCodegen(t.(ast.TypeInfo), c.writer, c).GenerateCode()
//	}).Err
//}

func (c *codegen) gen(templateName string, text string) (err error) {
	t := template.New(templateName)
	t, err = t.Parse(text)
	if err != nil {
		return
	}
	err = t.Execute(c.writer, c.depPkgPathInfo)
	return
}
