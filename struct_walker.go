package syrinx

import (
	"github.com/golang-collections/collections/stack"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
)

type FieldInfo struct {
	Name string
	Type types.Type
}

func (structInfo *StructInfo) addFields(names []*ast.Ident, fieldType types.Type) {
	for _, fieldName := range names {
		structInfo.addField(fieldName.Name, fieldType)
	}
}

func (structInfo *StructInfo) addField(fieldName string, fieldType types.Type) {
	structInfo.Fields = append(structInfo.Fields, &FieldInfo{
		Name: fieldName,
		Type: fieldType,
	})
}

type StructInfo struct {
	Name   string
	Fields []*FieldInfo
	Type   types.Type
}

type structWalker struct {
	structs         []*StructInfo
	structInfoStack stack.Stack
	typeNameStack   stack.Stack
	typeInfo        types.Info
}

func (walker *structWalker) Parse(sourceCode string) error {
	fileset := token.NewFileSet()
	astFile, err := parser.ParseFile(fileset, "src.go", sourceCode, 0)
	if err != nil {
		return err
	}
	typeInfo, err := walker.parseTypeInfo(fileset, astFile)
	if err != nil {
		return err
	}
	walker.typeInfo = typeInfo
	Visit(walker, astFile)
	return nil
}

func (*structWalker) parseTypeInfo(fset *token.FileSet, astFile *ast.File) (types.Info, error) {
	typeInfo := types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
	_, err := (&types.Config{Importer: importer.Default()}).Check("mypkg", fset, []*ast.File{astFile}, &typeInfo)
	return typeInfo, err
}

func (walker *structWalker) Structs() []*StructInfo {
	if walker.structs != nil {
		return walker.structs
	}
	r := make([]*StructInfo, walker.structInfoStack.Len())

	for i := len(r); i > 0; i-- {
		structInfo := walker.structInfoStack.Pop()
		r[i-1] = structInfo.(*StructInfo)
	}
	walker.structs = r
	return r
}

func (walker *structWalker) BeginWalk(node ast.Node) (w Walker) {
	if node == nil {
		return nil
	}
	return walker
}

func (*structWalker) EndWalk(node ast.Node) {
}

func (*structWalker) WalkComment(comment *ast.Comment) {
}

func (*structWalker) WalkCommentGroup(group *ast.CommentGroup) {
}

func (*structWalker) EndWalkComment(comment *ast.Comment) {
}

func (*structWalker) EndWalkCommentGroup(group *ast.CommentGroup) {
}

func (walker *structWalker) WalkField(field *ast.Field) {
	structInfo := walker.structInfoStack.Peek().(*StructInfo)
	fieldType := walker.typeInfo.Types[field.Type].Type
	emitTypeNameIfFiledIsNestedStruct(walker, fieldType)
	structInfo.addFields(field.Names, fieldType)
}

func (*structWalker) EndWalkField(field *ast.Field) {
}

func (*structWalker) WalkFieldList(list *ast.FieldList) {
}

func (*structWalker) EndWalkFieldList(list *ast.FieldList) {
}

func (*structWalker) WalkEllipsis(ellipsis *ast.Ellipsis) {
}

func (*structWalker) EndWalkEllipsis(ellipsis *ast.Ellipsis) {
}

func (*structWalker) WalkFuncLit(lit *ast.FuncLit) {
}

func (*structWalker) EndWalkFuncLit(lit *ast.FuncLit) {
}

func (*structWalker) WalkCompositeLit(lit *ast.CompositeLit) {
}

func (*structWalker) EndWalkCompositeLit(lit *ast.CompositeLit) {
}

func (*structWalker) WalkParenExpr(expr *ast.ParenExpr) {
}

func (*structWalker) EndWalkParenExpr(expr *ast.ParenExpr) {
}

func (*structWalker) WalkSelectorExpr(expr *ast.SelectorExpr) {
}

func (*structWalker) EndWalkSelectorExpr(expr *ast.SelectorExpr) {
}

func (*structWalker) WalkIndexExpr(expr *ast.IndexExpr) {
}

func (*structWalker) EndWalkIndexExpr(expr *ast.IndexExpr) {
}

func (*structWalker) WalkSliceExpr(expr *ast.SliceExpr) {
}

func (*structWalker) EndWalkSliceExpr(expr *ast.SliceExpr) {
}

