// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syrinx

import (
	"fmt"
	. "go/ast"
	"reflect"
)

// A Walker's BeginWalk method is invoked for each node encountered by Visit.
// If the result visitor w is not nil, Visit visits each of the children
// of node with the visitor w, followed by a call of w.BeginWalk(nil).
type Walker interface {
	BeginWalk(node Node) (w Walker)
	EndWalk(node Node)
	WalkComment(comment *Comment)
	WalkCommentGroup(group *CommentGroup)
	EndWalkComment(comment *Comment)
	EndWalkCommentGroup(group *CommentGroup)
	WalkField(field *Field)
	EndWalkField(field *Field)
	WalkFieldList(list *FieldList)
	EndWalkFieldList(list *FieldList)
	WalkEllipsis(ellipsis *Ellipsis)
	EndWalkEllipsis(ellipsis *Ellipsis)
	WalkFuncLit(lit *FuncLit)
	EndWalkFuncLit(lit *FuncLit)
	WalkCompositeLit(lit *CompositeLit)
	EndWalkCompositeLit(lit *CompositeLit)
	WalkParenExpr(expr *ParenExpr)
	EndWalkParenExpr(expr *ParenExpr)
	WalkSelectorExpr(expr *SelectorExpr)
	EndWalkSelectorExpr(expr *SelectorExpr)
	WalkIndexExpr(expr *IndexExpr)
	EndWalkIndexExpr(expr *IndexExpr)
	WalkSliceExpr(expr *SliceExpr)
	EndWalkSliceExpr(expr *SliceExpr)
	WalkTypeAssertExpr(expr *TypeAssertExpr)
	EndWalkTypeAssertExpr(expr *TypeAssertExpr)
	WalkCallExpr(expr *CallExpr)
	EndWalkCallExpr(expr *CallExpr)
	WalkStarExpr(expr *StarExpr)
	EndWalkStarExpr(expr *StarExpr)
	WalkUnaryExpr(expr *UnaryExpr)
	EndWalkUnaryExpr(expr *UnaryExpr)
	WalkBinaryExpr(expr *BinaryExpr)
	EndWalkBinaryExpr(expr *BinaryExpr)
	WalkKeyValueExpr(expr *KeyValueExpr)
	EndWalkKeyValueExpr(expr *KeyValueExpr)
	WalkArrayType(arrayType *ArrayType)
	EndWalkArrayType(arrayType *ArrayType)
	WalkStructType(structType *StructType)
	EndWalkStructType(structType *StructType)
	WalkFuncType(funcType *FuncType)
	EndWalkFuncType(funcType *FuncType)
	WalkInterfaceType(interfaceType *InterfaceType)
	EndWalkInterfaceType(interfaceType *InterfaceType)
	WalkMapType(mapType *MapType)
	EndWalkMapType(mapType *MapType)
	WalkChanType(chanType *ChanType)
	EndWalkChanType(chanType *ChanType)
	WalkDeclStmt(stmt *DeclStmt)
	EndWalkDeclStmt(stmt *DeclStmt)
	WalkLabeledStmt(stmt *LabeledStmt)
	EndWalkLabeledStmt(stmt *LabeledStmt)
	WalkExprStmt(stmt *ExprStmt)
	EndWalkExprStmt(stmt *ExprStmt)
	WalkSendStmt(stmt *SendStmt)
	EndWalkSendStmt(stmt *SendStmt)
	WalkIncDecStmt(stmt *IncDecStmt)
	EndWalkIncDecStmt(stmt *IncDecStmt)
	WalkAssignStmt(stmt *AssignStmt)
	EndWalkAssignStmt(stmt *AssignStmt)
	WalkGoStmt(stmt *GoStmt)
	EndWalkGoStmt(stmt *GoStmt)
	WalkDeferStmt(stmt *DeferStmt)
	EndWalkDeferStmt(stmt *DeferStmt)
	WalkReturnStmt(stmt *ReturnStmt)
	EndWalkReturnStmt(stmt *ReturnStmt)
	WalkBranchStmt(stmt *BranchStmt)
	EndWalkBranchStmt(stmt *BranchStmt)
	WalkBlockStmt(stmt *BlockStmt)
	EndWalkBlockStmt(stmt *BlockStmt)
	WalkIfStmt(stmt *IfStmt)
	EndWalkIfStmt(stmt *IfStmt)
	WalkCaseClause(clause *CaseClause)
	EndWalkCaseClause(clause *CaseClause)
	WalkSwitchStmt(stmt *SwitchStmt)
	EndWalkSwitchStmt(stmt *SwitchStmt)
	WalkTypeSwitchStmt(stmt *TypeSwitchStmt)
	EndWalkTypeSwitchStmt(stmt *TypeSwitchStmt)
	WalkCommClause(clause *CommClause)
	EndWalkCommClause(clause *CommClause)
	WalkSelectStmt(stmt *SelectStmt)
	EndWalkSelectStmt(stmt *SelectStmt)
	WalkForStmt(stmt *ForStmt)
	EndWalkForStmt(stmt *ForStmt)
	WalkRangeStmt(stmt *RangeStmt)
	EndWalkRangeStmt(stmt *RangeStmt)
	WalkImportSpec(spec *ImportSpec)
	EndWalkImportSpec(spec *ImportSpec)
	WalkValueSpec(spec *ValueSpec)
	EndWalkValueSpec(spec *ValueSpec)
	WalkTypeSpec(spec *TypeSpec)
	EndWalkTypeSpec(spec *TypeSpec)
	WalkGenDecl(decl *GenDecl)
	EndWalkGenDecl(decl *GenDecl)
	WalkFuncDecl(decl *FuncDecl)
	EndWalkFuncDecl(decl *FuncDecl)
	WalkFile(file *File)
	EndWalkFile(file *File)
	WalkPackage(n *Package)
	EndWalkPackage(n *Package)
}

