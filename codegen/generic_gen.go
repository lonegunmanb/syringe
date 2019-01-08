package codegen

import (
	"io"
	"text/template"
)

type WithPackageName interface {
	GetPkgName() string
}

type WithImportDecls interface {
	GenImportDecls() []string
}

const pkgDecl = `package {{.GetPkgName}}`

func genPkgDecl(writer io.Writer, data WithPackageName) error {
	return gen("pkg", pkgDecl, writer, data)
}

const importDecl = `
import (
    "github.com/lonegunmanb/syringe/ioc"
{{with .GenImportDecls}}{{range .}}    {{.}}
{{end}}{{end}})`

func genImportDecl(writer io.Writer, importDecls WithImportDecls) error {
	return gen("imports", importDecl, writer, importDecls)
}

func gen(templateName string, text string, writer io.Writer, data interface{}) (err error) {
	t := template.New(templateName)
	t, err = t.Parse(text)
	if err != nil {
		return
	}
	err = t.Execute(writer, data)
	return
}

var ContainerIdentName = "container"
var ProductIdentName = "product"
