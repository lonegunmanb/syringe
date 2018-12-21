package syrinx

import (
	"github.com/golang-collections/collections/stack"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"reflect"
)

type typeWalker struct {
	DefaultWalker
	types         []*TypeInfo
	typeInfoStack stack.Stack
	typeNameStack stack.Stack
	typeInfo      types.Info
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
	Visit(walker, astFile)
	return nil
}

func (walker *typeWalker) Types() []*TypeInfo {
	if walker.types != nil {
		return walker.types
	}
	r := make([]*TypeInfo, walker.typeInfoStack.Len())

	for i := len(r); i > 0; i-- {
		structInfo := walker.typeInfoStack.Pop()
		r[i-1] = structInfo.(*TypeInfo)
	}
	walker.types = r
	return r
}

func (walker *typeWalker) WalkField(field *ast.Field) {
	typeInfo := walker.typeInfoStack.Peek().(*TypeInfo)
	fieldType := walker.typeInfo.Types[field.Type].Type
	emitTypeNameIfFiledIsNestedStruct(walker, fieldType)
	typeInfo.addFields(field, fieldType)
}

func (walker *typeWalker) WalkStructType(structTypeExpr *ast.StructType) {
	addTypeInfo(walker, structTypeExpr, reflect.Struct)
}

func (walker *typeWalker) WalkInterfaceType(interfaceType *ast.InterfaceType) {
	addTypeInfo(walker, interfaceType, reflect.Interface)
}

func (walker *typeWalker) WalkTypeSpec(spec *ast.TypeSpec) {
	walker.typeNameStack.Push(spec.Name.Name)
}

func NewTypeWalker() *typeWalker {
	return &typeWalker{}
}

func (*typeWalker) parseTypeInfo(pkgPath string, fset *token.FileSet, astFile *ast.File) (types.Info, error) {
	typeInfo := types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
	_, err := (&types.Config{Importer: importer.Default()}).Check(pkgPath, fset, []*ast.File{astFile}, &typeInfo)
	return typeInfo, err
}

func addTypeInfo(walker *typeWalker, structTypeExpr ast.Expr, kind reflect.Kind) {
	typeName := walker.typeNameStack.Pop().(string)
	resolvedType := walker.typeInfo.Types[structTypeExpr].Type
	walker.addTypeInfo(typeName, resolvedType, kind)
}

func (walker *typeWalker) addTypeInfo(structName string, structType types.Type, kind reflect.Kind) {
	walker.typeInfoStack.Push(&TypeInfo{
		Name: structName,
		Type: structType,
		Kind: kind,
	})
}

func emitTypeNameIfFiledIsNestedStruct(walker *typeWalker, fieldType types.Type) {
	switch fieldType.(type) {
	case *types.Struct:
		{
			typeName := fieldType.String()
			walker.typeNameStack.Push(typeName)
		}
	}
}

func getTag(field *ast.Field) string {
	if field.Tag == nil {
		return ""
	}
	return field.Tag.Value
}

func (structInfo *TypeInfo) addFields(field *ast.Field, fieldType types.Type) {
	names := field.Names
	for _, fieldName := range names {
		structInfo.Fields = append(structInfo.Fields, &FieldInfo{
			Name: fieldName.Name,
			Type: fieldType,
			Tag:  getTag(field),
		})
	}
}
