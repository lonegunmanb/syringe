package codegen

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/lonegunmanb/syringe/util"
	"github.com/lonegunmanb/varys/ast"
	"io"
	"text/template"
)

type RegisterCodegen interface {
	GenerateCode() error
	GetPkgName() string
	GetPkgPath() string
}

type registerCodegen struct {
	writer              io.Writer
	pkgName             string
	pkgPath             string
	pkgNameArbitrator   PkgNameArbitrator
	registerCodeWriters []RegisterCodeWriter
}

func (c *registerCodegen) GetPkgName() string {
	return c.pkgName
}

func (c *registerCodegen) GetPkgPath() string {
	return c.pkgPath
}

func (c *registerCodegen) GetPkgNameFromPkgPath(pkgPath string) string {
	return c.pkgNameArbitrator.GetPkgNameFromPkgPath(pkgPath)
}

func (c *registerCodegen) GetRegisterCodeWriters() []RegisterCodeWriter {
	return c.registerCodeWriters
}

//noinspection GoUnusedExportedFunction
func NewRegisterCodegen(writer io.Writer, typeInfos []ast.TypeInfo, pkgName string, workingPkgPath string) RegisterCodegen {
	var typeInfosWraps []RegisterCodeWriter
	pkgNameArbitrator := newPkgNameArbitrator(typeInfos, workingPkgPath, RegisterCodegenMode)
	linq.From(typeInfos).Select(func(typeInfo interface{}) interface{} {
		typeInfoWrap := NewTypeInfoWrapWithDepPkgPath(typeInfo.(ast.TypeInfo), pkgNameArbitrator)
		return &registerCodeWriter{
			typeInfo:       typeInfoWrap,
			workingPkgPath: workingPkgPath,
		}
	}).ToSlice(&typeInfosWraps)
	return &registerCodegen{
		writer:              writer,
		pkgNameArbitrator:   pkgNameArbitrator,
		pkgName:             pkgName,
		pkgPath:             workingPkgPath,
		registerCodeWriters: typeInfosWraps,
	}
}

func (c *registerCodegen) genPkgDecl() error {
	return genPkgDecl(c.writer, c)
}

func (c *registerCodegen) genImportDecls() error {
	return genImportDecl(c.writer, c.pkgNameArbitrator)
}

const createIocTemplate = `
func CreateIoc() ioc.Container {
    %s := ioc.NewContainer()
    RegisterCodeWriter(%s)
    return %s
}
func RegisterCodeWriter(%s ioc.Container) {
{{with .GetRegisterCodeWriters}}{{range .}}    {{.RegisterCode}}
{{end}}{{end}}}`

func (c *registerCodegen) genRegister() error {
	return gen("registerType", fmt.Sprintf(createIocTemplate, ContainerIdentName, ContainerIdentName, ContainerIdentName, ContainerIdentName), c.writer, c)
}

func (c *registerCodegen) GenerateCode() error {
	return util.Call(func() error {
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
