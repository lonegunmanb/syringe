package rover

import (
	"github.com/ahmetb/go-linq"
	"github.com/lonegunmanb/syringe/ast"
	"os"
	"reflect"
	"strings"
)

type codeRover struct {
	roverStartingPath string
	destinationPath   string
	packageName       string
	goPathEnv         ast.GoPathEnv
	fileRetriever     FileRetriever
	walkerFactory     func() ast.TypeWalker
}

func newCodeRover(roverStartingPath string) *codeRover {
	return &codeRover{
		roverStartingPath: roverStartingPath,
		goPathEnv:         ast.NewGoPathEnv(),
		fileRetriever:     &fileRetriever{},
		walkerFactory: func() ast.TypeWalker {
			return ast.NewTypeWalker()
		},
	}
}

func (r *codeRover) getStructTypes() ([]ast.TypeInfo, error) {
	types, err := r.getTypeInfos()
	if err != nil {
		return nil, err
	}
	var results []ast.TypeInfo
	linq.From(types).Where(func(t interface{}) bool {
		return t.(ast.TypeInfo).GetKind() == reflect.Struct
	}).Where(func(t interface{}) bool {
		return t.(ast.TypeInfo).GetPkgName() != "main"
	}).ToSlice(&results)
	return results, nil
}

func isGoFile(info os.FileInfo) bool {
	return !info.IsDir() && isGoSrcFile(info.Name()) && !isTestFile(info.Name())
}

func (r *codeRover) getTypeInfos() ([]ast.TypeInfo, error) {
	files, err := r.fileRetriever.GetFiles(r.roverStartingPath, isGoFile)
	if err != nil {
		return nil, err
	}
	types := make([]ast.TypeInfo, 0, len(files))
	for _, file := range files {
		walker := r.walkerFactory()
		err := walker.ParseFile(file.Path(), file.Name())
		if err != nil {
			return nil, err
		}
		for _, typeInfo := range walker.GetTypes() {
			types = append(types, typeInfo)
		}
	}
	return types, nil
}

func isTestFile(fileName string) bool {
	return strings.HasSuffix(strings.TrimSuffix(fileName, ".go"), "test")
}

func isGoSrcFile(fileName string) bool {
	return strings.HasSuffix(fileName, ".go")
}