func (*structWalker) WalkTypeAssertExpr(expr *ast.TypeAssertExpr) {
}

func (*structWalker) EndWalkTypeAssertExpr(expr *ast.TypeAssertExpr) {
}

func (*structWalker) WalkCallExpr(expr *ast.CallExpr) {
}

func (*structWalker) EndWalkCallExpr(expr *ast.CallExpr) {
}

func (*structWalker) WalkStarExpr(expr *ast.StarExpr) {
}

func (*structWalker) EndWalkStarExpr(expr *ast.StarExpr) {
}

func (*structWalker) WalkUnaryExpr(expr *ast.UnaryExpr) {
}

func (*structWalker) EndWalkUnaryExpr(expr *ast.UnaryExpr) {
}

func (*structWalker) WalkBinaryExpr(expr *ast.BinaryExpr) {
}

func (*structWalker) EndWalkBinaryExpr(expr *ast.BinaryExpr) {
}

func (*structWalker) WalkKeyValueExpr(expr *ast.KeyValueExpr) {
}

func (*structWalker) EndWalkKeyValueExpr(expr *ast.KeyValueExpr) {
}

func (*structWalker) WalkArrayType(arrayType *ast.ArrayType) {
}

func (*structWalker) EndWalkArrayType(arrayType *ast.ArrayType) {
}

func (walker *structWalker) WalkStructType(structTypeExpr *ast.StructType) {
	structName := walker.typeNameStack.Pop().(string)
	structType := walker.typeInfo.Types[structTypeExpr].Type
	walker.addStructInfo(structName, structType)
}

func (walker *structWalker) EndWalkStructType(structType *ast.StructType) {
}

func (*structWalker) WalkFuncType(funcType *ast.FuncType) {
}

func (*structWalker) EndWalkFuncType(funcType *ast.FuncType) {
}

func (*structWalker) WalkInterfaceType(interfaceType *ast.InterfaceType) {
}

func (*structWalker) EndWalkInterfaceType(interfaceType *ast.InterfaceType) {
}

func (*structWalker) WalkMapType(mapType *ast.MapType) {
}

func (*structWalker) EndWalkMapType(mapType *ast.MapType) {
}

func (*structWalker) WalkChanType(chanType *ast.ChanType) {
}

func (*structWalker) EndWalkChanType(chanType *ast.ChanType) {
}

func (*structWalker) WalkDeclStmt(stmt *ast.DeclStmt) {
}

func (*structWalker) EndWalkDeclStmt(stmt *ast.DeclStmt) {
}

func (*structWalker) WalkLabeledStmt(stmt *ast.LabeledStmt) {
}

func (*structWalker) EndWalkLabeledStmt(stmt *ast.LabeledStmt) {
}

func (*structWalker) WalkExprStmt(stmt *ast.ExprStmt) {
}

func (*structWalker) EndWalkExprStmt(stmt *ast.ExprStmt) {
}

func (*structWalker) WalkSendStmt(stmt *ast.SendStmt) {
}

func (*structWalker) EndWalkSendStmt(stmt *ast.SendStmt) {
}

func (*structWalker) WalkIncDecStmt(stmt *ast.IncDecStmt) {
}

func (*structWalker) EndWalkIncDecStmt(stmt *ast.IncDecStmt) {
}

func (*structWalker) WalkAssignStmt(stmt *ast.AssignStmt) {
}

func (*structWalker) EndWalkAssignStmt(stmt *ast.AssignStmt) {
}

func (*structWalker) WalkGoStmt(stmt *ast.GoStmt) {
}

func (*structWalker) EndWalkGoStmt(stmt *ast.GoStmt) {
}

func (*structWalker) WalkDeferStmt(stmt *ast.DeferStmt) {
}

func (*structWalker) EndWalkDeferStmt(stmt *ast.DeferStmt) {
}

func (*structWalker) WalkReturnStmt(stmt *ast.ReturnStmt) {
}

func (*structWalker) EndWalkReturnStmt(stmt *ast.ReturnStmt) {
}

func (*structWalker) WalkBranchStmt(stmt *ast.BranchStmt) {
}

func (*structWalker) EndWalkBranchStmt(stmt *ast.BranchStmt) {
}

func (*structWalker) WalkBlockStmt(stmt *ast.BlockStmt) {
}

