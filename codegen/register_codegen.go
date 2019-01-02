package codegen

import (
	"github.com/lonegunmanb/syrinx/ast"
	"io"
	"text/template"
)

type RegisterCodegen interface {
	GenerateCode() error
	GetPkgName() string
	GetPkgPath() string
	//GetPkgNameFromPkgPath(pkgPath string) string
}

type registerCodegen struct {
	writer         io.Writer
	pkgName        string
	pkgPath        string
	depPkgPathInfo DepPkgPathInfo
}

func (c *registerCodegen) GetPkgName() string {
	return c.pkgName
}

func (c *registerCodegen) GetPkgPath() string {
	return c.pkgPath
}

func (c *registerCodegen) GetPkgNameFromPkgPath(pkgPath string) string {
	return c.depPkgPathInfo.GetPkgNameFromPkgPath(pkgPath)
}

//noinspection GoUnusedExportedFunction
func NewCodegen(writer io.Writer, typeInfos []ast.TypeInfo, pkgName string, pkgPath string) RegisterCodegen {
	return &registerCodegen{
		writer:         writer,
		depPkgPathInfo: NewDepPkgPathInfo(typeInfos),
		pkgName:        pkgName,
		pkgPath:        pkgPath,
	}
}

func (c *registerCodegen) genPkgDecl() error {
	return genPkgDecl(c.writer, c)
}

func (c *registerCodegen) genImportDecls() error {
	return genImportDecl(c.writer, c.depPkgPathInfo)
}

func (c *registerCodegen) GenerateCode() error {
	return Call(func() error {
		return c.genPkgDecl()
	}).Call(func() error {
		return c.genImportDecls()
	}).Err
}

func (c *registerCodegen) gen(templateName string, text string, data interface{}) (err error) {
	t := template.New(templateName)
	t, err = t.Parse(text)
	if err != nil {
		return
	}
	err = t.Execute(c.writer, data)
	return
}
