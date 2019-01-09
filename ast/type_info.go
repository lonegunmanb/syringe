package ast

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"go/ast"
	"go/types"
	"reflect"
	"regexp"
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
	GetDepPkgPaths(fieldTagFilter string) []string
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

func (typeInfo *typeInfo) GetDepPkgPaths(fieldTagFilter string) []string {
	var fieldFilter *regexp.Regexp
	if fieldTagFilter != "" {
		fieldFilter = regexp.MustCompile(fmt.Sprintf("%s:\".*\"", fieldTagFilter))
	}
	result := make([]string, 0)
	linq.From(typeInfo.EmbeddedTypes).Select(
		func(embeddedType interface{}) interface{} {
			return embeddedType.(EmbeddedType).GetPkgPath()
		}).Union(linq.From(typeInfo.Fields).Where(func(fieldInfo interface{}) bool {
		if fieldTagFilter != "" {
			return fieldFilter.MatchString(fieldInfo.(FieldInfo).GetTag())
		}
		return true
	}).SelectMany(
		func(fieldInfo interface{}) linq.Query {
			paths := fieldInfo.(FieldInfo).GetDepPkgPaths()
			return linq.From(paths)
		})).Distinct().Where(func(path interface{}) bool {
		p := path.(string)
		if typeInfo.GetPkgPath() == "" {
			return p != typeInfo.GetPkgName()
		}
		return p != typeInfo.GetPkgPath()
	}).ToSlice(&result)
	return result
}

func (typeInfo *typeInfo) processField(field *ast.Field, fieldType types.Type) {
	if isEmbeddedField(field) {
		typeInfo.addInheritance(field, fieldType)
	} else {
		typeInfo.addFieldInfos(field, fieldType)
	}
}

func (typeInfo *typeInfo) addFieldInfos(field *ast.Field, fieldType types.Type) {
	names := field.Names
	for _, fieldName := range names {
		typeInfo.Fields = append(typeInfo.Fields, &fieldInfo{
			Name:          fieldName.Name,
			Type:          fieldType,
			Tag:           getTag(field),
			ReferenceFrom: typeInfo,
		})
	}
}

func (typeInfo *typeInfo) addInheritance(field *ast.Field, fieldType types.Type) {
	var kind EmbeddedKind
	var packagePath string
	switch t := fieldType.(type) {
	case *types.Named:
		{
			if isStructType(t) {
				kind = EmbeddedByStruct
			} else {
				kind = EmbeddedByInterface
			}
			packagePath = GetNamedTypePkg(t).Path()
		}
	case *types.Pointer:
		{
			elemType, ok := t.Elem().(*types.Named)
			if !ok {
				panic(fmt.Sprintf("unknown embedded type %s", fieldType.String()))
			}
			kind = EmbeddedByPointer
			packagePath = GetNamedTypePkg(elemType).Path()
		}
	default:
		panic(fmt.Sprintf("unknown embedded type %s", t.String()))
	}

	embeddedType := &embeddedType{
		Kind:          kind,
		FullName:      fieldType.String(),
		PkgPath:       packagePath,
		Tag:           getTag(field),
		Type:          fieldType,
		ReferenceFrom: typeInfo,
	}
	typeInfo.EmbeddedTypes = append(typeInfo.EmbeddedTypes, embeddedType)
}
