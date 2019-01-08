package rover

import (
	"github.com/ahmetb/go-linq"
	"github.com/lonegunmanb/syringe/ast"
	"reflect"
)

type codeRover struct {
	roverStartingPath string
	destinationPath   string
	packageName       string
	goPathEnv         ast.GoPathEnv
	ignorePatten      string
	walkerFactory     func() ast.TypeWalker
}

func newCodeRover(roverStartingPath string) *codeRover {
	return &codeRover{
		roverStartingPath: roverStartingPath,
		goPathEnv:         ast.NewGoPathEnv(),
		walkerFactory: func() ast.TypeWalker {
			return ast.NewTypeWalker()
		},
	}
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
	walker := r.walkerFactory()
	err := walker.ParseDir(r.roverStartingPath, r.ignorePatten)
	if err != nil {
		return nil, err
	}
	return walker.GetTypes(), nil
}
