package codegen

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/lonegunmanb/syringe/ast"
)

type TypeInfoWrap interface {
	GetName() string
	GetPkgName() string
	GetDepPkgPaths(fieldTagFilter string) []string
	GetFieldAssigns() []Assembler
	GetEmbeddedTypeAssigns() []Assembler
	GetPkgNameFromPkgPath(pkgPath string) string
	GenImportDecls() []string
	RegisterCode() string
}

type typeInfoWrap struct {
	ast.TypeInfo
	depPkgPathInfo DepPkgPathInfo
}

func (t *typeInfoWrap) RegisterCode() string {
	return fmt.Sprintf("%s.Register_%s(container)", t.depPkgPathInfo.GetPkgNameFromPkgPath(t.GetPkgPath()), t.GetName())
}

func NewTypeInfoWrap(typeInfo ast.TypeInfo) TypeInfoWrap {
	return NewTypeInfoWrapWithDepPkgPath(typeInfo, NewDepPkgPathInfo([]ast.TypeInfo{typeInfo}, typeInfo.GetPkgPath()))
}

func NewTypeInfoWrapWithDepPkgPath(typeInfo ast.TypeInfo, depPkgPathInfo DepPkgPathInfo) TypeInfoWrap {
	return &typeInfoWrap{
		TypeInfo:       typeInfo,
		depPkgPathInfo: depPkgPathInfo,
	}
}

func (t *typeInfoWrap) GenImportDecls() []string {
	return t.depPkgPathInfo.GenImportDecls()
}

func (t *typeInfoWrap) GetPkgNameFromPkgPath(pkgPath string) string {
	return t.depPkgPathInfo.GetPkgNameFromPkgPath(pkgPath)
}

func (t *typeInfoWrap) GetPkgName() string {
	return t.TypeInfo.GetPkgName()
}

func (t *typeInfoWrap) GetFieldAssigns() []Assembler {
	fields := t.TypeInfo.GetFields()
	results := make([]Assembler, 0, len(fields))
	linq.From(fields).Select(func(fieldInfo interface{}) interface{} {
		return &productFieldInfoWrap{FieldInfo: fieldInfo.(ast.FieldInfo), typeInfo: t}
	}).Where(func(fieldInfo interface{}) bool {
		return hasInjectTag(fieldInfo.(ast.FieldInfo).GetTag())
	}).ToSlice(&results)
	return results
}

func (t *typeInfoWrap) GetEmbeddedTypeAssigns() []Assembler {
	embeddedTypes := t.TypeInfo.GetEmbeddedTypes()
	results := make([]Assembler, 0, len(embeddedTypes))
	linq.From(embeddedTypes).Select(func(embeddedType interface{}) interface{} {
		return &productEmbeddedTypeWrap{EmbeddedType: embeddedType.(ast.EmbeddedType), typeInfoWrap: t}
	}).ToSlice(&results)
	return results
}
