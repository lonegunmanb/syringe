package ast

import "go/ast"

type DefaultWalker struct {
}

func (walker *DefaultWalker) BeginWalk(node ast.Node) (w Walker) {
	if node == nil {
		w = nil
	} else {
		w = walker
	}
	return
}

func (*DefaultWalker) EndWalk(node ast.Node) {

}

func (*DefaultWalker) WalkComment(comment *ast.Comment) {

}

func (*DefaultWalker) WalkCommentGroup(group *ast.CommentGroup) {

}

func (*DefaultWalker) EndWalkComment(comment *ast.Comment) {

}

func (*DefaultWalker) EndWalkCommentGroup(group *ast.CommentGroup) {

}

func (*DefaultWalker) WalkField(field *ast.Field) {

}

func (*DefaultWalker) EndWalkField(field *ast.Field) {

}

func (*DefaultWalker) WalkFieldList(list *ast.FieldList) {

}

func (*DefaultWalker) EndWalkFieldList(list *ast.FieldList) {

}

func (*DefaultWalker) WalkEllipsis(ellipsis *ast.Ellipsis) {

}

func (*DefaultWalker) EndWalkEllipsis(ellipsis *ast.Ellipsis) {

}

func (*DefaultWalker) WalkFuncLit(lit *ast.FuncLit) {

}

func (*DefaultWalker) EndWalkFuncLit(lit *ast.FuncLit) {

}

func (*DefaultWalker) WalkCompositeLit(lit *ast.CompositeLit) {

}

func (*DefaultWalker) EndWalkCompositeLit(lit *ast.CompositeLit) {

}

func (*DefaultWalker) WalkParenExpr(expr *ast.ParenExpr) {

}

func (*DefaultWalker) EndWalkParenExpr(expr *ast.ParenExpr) {

}

func (*DefaultWalker) WalkSelectorExpr(expr *ast.SelectorExpr) {

}

func (*DefaultWalker) EndWalkSelectorExpr(expr *ast.SelectorExpr) {

}

func (*DefaultWalker) WalkIndexExpr(expr *ast.IndexExpr) {

}

func (*DefaultWalker) EndWalkIndexExpr(expr *ast.IndexExpr) {

}

func (*DefaultWalker) WalkSliceExpr(expr *ast.SliceExpr) {

}

func (*DefaultWalker) EndWalkSliceExpr(expr *ast.SliceExpr) {

}

func (*DefaultWalker) WalkTypeAssertExpr(expr *ast.TypeAssertExpr) {

}

func (*DefaultWalker) EndWalkTypeAssertExpr(expr *ast.TypeAssertExpr) {

}

func (*DefaultWalker) WalkCallExpr(expr *ast.CallExpr) {

}

func (*DefaultWalker) EndWalkCallExpr(expr *ast.CallExpr) {

}

func (*DefaultWalker) WalkStarExpr(expr *ast.StarExpr) {

}

func (*DefaultWalker) EndWalkStarExpr(expr *ast.StarExpr) {

}

func (*DefaultWalker) WalkUnaryExpr(expr *ast.UnaryExpr) {

}

func (*DefaultWalker) EndWalkUnaryExpr(expr *ast.UnaryExpr) {

}

func (*DefaultWalker) WalkBinaryExpr(expr *ast.BinaryExpr) {

}

func (*DefaultWalker) EndWalkBinaryExpr(expr *ast.BinaryExpr) {

}

func (*DefaultWalker) WalkKeyValueExpr(expr *ast.KeyValueExpr) {

}

func (*DefaultWalker) EndWalkKeyValueExpr(expr *ast.KeyValueExpr) {

}

func (*DefaultWalker) WalkArrayType(arrayType *ast.ArrayType) {

}

func (*DefaultWalker) EndWalkArrayType(arrayType *ast.ArrayType) {

}

func (*DefaultWalker) WalkStructType(structType *ast.StructType) {

}

func (*DefaultWalker) EndWalkStructType(structType *ast.StructType) {

}

func (*DefaultWalker) WalkFuncType(funcType *ast.FuncType) {

}

func (*DefaultWalker) EndWalkFuncType(funcType *ast.FuncType) {

}

func (*DefaultWalker) WalkInterfaceType(interfaceType *ast.InterfaceType) {

}

func (*DefaultWalker) EndWalkInterfaceType(interfaceType *ast.InterfaceType) {

}

func (*DefaultWalker) WalkMapType(mapType *ast.MapType) {

}

func (*DefaultWalker) EndWalkMapType(mapType *ast.MapType) {

}

func (*DefaultWalker) WalkChanType(chanType *ast.ChanType) {

}

func (*DefaultWalker) EndWalkChanType(chanType *ast.ChanType) {

}

func (*DefaultWalker) WalkDeclStmt(stmt *ast.DeclStmt) {

}

func (*DefaultWalker) EndWalkDeclStmt(stmt *ast.DeclStmt) {

}

func (*DefaultWalker) WalkLabeledStmt(stmt *ast.LabeledStmt) {

}

func (*DefaultWalker) EndWalkLabeledStmt(stmt *ast.LabeledStmt) {

}

func (*DefaultWalker) WalkExprStmt(stmt *ast.ExprStmt) {

}

func (*DefaultWalker) EndWalkExprStmt(stmt *ast.ExprStmt) {

}

func (*DefaultWalker) WalkSendStmt(stmt *ast.SendStmt) {

}

