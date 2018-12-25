package ast

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"go/types"
	"reflect"
)

type TypeInfo struct {
	Name          string
	PkgPath       string
	PkgName       string
	Fields        []*FieldInfo
	Kind          reflect.Kind
	Type          types.Type
	EmbeddedTypes []*EmbeddedType
}

func (typeInfo *TypeInfo) FullName() string {
	if typeInfo.PkgPath == "" {
		return typeInfo.Name
	}
	return fmt.Sprintf("%s.%s", typeInfo.PkgPath, typeInfo.Name)
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
