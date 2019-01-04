package ast

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"go/types"
	"reflect"
)

type TypeInfo interface {
	GetName() string
	GetPkgPath() string
	GetPkgName() string
	GetPhysicalPath() string
	GetFields() []FieldInfo
	GetKind() reflect.Kind
	GetType() types.Type
	GetEmbeddedTypes() []EmbeddedType
	GetFullName() string
	GetDepPkgPaths() []string
}

type typeInfo struct {
	Name          string
	PkgPath       string
	PkgName       string
	PhysicalPath  string
	Fields        []FieldInfo
	Kind          reflect.Kind
	Type          types.Type
	EmbeddedTypes []EmbeddedType
}

func (typeInfo *typeInfo) GetPhysicalPath() string {
	return typeInfo.PhysicalPath
}

func (typeInfo *typeInfo) GetName() string {
	return typeInfo.Name
}

func (typeInfo *typeInfo) GetPkgPath() string {
	return typeInfo.PkgPath
}

func (typeInfo *typeInfo) GetPkgName() string {
	return typeInfo.PkgName
}

func (typeInfo *typeInfo) GetFields() []FieldInfo {
	return typeInfo.Fields
}

func (typeInfo *typeInfo) GetKind() reflect.Kind {
	return typeInfo.Kind
}

func (typeInfo *typeInfo) GetType() types.Type {
	return typeInfo.Type
}

func (typeInfo *typeInfo) GetEmbeddedTypes() []EmbeddedType {
	return typeInfo.EmbeddedTypes
}

func (typeInfo *typeInfo) GetFullName() string {
	if typeInfo.PkgPath == "" {
		return typeInfo.Name
	}
	return fmt.Sprintf("%s.%s", typeInfo.PkgPath, typeInfo.Name)
}

func (typeInfo *typeInfo) GetDepPkgPaths() []string {
	result := make([]string, 0)
	linq.From(typeInfo.EmbeddedTypes).Select(
		func(embeddedType interface{}) interface{} {
			return embeddedType.(EmbeddedType).GetPkgPath()
		}).Union(linq.From(typeInfo.Fields).SelectMany(
		func(fieldInfo interface{}) linq.Query {
			paths := fieldInfo.(FieldInfo).GetDepPkgPaths()
			return linq.From(paths)
		})).Distinct().Where(func(path interface{}) bool {
		return path.(string) != typeInfo.PkgPath
	}).ToSlice(&result)
	return result
}
