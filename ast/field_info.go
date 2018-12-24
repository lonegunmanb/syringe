package ast

import (
	"fmt"
	"go/types"
)

type FieldInfo struct {
	Name   string
	Type   types.Type
	Tag    string
	Parent *TypeInfo
}

func (f *FieldInfo) DepPkgPaths() []string {
	return getDepPkgPaths(f, f.Type)
}

func getDepPkgPaths(fieldInfo *FieldInfo, t types.Type) []string {
	switch t.(type) {
	case *types.Basic:
		{
			return []string{}
		}
	case *types.Named:
		{
			namedType := t.(*types.Named)
			obj := namedType.Obj()
			if obj == nil {
				return []string{}
			}
			pkg := obj.Pkg()
			if pkg == nil {
				return []string{}
			}
			return []string{pkg.Path()}
		}
	case *types.Struct:
		{
			return []string{fieldInfo.Parent.PkgPath}
		}
	case *types.Interface:
		{
			return []string{fieldInfo.Parent.PkgPath}
		}
	case *types.Pointer:
		{
			return getDepPkgPaths(fieldInfo, t.(*types.Pointer).Elem())
		}
	case *types.Slice:
		{
			return getDepPkgPaths(fieldInfo, t.(*types.Slice).Elem())
		}
	case *types.Array:
		{
			return getDepPkgPaths(fieldInfo, t.(*types.Array).Elem())
		}
	case *types.Map:
		{
			mapType := t.(*types.Map)
			keyPaths := getDepPkgPaths(fieldInfo, mapType.Key())
			valuePaths := getDepPkgPaths(fieldInfo, mapType.Elem())
			return append(keyPaths, valuePaths...)
		}
	case *types.Chan:
		{
			chanType := t.(*types.Chan)
			return getDepPkgPaths(fieldInfo, chanType.Elem())
		}
	case *types.Signature:
		{
			funcType := t.(*types.Signature)
			depPaths := make([]string, 0, funcType.Params().Len()+funcType.Results().Len())
			depPaths = append(depPaths, tupleDeps(fieldInfo, funcType.Params())...)
			depPaths = append(depPaths, tupleDeps(fieldInfo, funcType.Results())...)
			return depPaths
		}
	default:
		panic(fmt.Sprintf("unsupported type %s", t.String()))
	}
}

func tupleDeps(fieldInfo *FieldInfo, tuple *types.Tuple) []string {
	depPaths := make([]string, 0, tuple.Len())
	for i := 0; i < tuple.Len(); i++ {
		depPaths = append(depPaths, getDepPkgPaths(fieldInfo, tuple.At(i).Type())...)
	}
	return depPaths
}