func (*DefaultWalker) EndWalkSendStmt(stmt *ast.SendStmt) {

}

func (*DefaultWalker) WalkIncDecStmt(stmt *ast.IncDecStmt) {

}

func (*DefaultWalker) EndWalkIncDecStmt(stmt *ast.IncDecStmt) {

}

func (*DefaultWalker) WalkAssignStmt(stmt *ast.AssignStmt) {

}

func (*DefaultWalker) EndWalkAssignStmt(stmt *ast.AssignStmt) {

}

func (*DefaultWalker) WalkGoStmt(stmt *ast.GoStmt) {

}

func (*DefaultWalker) EndWalkGoStmt(stmt *ast.GoStmt) {

}

func (*DefaultWalker) WalkDeferStmt(stmt *ast.DeferStmt) {

}

func (*DefaultWalker) EndWalkDeferStmt(stmt *ast.DeferStmt) {

}

func (*DefaultWalker) WalkReturnStmt(stmt *ast.ReturnStmt) {

}

func (*DefaultWalker) EndWalkReturnStmt(stmt *ast.ReturnStmt) {

}

func (*DefaultWalker) WalkBranchStmt(stmt *ast.BranchStmt) {

}

func (*DefaultWalker) EndWalkBranchStmt(stmt *ast.BranchStmt) {

}

func (*DefaultWalker) WalkBlockStmt(stmt *ast.BlockStmt) {

}

func (*DefaultWalker) EndWalkBlockStmt(stmt *ast.BlockStmt) {

}

func (*DefaultWalker) WalkIfStmt(stmt *ast.IfStmt) {

}

func (*DefaultWalker) EndWalkIfStmt(stmt *ast.IfStmt) {

}

func (*DefaultWalker) WalkCaseClause(clause *ast.CaseClause) {

}

func (*DefaultWalker) EndWalkCaseClause(clause *ast.CaseClause) {

}

func (*DefaultWalker) WalkSwitchStmt(stmt *ast.SwitchStmt) {

}

func (*DefaultWalker) EndWalkSwitchStmt(stmt *ast.SwitchStmt) {

}

func (*DefaultWalker) WalkTypeSwitchStmt(stmt *ast.TypeSwitchStmt) {

}

func (*DefaultWalker) EndWalkTypeSwitchStmt(stmt *ast.TypeSwitchStmt) {

}

func (*DefaultWalker) WalkCommClause(clause *ast.CommClause) {

}

func (*DefaultWalker) EndWalkCommClause(clause *ast.CommClause) {

}

func (*DefaultWalker) WalkSelectStmt(stmt *ast.SelectStmt) {

}

func (*DefaultWalker) EndWalkSelectStmt(stmt *ast.SelectStmt) {

}

func (*DefaultWalker) WalkForStmt(stmt *ast.ForStmt) {

}

func (*DefaultWalker) EndWalkForStmt(stmt *ast.ForStmt) {

}

func (*DefaultWalker) WalkRangeStmt(stmt *ast.RangeStmt) {

}

func (*DefaultWalker) EndWalkRangeStmt(stmt *ast.RangeStmt) {

}

func (*DefaultWalker) WalkImportSpec(spec *ast.ImportSpec) {

}

func (*DefaultWalker) EndWalkImportSpec(spec *ast.ImportSpec) {

}

func (*DefaultWalker) WalkValueSpec(spec *ast.ValueSpec) {

}

func (*DefaultWalker) EndWalkValueSpec(spec *ast.ValueSpec) {

}

func (*DefaultWalker) WalkTypeSpec(spec *ast.TypeSpec) {

}

func (*DefaultWalker) EndWalkTypeSpec(spec *ast.TypeSpec) {

}

func (*DefaultWalker) WalkGenDecl(decl *ast.GenDecl) {

}

func (*DefaultWalker) EndWalkGenDecl(decl *ast.GenDecl) {

}

func (*DefaultWalker) WalkFuncDecl(decl *ast.FuncDecl) {

}

func (*DefaultWalker) EndWalkFuncDecl(decl *ast.FuncDecl) {

}

func (*DefaultWalker) WalkFile(file *ast.File) {

}

func (*DefaultWalker) EndWalkFile(file *ast.File) {

}

func (*DefaultWalker) WalkPackage(n *ast.Package) {

}

func (*DefaultWalker) EndWalkPackage(n *ast.Package) {

}

func (*DefaultWalker) WalkBadExpr(n *ast.BadExpr) {

}

func (*DefaultWalker) EndWalkBadExpr(n *ast.BadExpr) {

}

func (*DefaultWalker) WalkIdent(n *ast.Ident) {

}

func (*DefaultWalker) EndWalkIdent(n *ast.Ident) {

}

func (*DefaultWalker) WalkBasicLit(n *ast.BasicLit) {

}

func (*DefaultWalker) EndWalkBasicLit(n *ast.BasicLit) {

}

func (*DefaultWalker) WalkBadStmt(n *ast.BadStmt) {

}

func (*DefaultWalker) EndWalkBadStmt(n *ast.BadStmt) {

}

func (*DefaultWalker) WalkEmptyStmt(n *ast.EmptyStmt) {

}

func (*DefaultWalker) EndWalkEmptyStmt(n *ast.EmptyStmt) {

}

func (*DefaultWalker) WalkBadDecl(n *ast.BadDecl) {

}

func (*DefaultWalker) EndWalkBadDecl(n *ast.BadDecl) {

}
