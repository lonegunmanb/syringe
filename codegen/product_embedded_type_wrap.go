package codegen

import (
	"fmt"
	"github.com/lonegunmanb/varys/ast"
	"go/types"
)

type productEmbeddedTypeWrap struct {
	ast.EmbeddedType
	typeInfoWrap TypeInfoWrap
}

const embeddedTypeInitTemplate = `%s.%s = %s%s.Resolve("%s").(%s)`

func (f *productEmbeddedTypeWrap) AssembleCode() string {
	pkgPath := f.GetReferenceFrom().GetPkgPath()
	typeDecl := getDeclType(pkgPath, f.GetType(), func(p *types.Package) string {
		return f.typeInfoWrap.GetPkgNameFromPkgPath(p.Path())
	})

	var name string
	embeddedType := f.GetType()
	switch t := embeddedType.(type) {
	case *types.Pointer:
		{
			embeddedType = t.Elem()
		}
	}
	if f.GetKind() == ast.EmbeddedByPointer {
		name = f.GetType().(*types.Pointer).Elem().(*types.Named).Obj().Name()
	} else {
		name = f.GetType().(*types.Named).Obj().Name()
	}
	star := ""
	if f.GetKind() == ast.EmbeddedByStruct {
		star = "*"
		typeDecl = fmt.Sprintf("*%s", typeDecl)
	}
	return fmt.Sprintf(embeddedTypeInitTemplate, ProductIdentName, name, star, ContainerIdentName, embeddedType.String(), typeDecl)
}
