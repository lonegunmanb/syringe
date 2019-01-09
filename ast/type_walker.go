package ast

import (
	"github.com/ahmetb/go-linq"
	"github.com/golang-collections/collections/stack"
	"github.com/lonegunmanb/johnnie"
	"github.com/lonegunmanb/syringe/ioc"
	"github.com/lonegunmanb/syringe/util"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type opsKind string

var analyzingType opsKind = "isAnalyzingType"
var analyzingFunc opsKind = "analyzingFunc"

var iocContainer = ioc.NewContainer()

type TypeWalker interface {
	johnnie.Walker
	GetTypes() []TypeInfo
	Parse(pkgPath string, sourceCode string) error
	ParseDir(dirPath string, ignorePatten string) error
	SetTypeInfo(i *types.Info)
	SetPhysicalPath(p string)
	ParseAst(path string, fileAst *ast.File) error
}

type typeWalker struct {
	johnnie.DefaultWalker
	osEnv         util.GoPathEnv
	types         []*typeInfo
	typeInfoStack stack.Stack
	opsStack      stack.Stack
	typeInfo      *types.Info
	pkgPath       string
	pkgName       string
	physicalPath  string
}

func (walker *typeWalker) SetPhysicalPath(p string) {
	walker.physicalPath = p
}

func (walker *typeWalker) SetTypeInfo(i *types.Info) {
	walker.typeInfo = i
}

func (walker *typeWalker) GetTypes() []TypeInfo {
	result := make([]TypeInfo, 0, len(walker.types))
	linq.From(walker.types).Select(func(t interface{}) interface{} {
		return t.(TypeInfo)
	}).ToSlice(&result)
	return result
}

func (walker *typeWalker) Parse(pkgPath string, sourceCode string) error {
	return walker.parse(pkgPath, "src.go", sourceCode)
}

func (walker *typeWalker) ParseDir(dirPath string, ignorePatten string) error {
	fSet := token.NewFileSet()
	osEnv := getOsEnv()
	fileAstMap, err := walker.parseFileAsts(dirPath, ignorePatten, fSet, osEnv)
	if err != nil {
		return err
	}
	info, err := parseTypes(fileAstMap, fSet, osEnv)
	if err != nil {
		return err
	}
	return walker.walkAsts(fileAstMap, info)
}

func (walker *typeWalker) ParseAst(path string, fileAst *ast.File) error {
	pkgPath, err := util.GetPkgPath(walker.osEnv, path)
	if err != nil {
		return err
	}
	return walker.parseAst(pkgPath, fileAst)
}

func (walker *typeWalker) Types() []*typeInfo {
	return walker.types
}

func (walker *typeWalker) WalkFile(f *ast.File) {
	walker.pkgName = f.Name.Name
}

func (walker *typeWalker) WalkField(field *ast.Field) {
	if walker.isAnalyzingType() {
		typeInfo := walker.typeInfoStack.Peek().(*typeInfo)
		t := walker.typeInfo.Types[field.Type]
		fieldType := t.Type
		emitTypeNameIfFiledIsNestedType(walker, fieldType)
		typeInfo.processField(field, fieldType)
	}
}

func (walker *typeWalker) WalkStructType(structType *ast.StructType) {
	if walker.opsStack.Peek() == analyzingType {
		walker.addTypeInfo(structType, reflect.Struct)
	}
}

func (walker *typeWalker) EndWalkStructType(structType *ast.StructType) {
	walker.typeInfoStack.Pop()
}

func (walker *typeWalker) WalkInterfaceType(interfaceType *ast.InterfaceType) {
	if walker.opsStack.Peek() == analyzingType {
		walker.addTypeInfo(interfaceType, reflect.Interface)
	}
}

func (walker *typeWalker) EndWalkInterfaceType(interfaceType *ast.InterfaceType) {
	walker.typeInfoStack.Pop()
}

func (walker *typeWalker) WalkTypeSpec(spec *ast.TypeSpec) {
	walker.typeInfoStack.Push(spec.Name.Name)
	walker.opsStack.Push(analyzingType)
}

func (walker *typeWalker) EndWalkTypeSpec(spec *ast.TypeSpec) {
	walker.opsStack.Pop()
}

func (walker *typeWalker) WalkFuncType(funcType *ast.FuncType) {
	walker.opsStack.Push(analyzingFunc)
}

func (walker *typeWalker) EndWalkFuncType(funcType *ast.FuncType) {
	walker.opsStack.Pop()
}

func NewTypeWalker() TypeWalker {
	return newTypeWalkerWithPhysicalPath("")
}

func newTypeWalkerWithPhysicalPath(physicalPath string) TypeWalker {
	return &typeWalker{
		types:        []*typeInfo{},
		osEnv:        util.NewGoPathEnv(),
		physicalPath: physicalPath,
	}
}

func (*typeWalker) parseTypeInfo(pkgPath string, fileSet *token.FileSet,
	astFile *ast.File) (*types.Info, error) {
	typeInfo := &types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
	_, err := (&types.Config{Importer: importer.For("source", nil)}).
		Check(pkgPath, fileSet, []*ast.File{astFile}, typeInfo)
	return typeInfo, err
}

func (walker *typeWalker) addTypeInfo(structTypeExpr ast.Expr, kind reflect.Kind) {

	item := walker.typeInfoStack.Pop()
	typeName, ok := item.(string)
	if !ok {
		println(typeName)
	}
	structType := walker.typeInfo.Types[structTypeExpr].Type
	typeInfo := &typeInfo{
		Name:         typeName,
		PkgPath:      walker.pkgPath,
		PkgName:      walker.pkgName,
		PhysicalPath: walker.physicalPath,
		Type:         structType,
		Kind:         kind,
	}
	walker.typeInfoStack.Push(typeInfo)
	walker.types = append(walker.types, typeInfo)
}

func (walker *typeWalker) isAnalyzingType() bool {
	return walker.opsStack.Peek() == analyzingType
}

func emitTypeNameIfFiledIsNestedType(walker *typeWalker, fieldType types.Type) {
	switch t := fieldType.(type) {
	case *types.Struct:
		{
			walker.typeInfoStack.Push(t.String())
		}
	case *types.Interface:
		{
			walker.typeInfoStack.Push(t.String())
		}
	case *types.Pointer:
		{
			emitTypeNameIfFiledIsNestedType(walker, t.Elem())
		}
	}
}

func (walker *typeWalker) parseAst(pkgPath string, astFile *ast.File) error {
	walker.pkgPath = pkgPath
	johnnie.Visit(walker, astFile)
	return nil
}

func (walker *typeWalker) getFiles(dirPath string, ignorePatten string) ([]util.FileInfo, error) {
	fileRetrieverKey := (*util.FileRetriever)(nil)
	fileRetriever := iocContainer.GetOrRegister(fileRetrieverKey, func(ioc ioc.Container) interface{} {
		return util.NewFileRetriever()
	}).(util.FileRetriever)

	ignoreRegex, err := parseIgnorePatten(ignorePatten)
	if err != nil {
		return nil, err
	}

	filter := func(info util.FileInfo) bool {
		if ignoreRegex != nil {
			return isGoFile(info) && !ignoreRegex.MatchString(info.Name())
		}
		return isGoFile(info)
	}
	files, err := fileRetriever.GetFiles(dirPath)
	if err != nil {
		return nil, err
	}
	filteredFiles := make([]util.FileInfo, 0)
	linq.From(files).Where(func(fileInfo interface{}) bool {
		return filter(fileInfo.(util.FileInfo))
	}).ToSlice(&filteredFiles)
	return filteredFiles, nil
}

func (walker *typeWalker) parse(pkgPath string, fileName string, sourceCode string) error {
	fileset := token.NewFileSet()

	astFile, err := parser.ParseFile(fileset, fileName, sourceCode, 0)
	if err != nil {
		return err
	}

	if walker.typeInfo == nil {
		typeInfo, err := walker.parseTypeInfo(pkgPath, fileset, astFile)
		if err != nil {
			return err
		}
		walker.typeInfo = typeInfo
	}

	return walker.parseAst(pkgPath, astFile)
}

func (walker *typeWalker) walkAsts(fileMap map[string][]*ast.File, info *types.Info) error {
	for path, fileAsts := range fileMap {
		walker.SetPhysicalPath(path)
		walker.SetTypeInfo(info)
		for _, fileAst := range fileAsts {
			err := walker.ParseAst(path, fileAst)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (walker *typeWalker) parseFileAsts(dirPath string, ignorePatten string, fSet *token.FileSet,
	osEnv util.GoPathEnv) (map[string][]*ast.File, error) {
	files, err := walker.getFiles(dirPath, ignorePatten)
	if err != nil {
		return nil, err
	}
	fileMap := make(map[string][]*ast.File)
	for _, file := range files {
		fileAst, err := parser.ParseFile(fSet, osEnv.ConcatFileNameWithPath(file.Path(), file.Name()), nil, 0)
		if err != nil {
			return nil, err
		}
		fileMap[file.Path()] = append(fileMap[file.Path()], fileAst)
	}
	return fileMap, nil
}

func isStructType(t types.Type) bool {
	_, ok := t.Underlying().(*types.Struct)
	return ok
}

func isEmbeddedField(field *ast.Field) bool {
	return field.Names == nil
}

func getTag(field *ast.Field) string {
	if field.Tag == nil {
		return ""
	}
	return field.Tag.Value
}

func isGoFile(info os.FileInfo) bool {
	return !info.IsDir() && isGoSrcFile(info.Name()) && !isTestFile(info.Name())
}

func isTestFile(fileName string) bool {
	return strings.HasSuffix(strings.TrimSuffix(fileName, ".go"), "test")
}

func isGoSrcFile(fileName string) bool {
	return strings.HasSuffix(fileName, ".go")
}

func parseIgnorePatten(ignorePatten string) (*regexp.Regexp, error) {
	var regex *regexp.Regexp
	if ignorePatten != "" {
		reg, err := regexp.Compile(ignorePatten)
		if err != nil {
			return nil, err
		}
		regex = reg
	}
	return regex, nil
}

func getOsEnv() util.GoPathEnv {
	return iocContainer.GetOrRegister((*util.GoPathEnv)(nil), func(ioc ioc.Container) interface{} {
		return util.NewGoPathEnv()
	}).(util.GoPathEnv)
}

func parseTypes(fileMap map[string][]*ast.File, fSet *token.FileSet, osEnv util.GoPathEnv) (*types.Info, error) {
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}
	for path, fileAsts := range fileMap {
		var conf = &types.Config{Importer: importer.For("source", nil)}
		goPath, err := util.GetPkgPath(osEnv, path)
		if err != nil {
			return nil, err
		}
		_, err = conf.Check(goPath, fSet, fileAsts, info)
		if err != nil {
			return nil, err
		}
	}
	return info, nil
}
