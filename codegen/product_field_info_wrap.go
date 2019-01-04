package codegen

import (
	"fmt"
	"github.com/lonegunmanb/syrinx/ast"
	"go/types"
	"regexp"
	"strings"
)

type productFieldInfoWrap struct {
	ast.FieldInfo
	typeInfo TypeInfoWrap
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
		tagKey := getInjectKeyFromTag(f.GetTag())
		if tagKey != "" {
			key = tagKey
		}
	}
	pkgPath := f.GetReferenceFrom().GetPkgPath()
	declType := getDeclType(pkgPath, fieldType, func(p *types.Package) string {
		return f.typeInfo.GetPkgNameFromPkgPath(p.Path())
	})
	star := "*"
	switch t := f.GetType().(type) {
	case *types.Basic:
		{
			star = ""
		}
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

var injectTagRegex = regexp.MustCompile("inject:\".*\"")

func hasInjectTag(tag string) bool {
	return injectTagRegex.MatchString(tag)
}

func getInjectKeyFromTag(tag string) string {
	injectTag := injectTagRegex.FindString(tag)
	if injectTag == "" {
		return injectTag
	}
	return strings.TrimPrefix(strings.TrimSuffix(injectTag, "\""), "inject:\"")
}
