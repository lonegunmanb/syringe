package codegen

import (
	"fmt"
	"github.com/lonegunmanb/varys/ast"
	"go/types"
)

func getDeclType(pkgPath string, t types.Type, qf func(*types.Package) string) string {
	switch t.(type) {
	case *types.Basic:
		{
			return t.String()
		}
	case *types.Named:
		{
			return parseNamesType(t.(*types.Named), pkgPath, qf)
		}
	case *types.Struct:
		{
			return types.TypeString(t, qf)
		}
	case *types.Interface:
		{
			panic("do not support nested interface field injection temporary")
		}
	case *types.Pointer:
		{
			return fmt.Sprintf("*%s", getDeclType(pkgPath, t.(*types.Pointer).Elem(), qf))
		}
	case *types.Slice:
		{
			return fmt.Sprintf("[]%s", getDeclType(pkgPath, t.(*types.Slice).Elem(), qf))
		}
	case *types.Array:
		{
			return fmt.Sprintf("[%d]%s", t.(*types.Array).Len(), getDeclType(pkgPath, t.(*types.Array).Elem(), qf))
		}
	case *types.Map:
		{
			mapType := t.(*types.Map)
			return parseMapType(mapType, pkgPath, qf)
		}
	case *types.Chan:
		{
			chanDecl := parseChannelType(t.(*types.Chan))
			return fmt.Sprintf(chanDecl, getDeclType(pkgPath, t.(*types.Chan).Elem(), qf))
		}
	case *types.Signature:
		{
			return types.TypeString(t, qf)
		}
	default:
		panic(fmt.Sprintf("unsupported type %s", t.String()))
	}
}

func parseNamesType(namedType *types.Named, pkgPath string, qf func(*types.Package) string) string {
	pkg := ast.GetNamedTypePkg(namedType)
	if pkg == nil {
		panic(fmt.Sprintf("cannot find package name for %s", namedType.String()))
	}
	typeName := namedType.Obj().Name()
	path := pkg.Path()
	if samePackage(path, pkgPath) {
		return typeName
	} else {
		return types.TypeString(namedType, qf)
	}
}

func parseMapType(mapType *types.Map, pkgPath string, qf func(*types.Package) string) string {
	keyDecl := getDeclType(pkgPath, mapType.Key(), qf)
	valueDecl := getDeclType(pkgPath, mapType.Elem(), qf)
	return fmt.Sprintf("map[%s]%s", keyDecl, valueDecl)
}

func parseChannelType(chanType *types.Chan) string {
	switch chanType.Dir() {
	case types.SendRecv:
		{
			return "chan %s"
		}
	case types.SendOnly:
		{
			return "chan<- %s"
		}
	default:
		{
			return "<-chan %s"
		}
	}
}

func samePackage(path string, pkgPath string) bool {
	return path == pkgPath
}
