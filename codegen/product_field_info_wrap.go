package codegen

import (
	"fmt"
	"github.com/lonegunmanb/syrinx/ast"
	"go/types"
)

type productFieldInfoWrap struct {
	ast.FieldInfo
}

//  By Interface
//	product.Decoration = container.Resolve("github.com/lonegunmanb/syrinx/test_code/fly_car.Decoration").(Decoration)
//  By Struct
//  product.Car = *container.Resolve("github.com/lonegunmanb/syrinx/test_code/fly_car.Decoration").(*car.Car)
//  By Pointer
//  product.Car = container.Resolve("github.com/lonegunmanb/syrinx/test_code/fly_car.Decoration").(*car.Car)
const fieldAssignTemplate = `product.%s = %scontainer.Resolve("%s").(%s)`

func (f *productFieldInfoWrap) AssembleCode() string {
	fieldType := f.GetType()
	key := fieldType.String()
	if f.GetTag() != "" {
		key = f.GetTag()
	}
	pkgPath := f.GetReferenceFrom().GetPkgPath()
	declType := getDeclType(pkgPath, fieldType)
	star := "*"
	switch t := f.GetType().(type) {
	case *types.Named:
		{
			switch t.Underlying().(type) {
			case *types.Struct:
				{
					declType = fmt.Sprintf("*%s", declType)
				}
			case *types.Interface:
				{
					star = ""
				}
			}
		}
	case *types.Pointer:
		{
			star = ""
			key = t.Elem().String()
		}
	}

	return fmt.Sprintf(fieldAssignTemplate, f.GetName(), star, key, declType)
}
