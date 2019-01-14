package codegen

import (
	"fmt"
	"github.com/lonegunmanb/varys/ast"
	"go/types"
	"regexp"
	"strings"
)

type productFieldInfoWrap struct {
	ast.FieldInfo
	typeInfo TypeInfoWrap
}

//  By Interface
//	product.Decoration = container.Resolve("github.com/lonegunmanb/syringe/test_code/fly_car.Decoration").(Decoration)
//  By Struct
//  product.Car = *container.Resolve("github.com/lonegunmanb/syringe/test_code/fly_car.Decoration").(*car.Car)
//  By Pointer
//  product.Car = container.Resolve("github.com/lonegunmanb/syringe/test_code/fly_car.Decoration").(*car.Car)
const fieldAssignTemplate = `%s.%s = %s%s.Resolve("%s").(%s)`

func (f *productFieldInfoWrap) AssembleCode() string {
	fieldType := f.GetType()
	key := fieldType.String()
	if f.GetTag() != "" {
		tagKey := getInjectKeyFromTag(f.GetTag())
		if tagKey != "" {
			key = tagKey
		}
	}
	pkgPath := f.GetReferenceFromType().GetPkgPath()
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
			if isInterface(t) {
				star = ""
			} else {
				declType = fmt.Sprintf("*%s", declType)
			}
		}
	case *types.Pointer:
		{
			star = ""
			key = t.Elem().String()
		}
	default:
		panic(fmt.Sprintf("unknown type %s", t.String()))
	}

	return fmt.Sprintf(fieldAssignTemplate, ProductIdentName, f.GetName(), star, ContainerIdentName, key, declType)
}

func isInterface(t *types.Named) bool {
	switch t.Underlying().(type) {
	case *types.Interface:
		return true
	default:
		return false
	}
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
