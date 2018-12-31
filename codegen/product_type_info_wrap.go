package codegen

import (
	"github.com/ahmetb/go-linq"
	"github.com/lonegunmanb/syrinx/ast"
)

type TypeCodegen interface {
	GetName() string
	GetPkgName() string
	GetDepPkgPaths() []string
	GetFieldAssigns() []Assembler
	GetEmbeddedTypeAssigns() []Assembler
	GetPkgNameFromPkgPath(pkgPath string) string
}

type typeInfoWrap struct {
	ast.TypeInfo
	codegen ProductCodegen
}

func (t *typeInfoWrap) GetPkgNameFromPkgPath(pkgPath string) string {
	return t.codegen.GetPkgNameFromPkgPath(pkgPath)
}

func (t *typeInfoWrap) GetPkgName() string {
	return t.codegen.GetPkgNameFromPkgPath(t.GetPkgPath())
}

func (t *typeInfoWrap) GetFieldAssigns() []Assembler {
	fields := t.TypeInfo.GetFields()
	results := make([]Assembler, 0, len(fields))
	linq.From(fields).Select(func(fieldInfo interface{}) interface{} {
		return &productFieldInfoWrap{FieldInfo: fieldInfo.(ast.FieldInfo)}
	}).ToSlice(&results)
	return results
}

func (t *typeInfoWrap) GetEmbeddedTypeAssigns() []Assembler {
	embeddedTypes := t.TypeInfo.GetEmbeddedTypes()
	results := make([]Assembler, 0, len(embeddedTypes))
	linq.From(embeddedTypes).Select(func(embeddedType interface{}) interface{} {
		return &productEmbeddedTypeWrap{EmbeddedType: embeddedType.(ast.EmbeddedType), typeCodegen: t}
	}).ToSlice(&results)
	return results
}
