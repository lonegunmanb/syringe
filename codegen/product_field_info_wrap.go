package codegen

import (
	"fmt"
	"github.com/lonegunmanb/syrinx/ast"
	"go/types"
)

type productFieldInfoWrap struct {
	ast.FieldInfo
}

var qf = func(p *types.Package) string {
	return p.Name()
}

//	product.Decoration = container.Resolve("github.com/lonegunmanb/syrinx/test_code/fly_car.Decoration").(Decoration)

const field_assign_template = `product.%s = container.Resolve("%s").(%s)`

func (f *productFieldInfoWrap) AssembleCode() string {
	fieldType := f.GetType()
	key := fieldType.String()
	if f.GetTag() != "" {
		key = f.GetTag()
	}
	pkgPath := f.GetReferenceFrom().GetPkgPath()
	return fmt.Sprintf(field_assign_template, f.GetName(), key, getDeclType(pkgPath, fieldType))
}
func getDeclType(pkgPath string, t types.Type) string {
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

			path := pkg.Path()
			if path != pkgPath {
				return types.TypeString(t, qf)
			} else {
				return typeName
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
			return fmt.Sprintf("*%s", getDeclType(pkgPath, t.(*types.Pointer).Elem()))
		}
	case *types.Slice:
		{
			return fmt.Sprintf("[]%s", getDeclType(pkgPath, t.(*types.Slice).Elem()))
		}
	case *types.Array:
		{
			return fmt.Sprintf("[%d]%s", t.(*types.Array).Len(), getDeclType(pkgPath, t.(*types.Array).Elem()))
		}
	case *types.Map:
		{
			mapType := t.(*types.Map)
			keyDecl := getDeclType(pkgPath, mapType.Key())
			valueDecl := getDeclType(pkgPath, mapType.Elem())
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
			return fmt.Sprintf(chanTmlt, getDeclType(pkgPath, chanType.Elem()))
		}
	case *types.Signature:
		{
			return types.TypeString(t, qf)
		}
	default:
		panic(fmt.Sprintf("unsupported type %s", t.String()))
	}
}
