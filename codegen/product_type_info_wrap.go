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
}

type productTypeInfoWrap struct {
	ast.TypeInfo
}

func (t *productTypeInfoWrap) GetFieldAssigns() []Assembler {
	fields := t.TypeInfo.GetFields()
	results := make([]Assembler, 0, len(fields))
	linq.From(fields).Select(func(fieldInfo interface{}) interface{} {
		return &productFieldInfoWrap{FieldInfo: fieldInfo.(ast.FieldInfo)}
	}).ToSlice(&results)
	return results
}

func (t *productTypeInfoWrap) GetEmbeddedTypeAssigns() []Assembler {
	embeddedTypes := t.TypeInfo.GetEmbeddedTypes()
	results := make([]Assembler, 0, len(embeddedTypes))
	linq.From(embeddedTypes).Select(func(embeddedType interface{}) interface{} {
		return &productEmbeddedTypeWrap{EmbeddedType: embeddedType.(ast.EmbeddedType)}
	}).ToSlice(&results)
	return results
}
