package ast

import (
	"fmt"
	"github.com/golang-collections/collections/stack"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"reflect"
)

type opsKind string

var analyzingType opsKind = "isAnalyzingType"
var analyzingFunc opsKind = "analyzingFunc"

type typeWalker struct {
	DefaultWalker
	types         []*TypeInfo
	typeInfoStack stack.Stack
	opsStack      stack.Stack
	typeInfo      types.Info
	pkgPath       string
}

func (walker *typeWalker) Parse(pkgPath string, sourceCode string) error {
	fileset := token.NewFileSet()
	astFile, err := parser.ParseFile(fileset, "src.go", sourceCode, 0)
	if err != nil {
		return err
	}
	typeInfo, err := walker.parseTypeInfo(pkgPath, fileset, astFile)
	if err != nil {
		return err
	}
	walker.typeInfo = typeInfo
	walker.pkgPath = pkgPath
	Visit(walker, astFile)
	return nil
}

func (walker *typeWalker) Types() []*TypeInfo {
	return walker.types
}

func (walker *typeWalker) WalkField(field *ast.Field) {
	if walker.isAnalyzingType() {
		typeInfo := walker.typeInfoStack.Peek().(*TypeInfo)
		fieldType := walker.typeInfo.Types[field.Type].Type
		emitTypeNameIfFiledIsNestedType(walker, fieldType)
		typeInfo.processField(field, fieldType)
	}
}

func (walker *typeWalker) WalkStructType(structType *ast.StructType) {
	walker.addTypeInfo(structType, reflect.Struct)
}

func (walker *typeWalker) EndWalkStructType(structType *ast.StructType) {
	walker.typeInfoStack.Pop()
}

func (walker *typeWalker) WalkInterfaceType(interfaceType *ast.InterfaceType) {
	walker.addTypeInfo(interfaceType, reflect.Interface)
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

func NewTypeWalker() *typeWalker {
	return &typeWalker{
		types: []*TypeInfo{},
	}
}

func (*typeWalker) parseTypeInfo(pkgPath string, fileSet *token.FileSet,
	astFile *ast.File) (types.Info, error) {
	typeInfo := types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
	_, err := (&types.Config{Importer: importer.Default()}).
		Check(pkgPath, fileSet, []*ast.File{astFile}, &typeInfo)
	return typeInfo, err
}

func (walker *typeWalker) addTypeInfo(structTypeExpr ast.Expr, kind reflect.Kind) {
	typeName := walker.typeInfoStack.Pop().(string)
	structType := walker.typeInfo.Types[structTypeExpr].Type
	typeInfo := &TypeInfo{
		Name:    typeName,
		PkgPath: walker.pkgPath,
		Type:    structType,
		Kind:    kind,
	}
	walker.typeInfoStack.Push(typeInfo)
	walker.types = append(walker.types, typeInfo)
}

func (walker *typeWalker) isAnalyzingType() bool {
	return walker.opsStack.Peek() == analyzingType
}

func emitTypeNameIfFiledIsNestedType(walker *typeWalker, fieldType types.Type) {
	switch fieldType.(type) {
	case *types.Struct:
		{
			typeName := fieldType.String()
			walker.typeInfoStack.Push(typeName)
		}
	case *types.Interface:
		{
			typeName := fieldType.String()
			walker.typeInfoStack.Push(typeName)
		}
	}
}

func getTag(field *ast.Field) string {
	if field.Tag == nil {
		return ""
	}
	return field.Tag.Value
}

func (typeInfo *TypeInfo) processField(field *ast.Field, fieldType types.Type) {
	if isEmbeddedField(field) {
		typeInfo.addInheritance(field, fieldType)
	} else {
		typeInfo.addFieldInfos(field, fieldType)
	}
}

func (typeInfo *TypeInfo) addFieldInfos(field *ast.Field, fieldType types.Type) {
	names := field.Names
	for _, fieldName := range names {
		typeInfo.Fields = append(typeInfo.Fields, &FieldInfo{
			Name:          fieldName.Name,
			Type:          fieldType,
			Tag:           getTag(field),
			ReferenceFrom: typeInfo,
		})
	}
}

func (typeInfo *TypeInfo) addInheritance(field *ast.Field, fieldType types.Type) {
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
			packagePath = getNamedTypePkg(t)
		}
	case *types.Pointer:
		{
			elemType, ok := t.Elem().(*types.Named)
			if !ok {
				panic(fmt.Sprintf("unknown embedded type %s", fieldType.String()))
			}
			kind = EmbeddedByPointer
			packagePath = getNamedTypePkg(elemType)
		}
	default:
		panic(fmt.Sprintf("unknown embedded type %s", t.String()))
	}

	embeddedType := &EmbeddedType{
		Kind:     kind,
		FullName: fieldType.String(),
		PkgPath:  packagePath,
		Tag:      getTag(field),
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