// Helper functions for common node lists. They may be empty.

func walkIdentList(v Walker, list []*Ident) {
	for _, x := range list {
		Visit(v, x)
	}
}

func walkExprList(v Walker, list []Expr) {
	for _, x := range list {
		Visit(v, x)
	}
}

func walkStmtList(v Walker, list []Stmt) {
	for _, x := range list {
		Visit(v, x)
	}
}

func walkDeclList(v Walker, list []Decl) {
	for _, x := range list {
		Visit(v, x)
	}
}

// TODO(gri): Investigate if providing a closure to Visit leads to
//            simpler use (and may help eliminate Inspect in turn).

// Visit traverses an AST in depth-first order: It starts by calling
// v.BeginWalk(node); node must not be nil. If the visitor w returned by
// v.BeginWalk(node) is not nil, Visit is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.BeginWalk(nil).
//
func Visit(v Walker, node Node) {
	if v = v.BeginWalk(node); v == nil {
		return
	}
	walkNode(v, node)

	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {
	// Comments and fields
	case *Comment:
		v.WalkComment(n)
		v.EndWalkComment(n)
		// nothing to do

	case *CommentGroup:
		v.WalkCommentGroup(n)
		for _, c := range n.List {
			Visit(v, c)
		}
		v.EndWalkCommentGroup(n)

	case *Field:
		v.WalkField(n)
		if n.Doc != nil {
			Visit(v, n.Doc)
		}
		walkIdentList(v, n.Names)
		Visit(v, n.Type)
		if n.Tag != nil {
			Visit(v, n.Tag)
		}
		if n.Comment != nil {
			Visit(v, n.Comment)
		}
		v.EndWalkField(n)

	case *FieldList:
		v.WalkFieldList(n)
		for _, f := range n.List {
			Visit(v, f)
		}
		v.EndWalkFieldList(n)

	// Expressions
	case *BadExpr, *Ident, *BasicLit:
		// nothing to do

	case *Ellipsis:
		if n.Elt != nil {
			Visit(v, n.Elt)
		}

	case *FuncLit:
		Visit(v, n.Type)
		Visit(v, n.Body)

	case *CompositeLit:
		if n.Type != nil {
			Visit(v, n.Type)
		}
		walkExprList(v, n.Elts)

	case *ParenExpr:
		Visit(v, n.X)

	case *SelectorExpr:
		Visit(v, n.X)
		Visit(v, n.Sel)

	case *IndexExpr:
		Visit(v, n.X)
		Visit(v, n.Index)

	case *SliceExpr:
		Visit(v, n.X)
		if n.Low != nil {
			Visit(v, n.Low)
		}
		if n.High != nil {
			Visit(v, n.High)
		}
		if n.Max != nil {
			Visit(v, n.Max)
		}

	case *TypeAssertExpr:
		Visit(v, n.X)
		if n.Type != nil {
			Visit(v, n.Type)
		}

	case *CallExpr:
		Visit(v, n.Fun)
		walkExprList(v, n.Args)

	case *StarExpr:
		Visit(v, n.X)

	case *UnaryExpr:
		Visit(v, n.X)

	case *BinaryExpr:
		Visit(v, n.X)
		Visit(v, n.Y)

	case *KeyValueExpr:
		Visit(v, n.Key)
		Visit(v, n.Value)

	// Types
	case *ArrayType:
		if n.Len != nil {
			Visit(v, n.Len)
		}
		Visit(v, n.Elt)

	case *StructType:
		Visit(v, n.Fields)

	case *FuncType:
		if n.Params != nil {
			Visit(v, n.Params)
		}
		if n.Results != nil {
			Visit(v, n.Results)
		}

	case *InterfaceType:
		Visit(v, n.Methods)

	case *MapType:
		Visit(v, n.Key)
		Visit(v, n.Value)

	case *ChanType:
		Visit(v, n.Value)

	// Statements
	case *BadStmt:
		// nothing to do

	case *DeclStmt:
		Visit(v, n.Decl)

	case *EmptyStmt:
		// nothing to do

	case *LabeledStmt:
		Visit(v, n.Label)
		Visit(v, n.Stmt)

	case *ExprStmt:
		Visit(v, n.X)

	case *SendStmt:
		Visit(v, n.Chan)
		Visit(v, n.Value)

	case *IncDecStmt:
		Visit(v, n.X)

	case *AssignStmt:
		walkExprList(v, n.Lhs)
		walkExprList(v, n.Rhs)

	case *GoStmt:
		Visit(v, n.Call)

	case *DeferStmt:
		Visit(v, n.Call)

	case *ReturnStmt:
		walkExprList(v, n.Results)

	case *BranchStmt:
		if n.Label != nil {
			Visit(v, n.Label)
		}

	case *BlockStmt:
		walkStmtList(v, n.List)

	case *IfStmt:
		if n.Init != nil {
			Visit(v, n.Init)
		}
		Visit(v, n.Cond)
		Visit(v, n.Body)
		if n.Else != nil {
			Visit(v, n.Else)
		}

	case *CaseClause:
		walkExprList(v, n.List)
		walkStmtList(v, n.Body)

	case *SwitchStmt:
		if n.Init != nil {
			Visit(v, n.Init)
		}
		if n.Tag != nil {
			Visit(v, n.Tag)
		}
		Visit(v, n.Body)

	case *TypeSwitchStmt:
		if n.Init != nil {
			Visit(v, n.Init)
		}
		Visit(v, n.Assign)
		Visit(v, n.Body)

	case *CommClause:
		if n.Comm != nil {
			Visit(v, n.Comm)
		}
		walkStmtList(v, n.Body)

	case *SelectStmt:
		Visit(v, n.Body)

	case *ForStmt:
		if n.Init != nil {
			Visit(v, n.Init)
		}
		if n.Cond != nil {
			Visit(v, n.Cond)
		}
		if n.Post != nil {
			Visit(v, n.Post)
		}
		Visit(v, n.Body)

	case *RangeStmt:
		if n.Key != nil {
			Visit(v, n.Key)
		}
		if n.Value != nil {
			Visit(v, n.Value)
		}
		Visit(v, n.X)
		Visit(v, n.Body)

	// Declarations
	case *ImportSpec:
		if n.Doc != nil {
			Visit(v, n.Doc)
		}
		if n.Name != nil {
			Visit(v, n.Name)
		}
		Visit(v, n.Path)
		if n.Comment != nil {
			Visit(v, n.Comment)
		}

	case *ValueSpec:
		if n.Doc != nil {
			Visit(v, n.Doc)
		}
		walkIdentList(v, n.Names)
		if n.Type != nil {
			Visit(v, n.Type)
		}
		walkExprList(v, n.Values)
		if n.Comment != nil {
			Visit(v, n.Comment)
		}

	case *TypeSpec:
		if n.Doc != nil {
			Visit(v, n.Doc)
		}
		Visit(v, n.Name)
		Visit(v, n.Type)
		if n.Comment != nil {
			Visit(v, n.Comment)
		}

	case *BadDecl:
		// nothing to do

	case *GenDecl:
		if n.Doc != nil {
			Visit(v, n.Doc)
		}
		for _, s := range n.Specs {
			Visit(v, s)
		}

	case *FuncDecl:
		if n.Doc != nil {
			Visit(v, n.Doc)
		}
		if n.Recv != nil {
			Visit(v, n.Recv)
		}
		Visit(v, n.Name)
		Visit(v, n.Type)
		if n.Body != nil {
			Visit(v, n.Body)
		}

	// Files and packages
	case *File:
		if n.Doc != nil {
			Visit(v, n.Doc)
		}
		Visit(v, n.Name)
		walkDeclList(v, n.Decls)
		// don't walk n.Comments - they have been
		// visited already through the individual
		// nodes

	case *Package:
		for _, f := range n.Files {
			Visit(v, f)
		}

	default:
		panic(fmt.Sprintf("ast.Visit: unexpected node type %T", n))
	}

	endWalkNode(v, node)

	v.BeginWalk(nil)
}

func walkNode(v Walker, node Node) {
	callMethod(node, v, getWalkNodeMethodName(node))
}

func endWalkNode(v Walker, node Node) {
	callMethod(node, v, getEndWalkNodeMethodName(node))
}

func getWalkNodeMethodName(node Node) string {
	nodeName := getNodeName(node)
	return fmt.Sprintf("Walk%s", nodeName)
}

func getEndWalkNodeMethodName(node Node) string {
	nodeName := getNodeName(node)
	return fmt.Sprintf("EndWalk%s", nodeName)
}

func callMethod(node Node, v Walker, walkMethodName string) {
	walkMethod := getMethod(v, walkMethodName)
	walkMethod.Call([]reflect.Value{reflect.ValueOf(node)})
}

func getMethod(v Walker, walkMethodName string) reflect.Value {
	value := reflect.ValueOf(v)
	walkMethod := value.MethodByName(walkMethodName)
	return walkMethod
}

func getNodeName(node Node) string {
	nodeType := reflect.TypeOf(node).Elem()
	nodeTypeName := nodeType.Name()
	return nodeTypeName
}
