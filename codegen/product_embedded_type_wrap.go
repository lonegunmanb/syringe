package codegen

import (
	"fmt"
	"github.com/lonegunmanb/syrinx/ast"
	"go/types"
)

type productEmbeddedTypeWrap struct {
	ast.EmbeddedType
}

const embeddedTypeInitTemplate = `product.%s = %scontainer.Resolve("%s").(%s)`

func (f *productEmbeddedTypeWrap) AssembleCode() string {
	pkgPath := f.GetReferenceFrom().GetPkgPath()
	typeDecl := getDeclType(pkgPath, f.GetType())

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
	return fmt.Sprintf(embeddedTypeInitTemplate, name, star, embeddedType.String(), typeDecl)
}
