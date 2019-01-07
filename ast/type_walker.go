package ast

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/golang-collections/collections/stack"
	"github.com/lonegunmanb/johnnie"
	"github.com/lonegunmanb/syringe/util"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"reflect"
)

type opsKind string

var analyzingType opsKind = "isAnalyzingType"
var analyzingFunc opsKind = "analyzingFunc"

type TypeWalker interface {
	johnnie.Walker
	GetTypes() []TypeInfo
	Parse(pkgPath string, sourceCode string) error
	ParseFile(path string, fileName string) error
	SetTypeInfo(i *types.Info)
	SetPhysicalPath(p string)
}

type typeWalker struct {
	johnnie.DefaultWalker
	osEnv         GoPathEnv
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

func (walker *typeWalker) ParseFile(path string, fileName string) error {
	walker.physicalPath = path
	osEnv := walker.osEnv
	filePath := osEnv.ConcatFileNameWithPath(path, fileName)
	return util.CallSingleRet(func() (interface{}, error) {
		return ioutil.ReadFile(filePath)
	}).CallBiRet(func(buffer interface{}) (interface{}, interface{}, error) {
		sourceCode := string(buffer.([]byte))
		pkgPath, err := GetPkgPath(osEnv, path)
		return sourceCode, pkgPath, err
	}).Call(func(sourceCode interface{}, pkgPath interface{}) error {
		return walker.parse(pkgPath.(string), fileName, sourceCode.(string))
	}).Err
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

	walker.pkgPath = pkgPath
	johnnie.Visit(walker, astFile)
	return nil
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
		osEnv:        &envImpl{},
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

func getTag(field *ast.Field) string {
	if field.Tag == nil {
		return ""
	}
	return field.Tag.Value
}

func (typeInfo *typeInfo) processField(field *ast.Field, fieldType types.Type) {
	if isEmbeddedField(field) {
		typeInfo.addInheritance(field, fieldType)
	} else {
		typeInfo.addFieldInfos(field, fieldType)
	}
}

func (typeInfo *typeInfo) addFieldInfos(field *ast.Field, fieldType types.Type) {
	names := field.Names
	for _, fieldName := range names {
		typeInfo.Fields = append(typeInfo.Fields, &fieldInfo{
			Name:          fieldName.Name,
			Type:          fieldType,
			Tag:           getTag(field),
			ReferenceFrom: typeInfo,
		})
	}
}

func (typeInfo *typeInfo) addInheritance(field *ast.Field, fieldType types.Type) {
	var kind EmbeddedKind
	var packagePath string
	switch t := fieldType.(type) {
	case *types.Named:
		{
			if isStructType(t) {
				kind = EmbeddedByStruct
			} else {
				kind = EmbeddedByInterface
			}
			packagePath = GetNamedTypePkg(t).Path()
		}
	case *types.Pointer:
		{
			elemType, ok := t.Elem().(*types.Named)
			if !ok {
				panic(fmt.Sprintf("unknown embedded type %s", fieldType.String()))
			}
			kind = EmbeddedByPointer
			packagePath = GetNamedTypePkg(elemType).Path()
		}
	default:
		panic(fmt.Sprintf("unknown embedded type %s", t.String()))
	}

	embeddedType := &embeddedType{
		Kind:          kind,
		FullName:      fieldType.String(),
		PkgPath:       packagePath,
		Tag:           getTag(field),
		Type:          fieldType,
		ReferenceFrom: typeInfo,
	}
	typeInfo.EmbeddedTypes = append(typeInfo.EmbeddedTypes, embeddedType)
}

func isStructType(t types.Type) bool {
	_, ok := t.Underlying().(*types.Struct)
	return ok
}

func isEmbeddedField(field *ast.Field) bool {
	return field.Names == nil
}
