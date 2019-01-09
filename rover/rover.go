package rover

import (
	"github.com/ahmetb/go-linq"
	"github.com/lonegunmanb/syringe/ioc"
	"github.com/lonegunmanb/varys/ast"
	"reflect"
)

type codeRover struct {
	roverStartingPath string
	destinationPath   string
	packageName       string
	goPathEnv         ast.GoPathEnv
	ignorePattern     string
}

func newCodeRover(roverStartingPath string) *codeRover {
	return &codeRover{
		roverStartingPath: roverStartingPath,
		goPathEnv:         ast.NewGoPathEnv(),
	}
}

func getTypeWalker() ast.TypeWalker {
	return roverContainer.GetOrRegister((*ast.TypeWalker)(nil), func(ioc ioc.Container) interface{} {
		return ast.NewTypeWalker()
	}).(ast.TypeWalker)
}

func (r *codeRover) getStructTypes() ([]ast.TypeInfo, error) {
	typeInfos, err := r.getTypeInfos()
	if err != nil {
		return nil, err
	}
	var results []ast.TypeInfo
	linq.From(typeInfos).Where(func(t interface{}) bool {
		return t.(ast.TypeInfo).GetKind() == reflect.Struct
	}).Where(func(t interface{}) bool {
		return t.(ast.TypeInfo).GetPkgName() != "main"
	}).ToSlice(&results)
	return results, nil
}

func (r *codeRover) getTypeInfos() ([]ast.TypeInfo, error) {
	walker := getTypeWalker()
	err := walker.ParseDir(r.roverStartingPath, r.ignorePattern)
	if err != nil {
		return nil, err
	}
	return walker.GetTypes(), nil
}
