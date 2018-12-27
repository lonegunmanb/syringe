package codegen

import (
	"github.com/ahmetb/go-linq"
	"github.com/lonegunmanb/syrinx/ast"
)

type productTypeInfoWrap struct {
	ast.TypeInfo
}

func (t *productTypeInfoWrap) GetFields() []ast.FieldInfo {
	fields := t.TypeInfo.GetFields()
	results := make([]ast.FieldInfo, 0, len(fields))
	linq.From(fields).Select(func(fieldInfo interface{}) interface{} {
		return &productFieldInfoWrap{FieldInfo: fieldInfo.(ast.FieldInfo)}
	}).ToSlice(&results)
	return results
}
