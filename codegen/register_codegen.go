package codegen

import (
	"github.com/ahmetb/go-linq"
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
	typeInfos      []Register
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

func (c *registerCodegen) GetRegisters() []Register {
	return c.typeInfos
}

//noinspection GoUnusedExportedFunction
func NewRegisterCodegen(writer io.Writer, typeInfos []ast.TypeInfo, pkgName string, pkgPath string) RegisterCodegen {
	var typeInfosWraps []Register
	pkgPathInfo := NewDepPkgPathInfo(typeInfos)
	linq.From(typeInfos).Select(func(typeInfo interface{}) interface{} {
		return NewTypeInfoWrapWithDepPkgPath(typeInfo.(ast.TypeInfo), pkgPathInfo)
	}).ToSlice(&typeInfosWraps)
	return &registerCodegen{
		writer:         writer,
		depPkgPathInfo: pkgPathInfo,
		pkgName:        pkgName,
		pkgPath:        pkgPath,
		typeInfos:      typeInfosWraps,
	}
}

func (c *registerCodegen) genPkgDecl() error {
	return genPkgDecl(c.writer, c)
}

func (c *registerCodegen) genImportDecls() error {
	return genImportDecl(c.writer, c.depPkgPathInfo)
}

const createIocTemplate = `
func CreateIoc() ioc.Container {
    container := ioc.NewContainer()
{{with .GetRegisters}}{{range .}}    {{.RegisterCode}}
{{end}}{{end}}    return container
}`

func (c *registerCodegen) genRegister() error {
	return gen("registerType", createIocTemplate, c.writer, c)
}

func (c *registerCodegen) GenerateCode() error {
	return Call(func() error {
		return c.genPkgDecl()
	}).Call(func() error {
		return c.genImportDecls()
	}).Call(func() error {
		return c.genRegister()
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
