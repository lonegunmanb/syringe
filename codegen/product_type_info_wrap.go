package codegen

import (
	"github.com/ahmetb/go-linq"
	"github.com/lonegunmanb/varys/ast"
)

type TypeInfoWrap interface {
	GetName() string
	GetPkgName() string
	GetPkgPath() string
	GetDepPkgPaths(fieldTagFilter string) []string
	GetFieldAssigns() []Assembler
	GetEmbeddedTypeAssigns() []Assembler
	GetPkgNameFromPkgPath(pkgPath string) string
	GenImportDecls() []string
}

type typeInfoWrap struct {
	ast.TypeInfo
	pkgNameArbitrator PkgNameArbitrator
}

func NewTypeInfoWrap(typeInfo ast.TypeInfo) TypeInfoWrap {
	return NewTypeInfoWrapWithDepPkgPath(typeInfo,
		newPkgNameArbitrator([]ast.TypeInfo{typeInfo}, typeInfo.GetPkgPath(), ProductCodegenMode))
}

func NewTypeInfoWrapWithDepPkgPath(typeInfo ast.TypeInfo, pkgNameArbitrator PkgNameArbitrator) TypeInfoWrap {
	return &typeInfoWrap{
		TypeInfo:          typeInfo,
		pkgNameArbitrator: pkgNameArbitrator,
	}
}

func (t *typeInfoWrap) GenImportDecls() []string {
	return t.pkgNameArbitrator.GenImportDecls()
}

func (t *typeInfoWrap) GetPkgNameFromPkgPath(pkgPath string) string {
	return t.pkgNameArbitrator.GetPkgNameFromPkgPath(pkgPath)
}

func (t *typeInfoWrap) GetPkgName() string {
	return t.TypeInfo.GetPkgName()
}

func (t *typeInfoWrap) GetPkgPath() string {
	return t.TypeInfo.GetPkgPath()
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
	}).Where(func(embeddedType interface{}) bool {
		return hasInjectTag(embeddedType.(ast.EmbeddedType).GetTag())
	}).ToSlice(&results)
	return results
}