func (*structWalker) EndWalkBlockStmt(stmt *ast.BlockStmt) {
}

func (*structWalker) WalkIfStmt(stmt *ast.IfStmt) {
}

func (*structWalker) EndWalkIfStmt(stmt *ast.IfStmt) {
}

func (*structWalker) WalkCaseClause(clause *ast.CaseClause) {
}

func (*structWalker) EndWalkCaseClause(clause *ast.CaseClause) {
}

func (*structWalker) WalkSwitchStmt(stmt *ast.SwitchStmt) {
}

func (*structWalker) EndWalkSwitchStmt(stmt *ast.SwitchStmt) {
}

func (*structWalker) WalkTypeSwitchStmt(stmt *ast.TypeSwitchStmt) {
}

func (*structWalker) EndWalkTypeSwitchStmt(stmt *ast.TypeSwitchStmt) {
}

func (*structWalker) WalkCommClause(clause *ast.CommClause) {
}

func (*structWalker) EndWalkCommClause(clause *ast.CommClause) {
}

func (*structWalker) WalkSelectStmt(stmt *ast.SelectStmt) {
}

func (*structWalker) EndWalkSelectStmt(stmt *ast.SelectStmt) {
}

func (*structWalker) WalkForStmt(stmt *ast.ForStmt) {
}

func (*structWalker) EndWalkForStmt(stmt *ast.ForStmt) {
}

func (*structWalker) WalkRangeStmt(stmt *ast.RangeStmt) {
}

func (*structWalker) EndWalkRangeStmt(stmt *ast.RangeStmt) {
}

func (*structWalker) WalkImportSpec(spec *ast.ImportSpec) {
}

func (*structWalker) EndWalkImportSpec(spec *ast.ImportSpec) {
}

func (*structWalker) WalkValueSpec(spec *ast.ValueSpec) {
}

func (*structWalker) EndWalkValueSpec(spec *ast.ValueSpec) {
}

func (walker *structWalker) WalkTypeSpec(spec *ast.TypeSpec) {
	walker.typeNameStack.Push(spec.Name.Name)
}

func (walker *structWalker) EndWalkTypeSpec(spec *ast.TypeSpec) {
}

func (*structWalker) WalkGenDecl(decl *ast.GenDecl) {
}

func (*structWalker) EndWalkGenDecl(decl *ast.GenDecl) {
}

func (*structWalker) WalkFuncDecl(decl *ast.FuncDecl) {
}

func (*structWalker) EndWalkFuncDecl(decl *ast.FuncDecl) {
}

func (walker *structWalker) WalkFile(file *ast.File) {
}

func (*structWalker) EndWalkFile(file *ast.File) {
}

func (*structWalker) WalkPackage(n *ast.Package) {
}

func (*structWalker) EndWalkPackage(n *ast.Package) {
}

func (*structWalker) WalkBadExpr(n *ast.BadExpr) {
}

func (*structWalker) EndWalkBadExpr(n *ast.BadExpr) {
}

func (*structWalker) WalkIdent(n *ast.Ident) {
}

func (*structWalker) EndWalkIdent(n *ast.Ident) {
}

func (*structWalker) WalkBasicLit(n *ast.BasicLit) {
}

func (*structWalker) EndWalkBasicLit(n *ast.BasicLit) {
}

func (*structWalker) WalkBadStmt(n *ast.BadStmt) {
}

func (*structWalker) EndWalkBadStmt(n *ast.BadStmt) {
}

func (*structWalker) WalkEmptyStmt(n *ast.EmptyStmt) {
}

func (*structWalker) EndWalkEmptyStmt(n *ast.EmptyStmt) {
}

func (*structWalker) WalkBadDecl(n *ast.BadDecl) {
}

func (*structWalker) EndWalkBadDecl(n *ast.BadDecl) {
}

func NewStructWalker() *structWalker {
	return &structWalker{}
}

func emitTypeNameIfFiledIsNestedStruct(walker *structWalker, fieldType types.Type) {
	switch fieldType.(type) {
	case *types.Struct:
		{
			structName := fieldType.String()
			walker.typeNameStack.Push(structName)
		}
	}
}

func (walker *structWalker) addStructInfo(structName string, structType types.Type) {
	walker.structInfoStack.Push(&StructInfo{
		Name: structName,
		Type: structType,
	})
}
