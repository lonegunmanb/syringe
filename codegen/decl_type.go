package codegen

import (
	"fmt"
	"github.com/lonegunmanb/syringe/ast"
	"go/types"
)

func getDeclType(pkgPath string, t types.Type, qf func(p *types.Package) string) string {
	switch t.(type) {
	case *types.Basic:
		{
			return t.String()
		}
	case *types.Named:
		{
			pkg := ast.GetNamedTypePkg(t.(*types.Named))
			if pkg == nil {
				panic(fmt.Sprintf("cannot find package name for %s", t.String()))
			}
			typeName := t.(*types.Named).Obj().Name()
			pkgName := pkg.Name()
			path := pkg.Path()
			if samePackage(path, pkgPath, pkgName) {
				return typeName
			} else {
				return types.TypeString(t, qf)
			}
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
			keyDecl := getDeclType(pkgPath, mapType.Key(), qf)
			valueDecl := getDeclType(pkgPath, mapType.Elem(), qf)
			return fmt.Sprintf("map[%s]%s", keyDecl, valueDecl)
		}
	case *types.Chan:
		{
			chanType := t.(*types.Chan)
			var chanTmlt string
			switch chanType.Dir() {
			case types.SendRecv:
				{
					chanTmlt = "chan %s"
				}
			case types.SendOnly:
				{
					chanTmlt = "chan<- %s"
				}
			default:
				{
					chanTmlt = "<-chan %s"
				}
			}
			return fmt.Sprintf(chanTmlt, getDeclType(pkgPath, chanType.Elem(), qf))
		}
	case *types.Signature:
		{
			return types.TypeString(t, qf)
		}
	default:
		panic(fmt.Sprintf("unsupported type %s", t.String()))
	}
}

func samePackage(path string, pkgPath string, pkgName string) bool {
	if pkgPath != "" {
		return path == pkgPath
	}
	return path == pkgName
}
