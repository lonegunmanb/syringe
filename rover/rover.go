package rover

import (
	"github.com/ahmetb/go-linq"
	"github.com/lonegunmanb/johnnie"
	"github.com/lonegunmanb/syringe/ast"
	goast "go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
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

func isGoFile(info os.FileInfo) bool {
	return !info.IsDir() && isGoSrcFile(info.Name()) && !isTestFile(info.Name())
}

func (r *codeRover) getTypeInfos() ([]ast.TypeInfo, error) {
	files, err := r.fileRetriever.GetFiles(r.roverStartingPath, isGoFile)
	if err != nil {
		return nil, err
	}
	fileMap := make(map[string][]*goast.File)
	fset := token.NewFileSet()
	osEnv := ast.NewGoPathEnv()
	for _, file := range files {
		fileAst, err := parser.ParseFile(fset, osEnv.ConcatFileNameWithPath(file.Path(), file.Name()), nil, 0)
		if err != nil {
			return nil, err
		}
		fileMap[file.Path()] = append(fileMap[file.Path()], fileAst)
	}
	info := &types.Info{
		Types: make(map[goast.Expr]types.TypeAndValue),
		Defs:  make(map[*goast.Ident]types.Object),
		Uses:  make(map[*goast.Ident]types.Object),
	}
	for path, fileAsts := range fileMap {
		var conf = &types.Config{Importer: importer.For("source", nil)}
		goPath, err := ast.GetPkgPath(osEnv, path)
		if err != nil {
			return nil, err
		}
		_, err = conf.Check(getPkgName(goPath), fset, fileAsts, info)
		if err != nil {
			return nil, err
		}
	}
	typeInfos := make([]ast.TypeInfo, 0, len(files))
	for path, fileAsts := range fileMap {
		for _, fileAst := range fileAsts {
			walker := r.walkerFactory()
			walker.SetPhysicalPath(path)
			walker.SetTypeInfo(info)
			//err := walker.ParseFile(file.Path(), file.Name())
			johnnie.Visit(walker, fileAst)
			if err != nil {
				return nil, err
			}
			for _, typeInfo := range walker.GetTypes() {
				typeInfo.GetPhysicalPath()
				typeInfos = append(typeInfos, typeInfo)
			}
		}
	}
	return typeInfos, nil
}

func isTestFile(fileName string) bool {
	return strings.HasSuffix(strings.TrimSuffix(fileName, ".go"), "test")
}

func isGoSrcFile(fileName string) bool {
	return strings.HasSuffix(fileName, ".go")
}
