package ast

import (
	"github.com/ahmetb/go-linq"
	"go/types"
	"reflect"
)

type TypeInfo struct {
	Name    string
	PkgPath string
	Fields  []*FieldInfo
	Type    types.Type
	Kind    reflect.Kind
}

func (typeInfo *TypeInfo) DepPkgPaths() []string {
	result := make([]string, 0)
	linq.From(typeInfo.Fields).SelectMany(
		func(fieldInfo interface{}) linq.Query {
			paths := fieldInfo.(*FieldInfo).DepPkgPaths()
			return linq.From(paths)
		}).
		Distinct().ToSlice(&result)
	return result
}
