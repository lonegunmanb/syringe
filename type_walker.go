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

type FieldInfo struct {
	Name string
	Type types.Type
	Tag  string
}

type TypeInfo struct {
	Name   string
	Fields []*FieldInfo
	Type   types.Type
	Kind   reflect.Kind
}

type typeWalker struct {
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

func (*typeWalker) parseTypeInfo(pkgPath string, fset *token.FileSet, astFile *ast.File) (types.Info, error) {
	typeInfo := types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
	_, err := (&types.Config{Importer: importer.Default()}).Check(pkgPath, fset, []*ast.File{astFile}, &typeInfo)
	return typeInfo, err
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

func (walker *typeWalker) BeginWalk(node ast.Node) (w Walker) {
	if node == nil {
		return nil
	}
	return walker
}

func (*typeWalker) EndWalk(node ast.Node) {
}

func (*typeWalker) WalkComment(comment *ast.Comment) {
}

func (*typeWalker) WalkCommentGroup(group *ast.CommentGroup) {
}

func (*typeWalker) EndWalkComment(comment *ast.Comment) {
}

func (*typeWalker) EndWalkCommentGroup(group *ast.CommentGroup) {
}

func (walker *typeWalker) WalkField(field *ast.Field) {
	typeInfo := walker.typeInfoStack.Peek().(*TypeInfo)
	fieldType := walker.typeInfo.Types[field.Type].Type
	emitTypeNameIfFiledIsNestedStruct(walker, fieldType)
	typeInfo.addFields(field, fieldType)
}

func (*typeWalker) EndWalkField(field *ast.Field) {
}

func (*typeWalker) WalkFieldList(list *ast.FieldList) {
}

func (*typeWalker) EndWalkFieldList(list *ast.FieldList) {
}

func (*typeWalker) WalkEllipsis(ellipsis *ast.Ellipsis) {
}

func (*typeWalker) EndWalkEllipsis(ellipsis *ast.Ellipsis) {
}

func (*typeWalker) WalkFuncLit(lit *ast.FuncLit) {
}

func (*typeWalker) EndWalkFuncLit(lit *ast.FuncLit) {
}

func (*typeWalker) WalkCompositeLit(lit *ast.CompositeLit) {
}

func (*typeWalker) EndWalkCompositeLit(lit *ast.CompositeLit) {
}

func (*typeWalker) WalkParenExpr(expr *ast.ParenExpr) {
}

func (*typeWalker) EndWalkParenExpr(expr *ast.ParenExpr) {
}

func (*typeWalker) WalkSelectorExpr(expr *ast.SelectorExpr) {
}

func (*typeWalker) EndWalkSelectorExpr(expr *ast.SelectorExpr) {
}

func (*typeWalker) WalkIndexExpr(expr *ast.IndexExpr) {
}

func (*typeWalker) EndWalkIndexExpr(expr *ast.IndexExpr) {
}

func (*typeWalker) WalkSliceExpr(expr *ast.SliceExpr) {
}

func (*typeWalker) EndWalkSliceExpr(expr *ast.SliceExpr) {
}

func (*typeWalker) WalkTypeAssertExpr(expr *ast.TypeAssertExpr) {
}

func (*typeWalker) EndWalkTypeAssertExpr(expr *ast.TypeAssertExpr) {
}

func (*typeWalker) WalkCallExpr(expr *ast.CallExpr) {
}

func (*typeWalker) EndWalkCallExpr(expr *ast.CallExpr) {
}

func (*typeWalker) WalkStarExpr(expr *ast.StarExpr) {
}

func (*typeWalker) EndWalkStarExpr(expr *ast.StarExpr) {
}

func (*typeWalker) WalkUnaryExpr(expr *ast.UnaryExpr) {
}

func (*typeWalker) EndWalkUnaryExpr(expr *ast.UnaryExpr) {
}

func (*typeWalker) WalkBinaryExpr(expr *ast.BinaryExpr) {
}

func (*typeWalker) EndWalkBinaryExpr(expr *ast.BinaryExpr) {
}

func (*typeWalker) WalkKeyValueExpr(expr *ast.KeyValueExpr) {
}

func (*typeWalker) EndWalkKeyValueExpr(expr *ast.KeyValueExpr) {
}

func (*typeWalker) WalkArrayType(arrayType *ast.ArrayType) {
}

func (*typeWalker) EndWalkArrayType(arrayType *ast.ArrayType) {
}

func (walker *typeWalker) WalkStructType(structTypeExpr *ast.StructType) {
	addTypeInfo(walker, structTypeExpr, reflect.Struct)
}

func addTypeInfo(walker *typeWalker, structTypeExpr ast.Expr, kind reflect.Kind) {
	typeName := walker.typeNameStack.Pop().(string)
	resolvedType := walker.typeInfo.Types[structTypeExpr].Type
	walker.addTypeInfo(typeName, resolvedType, kind)
}

func (walker *typeWalker) EndWalkStructType(structType *ast.StructType) {
}

func (*typeWalker) WalkFuncType(funcType *ast.FuncType) {
}

func (*typeWalker) EndWalkFuncType(funcType *ast.FuncType) {
}

func (walker *typeWalker) WalkInterfaceType(interfaceType *ast.InterfaceType) {
	addTypeInfo(walker, interfaceType, reflect.Interface)
}

func (*typeWalker) EndWalkInterfaceType(interfaceType *ast.InterfaceType) {
}

func (*typeWalker) WalkMapType(mapType *ast.MapType) {
}

func (*typeWalker) EndWalkMapType(mapType *ast.MapType) {
}

func (*typeWalker) WalkChanType(chanType *ast.ChanType) {
}

func (*typeWalker) EndWalkChanType(chanType *ast.ChanType) {
}

func (*typeWalker) WalkDeclStmt(stmt *ast.DeclStmt) {
}

func (*typeWalker) EndWalkDeclStmt(stmt *ast.DeclStmt) {
}

func (*typeWalker) WalkLabeledStmt(stmt *ast.LabeledStmt) {
}

func (*typeWalker) EndWalkLabeledStmt(stmt *ast.LabeledStmt) {
}

func (*typeWalker) WalkExprStmt(stmt *ast.ExprStmt) {
}

func (*typeWalker) EndWalkExprStmt(stmt *ast.ExprStmt) {
}

func (*typeWalker) WalkSendStmt(stmt *ast.SendStmt) {
}

func (*typeWalker) EndWalkSendStmt(stmt *ast.SendStmt) {
}

func (*typeWalker) WalkIncDecStmt(stmt *ast.IncDecStmt) {
}

func (*typeWalker) EndWalkIncDecStmt(stmt *ast.IncDecStmt) {
}

func (*typeWalker) WalkAssignStmt(stmt *ast.AssignStmt) {
}

func (*typeWalker) EndWalkAssignStmt(stmt *ast.AssignStmt) {
}

func (*typeWalker) WalkGoStmt(stmt *ast.GoStmt) {
}

func (*typeWalker) EndWalkGoStmt(stmt *ast.GoStmt) {
}

func (*typeWalker) WalkDeferStmt(stmt *ast.DeferStmt) {
}

func (*typeWalker) EndWalkDeferStmt(stmt *ast.DeferStmt) {
}

func (*typeWalker) WalkReturnStmt(stmt *ast.ReturnStmt) {
}

func (*typeWalker) EndWalkReturnStmt(stmt *ast.ReturnStmt) {
}

func (*typeWalker) WalkBranchStmt(stmt *ast.BranchStmt) {
}

func (*typeWalker) EndWalkBranchStmt(stmt *ast.BranchStmt) {
}

func (*typeWalker) WalkBlockStmt(stmt *ast.BlockStmt) {
}

func (*typeWalker) EndWalkBlockStmt(stmt *ast.BlockStmt) {
}

func (*typeWalker) WalkIfStmt(stmt *ast.IfStmt) {
}

func (*typeWalker) EndWalkIfStmt(stmt *ast.IfStmt) {
}

func (*typeWalker) WalkCaseClause(clause *ast.CaseClause) {
}

func (*typeWalker) EndWalkCaseClause(clause *ast.CaseClause) {
}

func (*typeWalker) WalkSwitchStmt(stmt *ast.SwitchStmt) {
}

func (*typeWalker) EndWalkSwitchStmt(stmt *ast.SwitchStmt) {
}

func (*typeWalker) WalkTypeSwitchStmt(stmt *ast.TypeSwitchStmt) {
}

func (*typeWalker) EndWalkTypeSwitchStmt(stmt *ast.TypeSwitchStmt) {
}

func (*typeWalker) WalkCommClause(clause *ast.CommClause) {
}

func (*typeWalker) EndWalkCommClause(clause *ast.CommClause) {
}

func (*typeWalker) WalkSelectStmt(stmt *ast.SelectStmt) {
}

func (*typeWalker) EndWalkSelectStmt(stmt *ast.SelectStmt) {
}

func (*typeWalker) WalkForStmt(stmt *ast.ForStmt) {
}

func (*typeWalker) EndWalkForStmt(stmt *ast.ForStmt) {
}

func (*typeWalker) WalkRangeStmt(stmt *ast.RangeStmt) {
}

func (*typeWalker) EndWalkRangeStmt(stmt *ast.RangeStmt) {
}

func (*typeWalker) WalkImportSpec(spec *ast.ImportSpec) {
}

func (*typeWalker) EndWalkImportSpec(spec *ast.ImportSpec) {
}

func (*typeWalker) WalkValueSpec(spec *ast.ValueSpec) {
}

func (*typeWalker) EndWalkValueSpec(spec *ast.ValueSpec) {
}

func (walker *typeWalker) WalkTypeSpec(spec *ast.TypeSpec) {
	walker.typeNameStack.Push(spec.Name.Name)
}

func (walker *typeWalker) EndWalkTypeSpec(spec *ast.TypeSpec) {
}

func (*typeWalker) WalkGenDecl(decl *ast.GenDecl) {
}

func (*typeWalker) EndWalkGenDecl(decl *ast.GenDecl) {
}

func (*typeWalker) WalkFuncDecl(decl *ast.FuncDecl) {
}

func (*typeWalker) EndWalkFuncDecl(decl *ast.FuncDecl) {
}

func (walker *typeWalker) WalkFile(file *ast.File) {
}

func (*typeWalker) EndWalkFile(file *ast.File) {
}

func (*typeWalker) WalkPackage(n *ast.Package) {
}

func (*typeWalker) EndWalkPackage(n *ast.Package) {
}

func (*typeWalker) WalkBadExpr(n *ast.BadExpr) {
}

func (*typeWalker) EndWalkBadExpr(n *ast.BadExpr) {
}

func (*typeWalker) WalkIdent(n *ast.Ident) {
}

func (*typeWalker) EndWalkIdent(n *ast.Ident) {
}

func (*typeWalker) WalkBasicLit(n *ast.BasicLit) {
}

func (*typeWalker) EndWalkBasicLit(n *ast.BasicLit) {
}

func (*typeWalker) WalkBadStmt(n *ast.BadStmt) {
}

func (*typeWalker) EndWalkBadStmt(n *ast.BadStmt) {
}

func (*typeWalker) WalkEmptyStmt(n *ast.EmptyStmt) {
}

func (*typeWalker) EndWalkEmptyStmt(n *ast.EmptyStmt) {
}

func (*typeWalker) WalkBadDecl(n *ast.BadDecl) {
}

func (*typeWalker) EndWalkBadDecl(n *ast.BadDecl) {
}

func NewTypeWalker() *typeWalker {
	return &typeWalker{}
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

func (walker *typeWalker) addTypeInfo(structName string, structType types.Type, kind reflect.Kind) {
	walker.typeInfoStack.Push(&TypeInfo{
		Name: structName,
		Type: structType,
		Kind: kind,
	})
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
