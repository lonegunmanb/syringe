package syrinx

import (
	"github.com/golang/mock/gomock"
	"go/ast"
	"go/token"
	"testing"
)

func TestWalkComment(t *testing.T) {
	comment := &ast.Comment{}
	testVisit(t, comment, func(walker *MockWalker) {
		walker.EXPECT().WalkComment(comment).Times(1)
		walker.EXPECT().EndWalkComment(comment).Times(1)
	})
}

/*
case *CommentGroup:
		v.WalkCommentGroup(n)
		for _, c := range n.List {
			Visit(v, c)
		}
		v.EndWalkCommentGroup(n)
*/
func TestWalkCommentGroup(t *testing.T) {
	commentGroup := &ast.CommentGroup{
		List: []*ast.Comment{{}},
	}
	testVisit(t, commentGroup, func(walker *MockWalker) {
		walker.ExpectWalkOnce(commentGroup.List[0])
		walker.EXPECT().WalkCommentGroup(commentGroup).Times(1)
		walker.EXPECT().EndWalkCommentGroup(commentGroup).Times(1)
	})
}

/*
case *Field:
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
*/
func TestWalkField(t *testing.T) {
	field := &ast.Field{
		Doc:     dummyCommentGroup(),
		Names:   []*ast.Ident{{}},
		Type:    dummyExpr(),
		Tag:     &ast.BasicLit{},
		Comment: dummyCommentGroup(),
	}
	testVisit(t, field, func(walker *MockWalker) {
		walker.ExpectWalkOnce(field.Doc)
		walker.ExpectWalkOnce(field.Names[0])
		walker.ExpectWalkOnce(field.Type)
		walker.ExpectWalkOnce(field.Tag)
		walker.ExpectWalkOnce(field.Comment)
		walker.EXPECT().WalkField(field).Times(1)
		walker.EXPECT().EndWalkField(field).Times(1)
	})
}

/*
case *FieldList:
	   for _, f := range n.List {
		   Visit(v, f)
	   }
*/
func TestWalkFieldList(t *testing.T) {
	fieldList := &ast.FieldList{
		List: []*ast.Field{{}},
	}
	testVisit(t, fieldList, func(walker *MockWalker) {
		walker.ExpectWalkOnce(fieldList.List[0])
		walker.EXPECT().WalkFieldList(fieldList).Times(1)
		walker.EXPECT().EndWalkFieldList(fieldList).Times(1)
	})
}

/*
case *BadExpr, *Ident, *BasicLit:
*/

func TestWalkBadExpr(t *testing.T) {
	badExpr := &ast.BadExpr{}
	testNothingToDoNode(t, badExpr)
}

func TestWalkIdent(t *testing.T) {
	ident := &ast.Ident{}
	testNothingToDoNode(t, ident)
}

func TestWalkBasicLit(t *testing.T) {
	basicLit := &ast.BasicLit{}
	testNothingToDoNode(t, basicLit)
}

/*
case *Ellipsis:
		if n.Elt != nil {
			Visit(v, n.Elt)
		}
*/
func TestWalkEllipsis(t *testing.T) {
	ellipsis := &ast.Ellipsis{
		Elt: dummyExpr(),
	}
	testVisit(t, ellipsis, func(walker *MockWalker) {
		walker.ExpectWalkOnce(ellipsis.Elt)
	})
}

/*
case *FuncLit:
	   Visit(v, n.Type)
	   Visit(v, n.Body)
*/
func TestFuncLit(t *testing.T) {
	funcLit := &ast.FuncLit{
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{},
	}
	testVisit(t, funcLit, func(walker *MockWalker) {
		walker.ExpectWalkOnce(funcLit.Type)
		walker.ExpectWalkOnce(funcLit.Body)
	})
}

/*
case *CompositeLit:
		if n.Type != nil {
			Visit(v, n.Type)
		}
		walkExprList(v, n.Elts)
*/
func TestCompositeLit(t *testing.T) {
	compositeLit := &ast.CompositeLit{
		Type: dummyExpr(),
		Elts: []ast.Expr{dummyExpr()},
	}
	testVisit(t, compositeLit, func(walker *MockWalker) {
		walker.ExpectWalkOnce(compositeLit.Type)
		walker.ExpectWalkOnce(compositeLit.Elts[0])
	})
}

/*
case *ParenExpr:
	   Visit(v, n.X)
*/
func TestWalkParenExpr(t *testing.T) {
	parenExpr := &ast.ParenExpr{
		X: dummyExpr(),
	}
	testVisit(t, parenExpr, func(walker *MockWalker) {
		walker.ExpectWalkOnce(parenExpr.X)
	})
}

/*
case *SelectorExpr:
	   Visit(v, n.X)
	   Visit(v, n.Sel)
*/

func TestWalkSelectorExpr(t *testing.T) {
	selectorExpr := &ast.SelectorExpr{
		X:   dummyExpr(),
		Sel: &ast.Ident{},
	}
	testVisit(t, selectorExpr, func(walker *MockWalker) {
		walker.ExpectWalkOnce(selectorExpr.X)
		walker.ExpectWalkOnce(selectorExpr.Sel)
	})
}

/*
case *IndexExpr:
	   Visit(v, n.X)
	   Visit(v, n.Index)
*/
func TestWalkIndexExpr(t *testing.T) {
	indexExpr := &ast.IndexExpr{
		X:     dummyExpr(),
		Index: dummyExpr(),
	}
	testVisit(t, indexExpr, func(walker *MockWalker) {
		walker.ExpectWalkOnce(indexExpr.X)
		walker.ExpectWalkOnce(indexExpr.Index)
	})
}

/*
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
*/
func TestWalkSliceExpr(t *testing.T) {
	sliceExpr := &ast.SliceExpr{
		X:    dummyExpr(),
		Low:  dummyExpr(),
		High: dummyExpr(),
		Max:  dummyExpr(),
	}
	testVisit(t, sliceExpr, func(walker *MockWalker) {
		walker.ExpectWalkOnce(sliceExpr.X)
		walker.ExpectWalkOnce(sliceExpr.Low)
		walker.ExpectWalkOnce(sliceExpr.High)
		walker.ExpectWalkOnce(sliceExpr.Max)
	})
}

/*
case *TypeAssertExpr:
		Visit(v, n.X)
		if n.Type != nil {
			Visit(v, n.Type)
		}
*/
func TestWalkTypeAssertExpr(t *testing.T) {
	typeAssertExpr := &ast.TypeAssertExpr{
		X:    dummyExpr(),
		Type: dummyExpr(),
	}
	testVisit(t, typeAssertExpr, func(walker *MockWalker) {
		walker.ExpectWalkOnce(typeAssertExpr.X)
		walker.ExpectWalkOnce(typeAssertExpr.Type)
	})
}

/*
case *CallExpr:
		Visit(v, n.Fun)
		walkExprList(v, n.Args)
*/
func TestWalkCallExpr(t *testing.T) {
	callExpr := &ast.CallExpr{
		Fun:  dummyExpr(),
		Args: []ast.Expr{dummyExpr()},
	}
	testVisit(t, callExpr, func(walker *MockWalker) {
		walker.ExpectWalkOnce(callExpr.Fun)
		walker.ExpectWalkOnce(callExpr.Args[0])
	})
}

/*
case *StarExpr:
	   Visit(v, n.X)
*/
func TestWalkStarExpr(t *testing.T) {
	starExpr := &ast.StarExpr{
		X: dummyExpr(),
	}
	testVisit(t, starExpr, func(walker *MockWalker) {
		walker.ExpectWalkOnce(starExpr.X)
	})
}

/*
case *UnaryExpr:
	  Visit(v, n.X)
*/
func TestWalkUnaryExpr(t *testing.T) {
	unaryExpr := &ast.UnaryExpr{
		X: dummyExpr(),
	}
	testVisit(t, unaryExpr, func(walker *MockWalker) {
		walker.ExpectWalkOnce(unaryExpr.X)
	})
}

/*
case *BinaryExpr:
		Visit(v, n.X)
		Visit(v, n.Y)
*/
func TestWalkBinaryExpr(t *testing.T) {
	binaryExpr := &ast.BinaryExpr{
		X: dummyExpr(),
		Y: dummyExpr(),
	}
	testVisit(t, binaryExpr, func(walker *MockWalker) {
		walker.ExpectWalkOnce(binaryExpr.X)
		walker.ExpectWalkOnce(binaryExpr.Y)
	})
}

/*
case *KeyValueExpr:
		Visit(v, n.Key)
		Visit(v, n.Value)
*/
func TestWalkKeyValueExpr(t *testing.T) {
	keyValueExpr := &ast.KeyValueExpr{
		Key:   dummyExpr(),
		Value: dummyExpr(),
	}
	testVisit(t, keyValueExpr, func(walker *MockWalker) {
		walker.ExpectWalkOnce(keyValueExpr.Key)
		walker.ExpectWalkOnce(keyValueExpr.Value)
	})
}

/*
case *ArrayType:
		if n.Len != nil {
			Visit(v, n.Len)
		}
		Visit(v, n.Elt)
*/
func TestWalkArrayType(t *testing.T) {
	arrayType := &ast.ArrayType{
		Len: dummyExpr(),
		Elt: dummyExpr(),
	}
	testVisit(t, arrayType, func(walker *MockWalker) {
		walker.ExpectWalkOnce(arrayType.Len)
		walker.ExpectWalkOnce(arrayType.Elt)
	})
}

/*
case *StructType:
		Visit(v, n.Fields)
*/
func TestWalkStuctType(t *testing.T) {
	structType := &ast.StructType{
		Fields: &ast.FieldList{},
	}
	testVisit(t, structType, func(walker *MockWalker) {
		walker.ExpectWalkOnce(structType.Fields)
	})
}

/*
case *FuncType:
		if n.Params != nil {
			Visit(v, n.Params)
		}
		if n.Results != nil {
			Visit(v, n.Results)
		}
*/
func TestWalkFuncType(t *testing.T) {
	funcType := &ast.FuncType{
		Params:  &ast.FieldList{Opening: 0},
		Results: &ast.FieldList{Opening: 1},
	}
	testVisit(t, funcType, func(walker *MockWalker) {
		walker.ExpectWalkOnce(funcType.Params)
		walker.ExpectWalkOnce(funcType.Results)
	})
}

/*
case *InterfaceType:
		Visit(v, n.Methods)
*/
func TestWalkInterfaceType(t *testing.T) {
	interfaceType := &ast.InterfaceType{
		Methods: &ast.FieldList{},
	}
	testVisit(t, interfaceType, func(walker *MockWalker) {
		walker.ExpectWalkOnce(interfaceType.Methods)
	})
}

/*
case *MapType:
		Visit(v, n.Key)
		Visit(v, n.Value)
*/
func TestWalkMapType(t *testing.T) {
	mapType := &ast.MapType{
		Key:   dummyExpr(),
		Value: dummyExpr(),
	}
	testVisit(t, mapType, func(walker *MockWalker) {
		walker.ExpectWalkOnce(mapType.Key)
		walker.ExpectWalkOnce(mapType.Value)
	})
}

/*
case *ChanType:
		Visit(v, n.Value)
*/
func TestWalkChanType(t *testing.T) {
	chanType := &ast.ChanType{
		Value: dummyExpr(),
	}
	testVisit(t, chanType, func(walker *MockWalker) {
		walker.ExpectWalkOnce(chanType.Value)
	})
}

/*
case *BadStmt:
		// nothing to do
*/
func TestWalkBadStmt(t *testing.T) {
	badStmt := &ast.BadStmt{}
	testNothingToDoNode(t, badStmt)
}

/*
case *DeclStmt:
		Visit(v, n.Decl)
*/
func TestWalkDeclStmt(t *testing.T) {
	declStmt := &ast.DeclStmt{
		Decl: dummyDecl(),
	}
	testVisit(t, declStmt, func(walker *MockWalker) {
		walker.ExpectWalkOnce(declStmt.Decl)
	})
}

/*
case *EmptyStmt:
		// nothing to do
*/
func TestWalkEmptyStmt(t *testing.T) {
	emptyStmt := &ast.EmptyStmt{}
	testNothingToDoNode(t, emptyStmt)
}

/*
case *LabeledStmt:
		Visit(v, n.Label)
		Visit(v, n.Stmt)
*/
func TestWalkLabeledStmt(t *testing.T) {
	labeledStmt := &ast.LabeledStmt{
		Label: &ast.Ident{},
		Stmt:  dummyStmt(),
	}
	testVisit(t, labeledStmt, func(walker *MockWalker) {
		walker.ExpectWalkOnce(labeledStmt.Label)
		walker.ExpectWalkOnce(labeledStmt.Stmt)
	})
}

/*
case *ExprStmt:
		Visit(v, n.X)
*/
func TestWalkExprStmt(t *testing.T) {
	exprStmt := &ast.ExprStmt{
		X: dummyExpr(),
	}
	testVisit(t, exprStmt, func(walker *MockWalker) {
		walker.ExpectWalkOnce(exprStmt.X)
	})
}

/*
case *SendStmt:
		Visit(v, n.Chan)
		Visit(v, n.Value)
*/
func TestWalkSendStmt(t *testing.T) {
	sendStmt := &ast.SendStmt{
		Chan:  dummyExpr(),
		Value: dummyExpr(),
	}
	testVisit(t, sendStmt, func(walker *MockWalker) {
		walker.ExpectWalkOnce(sendStmt.Chan)
		walker.ExpectWalkOnce(sendStmt.Value)
	})
}

/*
case *IncDecStmt:
		Visit(v, n.X)
*/
func TestWalkIncDecStmt(t *testing.T) {
	incDecStmt := &ast.IncDecStmt{
		X: dummyExpr(),
	}
	testVisit(t, incDecStmt, func(walker *MockWalker) {
		walker.ExpectWalkOnce(incDecStmt.X)
	})
}

/*
case *AssignStmt:
		walkExprList(v, n.Lhs)
		walkExprList(v, n.Rhs)
*/
func TestWalkAssignStmt(t *testing.T) {
	assignStmt := &ast.AssignStmt{
		Lhs: []ast.Expr{dummyExpr()},
		Rhs: []ast.Expr{dummyExpr()},
	}
	testVisit(t, assignStmt, func(walker *MockWalker) {
		walker.ExpectWalkOnce(assignStmt.Lhs[0])
		walker.ExpectWalkOnce(assignStmt.Rhs[0])
	})
}

/*
case *GoStmt:
		Visit(v, n.Call)
*/
func TestWalkGoStmt(t *testing.T) {
	goStmt := &ast.GoStmt{
		Call: &ast.CallExpr{},
	}
	testVisit(t, goStmt, func(walker *MockWalker) {
		walker.ExpectWalkOnce(goStmt.Call)
	})
}

/*
case *DeferStmt:
		Visit(v, n.Call)
*/
func TestWalkDeferStmt(t *testing.T) {
	deferStmt := &ast.DeferStmt{
		Call: &ast.CallExpr{},
	}
	testVisit(t, deferStmt, func(walker *MockWalker) {
		walker.ExpectWalkOnce(deferStmt.Call)
	})
}

/*
case *ReturnStmt:
		walkExprList(v, n.Results)
*/
func TestWalkReturnStmt(t *testing.T) {
	returnStmt := &ast.ReturnStmt{
		Results: []ast.Expr{dummyExpr()},
	}
	testVisit(t, returnStmt, func(walker *MockWalker) {
		walker.ExpectWalkOnce(returnStmt.Results[0])
	})
}

/*
case *BranchStmt:
		if n.Label != nil {
			Visit(v, n.Label)
		}
*/
func TestWalkBranchStmt(t *testing.T) {
	branchStmt := &ast.BranchStmt{
		Label: &ast.Ident{},
	}
	testVisit(t, branchStmt, func(walker *MockWalker) {
		walker.ExpectWalkOnce(branchStmt.Label)
	})
}

/*
case *BlockStmt:
		walkStmtList(v, n.List)
*/
func TestWalkBlockStmt(t *testing.T) {
	blockStmt := &ast.BlockStmt{
		List: []ast.Stmt{dummyStmt()},
	}
	testVisit(t, blockStmt, func(walker *MockWalker) {
		walker.ExpectWalkOnce(blockStmt.List[0])
	})
}

/*
case *IfStmt:
		if n.Init != nil {
			Visit(v, n.Init)
		}
		Visit(v, n.Cond)
		Visit(v, n.Body)
		if n.Else != nil {
			Visit(v, n.Else)
		}
*/
func TestWalkIfStmt(t *testing.T) {
	ifStmt := &ast.IfStmt{
		Init: dummyStmt(),
		Cond: dummyExpr(),
		Body: &ast.BlockStmt{},
		Else: dummyStmt(),
	}
	testVisit(t, ifStmt, func(walker *MockWalker) {
		walker.ExpectWalkOnce(ifStmt.Init)
		walker.ExpectWalkOnce(ifStmt.Cond)
		walker.ExpectWalkOnce(ifStmt.Body)
		walker.ExpectWalkOnce(ifStmt.Else)
	})
}

/*
case *CaseClause:
		walkExprList(v, n.List)
		walkStmtList(v, n.Body)
*/
func TestCaseClause(t *testing.T) {
	caseClause := &ast.CaseClause{
		List: []ast.Expr{dummyExpr()},
		Body: []ast.Stmt{dummyStmt()},
	}
	testVisit(t, caseClause, func(walker *MockWalker) {
		walker.ExpectWalkOnce(caseClause.List[0])
		walker.ExpectWalkOnce(caseClause.Body[0])
	})
}

/*
case *SwitchStmt:
		if n.Init != nil {
			Visit(v, n.Init)
		}
		if n.Tag != nil {
			Visit(v, n.Tag)
		}
		Visit(v, n.Body)
*/
func TestSwitchStmt(t *testing.T) {
	switchStmt := &ast.SwitchStmt{
		Init: dummyStmt(),
		Tag:  dummyExpr(),
		Body: &ast.BlockStmt{},
	}
	testVisit(t, switchStmt, func(walker *MockWalker) {
		walker.ExpectWalkOnce(switchStmt.Init)
		walker.ExpectWalkOnce(switchStmt.Tag)
		walker.ExpectWalkOnce(switchStmt.Body)
	})
}

/*
case *CommClause:
		if n.Comm != nil {
			Visit(v, n.Comm)
		}
		walkStmtList(v, n.Body)
*/
func TestWalkCommClause(t *testing.T) {
	commClause := &ast.CommClause{
		Comm: dummyStmt(),
		Body: []ast.Stmt{dummyStmt()},
	}
	testVisit(t, commClause, func(walker *MockWalker) {
		walker.ExpectWalkOnce(commClause.Comm)
		walker.ExpectWalkOnce(commClause.Body[0])
	})
}

/*
case *SelectStmt:
		Visit(v, n.Body)
*/
func TestWalkSelectStmt(t *testing.T) {
	selectStmt := &ast.SelectStmt{
		Body: &ast.BlockStmt{},
	}
	testVisit(t, selectStmt, func(walker *MockWalker) {
		walker.ExpectWalkOnce(selectStmt.Body)
	})
}

/*
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
*/
func TestWalkForStmt(t *testing.T) {
	forStmt := &ast.ForStmt{
		Init: dummyStmt(),
		Cond: dummyExpr(),
		Post: dummyStmt(),
		Body: &ast.BlockStmt{},
	}
	testVisit(t, forStmt, func(walker *MockWalker) {
		walker.ExpectWalkOnce(forStmt.Init)
		walker.ExpectWalkOnce(forStmt.Cond)
		walker.ExpectWalkOnce(forStmt.Post)
		walker.ExpectWalkOnce(forStmt.Body)
	})
}

/*
case *RangeStmt:
		if n.Key != nil {
			Visit(v, n.Key)
		}
		if n.Value != nil {
			Visit(v, n.Value)
		}
		Visit(v, n.X)
		Visit(v, n.Body)
*/
func TestWalkRangeStmt(t *testing.T) {
	rangeStmt := &ast.RangeStmt{
		Key:   dummyExpr(),
		Value: dummyExpr(),
		X:     dummyExpr(),
		Body:  &ast.BlockStmt{},
	}
	testVisit(t, rangeStmt, func(walker *MockWalker) {
		walker.ExpectWalkOnce(rangeStmt.Key)
		walker.ExpectWalkOnce(rangeStmt.Value)
		walker.ExpectWalkOnce(rangeStmt.X)
		walker.ExpectWalkOnce(rangeStmt.Body)
	})
}

/*
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

*/
func TestWalkImportSpec(t *testing.T) {
	importSpec := &ast.ImportSpec{
		Doc:     dummyCommentGroup(),
		Name:    &ast.Ident{},
		Path:    &ast.BasicLit{},
		Comment: dummyCommentGroup(),
	}
	testVisit(t, importSpec, func(walker *MockWalker) {
		walker.ExpectWalkOnce(importSpec.Doc)
		walker.ExpectWalkOnce(importSpec.Name)
		walker.ExpectWalkOnce(importSpec.Path)
		walker.ExpectWalkOnce(importSpec.Comment)
	})
}

/*
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
*/
func TestWalkValueSpec(t *testing.T) {
	valueSpec := &ast.ValueSpec{
		Doc:     dummyCommentGroup(),
		Names:   []*ast.Ident{{}},
		Type:    dummyExpr(),
		Values:  []ast.Expr{dummyExpr()},
		Comment: dummyCommentGroup(),
	}
	testVisit(t, valueSpec, func(walker *MockWalker) {
		walker.ExpectWalkOnce(valueSpec.Doc)
		walker.ExpectWalkOnce(valueSpec.Names[0])
		walker.ExpectWalkOnce(valueSpec.Type)
		walker.ExpectWalkOnce(valueSpec.Values[0])
		walker.ExpectWalkOnce(valueSpec.Comment)
	})
}

/*
case *TypeSpec:
		if n.Doc != nil {
			Visit(v, n.Doc)
		}
		Visit(v, n.Name)
		Visit(v, n.Type)
		if n.Comment != nil {
			Visit(v, n.Comment)
		}
*/
func TestWalkTypeSpec(t *testing.T) {
	typeSpec := &ast.TypeSpec{
		Doc:     dummyCommentGroup(),
		Name:    &ast.Ident{},
		Type:    dummyExpr(),
		Comment: dummyCommentGroup(),
	}
	testVisit(t, typeSpec, func(walker *MockWalker) {
		walker.ExpectWalkOnce(typeSpec.Doc)
		walker.ExpectWalkOnce(typeSpec.Name)
		walker.ExpectWalkOnce(typeSpec.Type)
		walker.ExpectWalkOnce(typeSpec.Comment)
	})
}

/*
case *BadDecl:
*/
func TestWalkBadDecl(t *testing.T) {
	badDecl := &ast.BadDecl{}
	testNothingToDoNode(t, badDecl)
}

/*
case *GenDecl:
		if n.Doc != nil {
			Visit(v, n.Doc)
		}
		for _, s := range n.Specs {
			Visit(v, s)
		}
*/
func TestWalkGenDecl(t *testing.T) {
	genDecl := &ast.GenDecl{
		Doc:   dummyCommentGroup(),
		Specs: []ast.Spec{&ast.TypeSpec{}},
	}
	testVisit(t, genDecl, func(walker *MockWalker) {
		walker.ExpectWalkOnce(genDecl.Doc)
		walker.ExpectWalkOnce(genDecl.Specs[0])
	})
}

/*
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
*/
func TestWalkFuncDecl(t *testing.T) {
	funcDecl := &ast.FuncDecl{
		Doc:  dummyCommentGroup(),
		Recv: &ast.FieldList{},
		Name: &ast.Ident{},
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{},
	}
	testVisit(t, funcDecl, func(walker *MockWalker) {
		walker.ExpectWalkOnce(funcDecl.Doc)
		walker.ExpectWalkOnce(funcDecl.Recv)
		walker.ExpectWalkOnce(funcDecl.Name)
		walker.ExpectWalkOnce(funcDecl.Type)
		walker.ExpectWalkOnce(funcDecl.Body)
	})
}

/*
case *File:
		if n.Doc != nil {
			Visit(v, n.Doc)
		}
		Visit(v, n.Name)
		walkDeclList(v, n.Decls)
*/
func TestWalkFile(t *testing.T) {
	file := &ast.File{
		Doc:   dummyCommentGroup(),
		Name:  &ast.Ident{},
		Decls: []ast.Decl{dummyDecl()},
	}
	testVisit(t, file, func(walker *MockWalker) {
		walker.ExpectWalkOnce(file.Doc)
		walker.ExpectWalkOnce(file.Name)
		walker.ExpectWalkOnce(file.Decls[0])
	})
}

/*
case *Package:
		for _, f := range n.Files {
			Visit(v, f)
		}
*/
func TestWalkPackage(t *testing.T) {
	file := &ast.File{}
	packageDecl := &ast.Package{
		Files: map[string]*ast.File{"": file},
	}
	testVisit(t, packageDecl, func(walker *MockWalker) {
		walker.ExpectWalkOnce(file)
	})
}

func (v *MockWalker) ExpectWalkOnce(i interface{}) {
	v.EXPECT().BeginWalk(gomock.Eq(i)).Times(1).Return(nil)
}

func testVisit(t *testing.T, sut ast.Node, setup func(*MockWalker)) {
	ctrl, walker := createMock(t, sut)
	defer ctrl.Finish()
	setupWalkAndEndWalk(sut, ctrl, walker)
	setup(walker)
	Visit(walker, sut)
}

func setupWalkAndEndWalk(sut ast.Node, ctrl *gomock.Controller, walker *MockWalker) {
	walkMethodName := getWalkNodeMethodName(sut)
	endWalkMethodName := getEndWalkNodeMethodName(sut)
	ctrl.RecordCall(walker.recorder.mock, walkMethodName, sut)
	ctrl.RecordCall(walker.recorder.mock, endWalkMethodName, sut)
}

func testNothingToDoNode(t *testing.T, node ast.Node) {
	ctrl, walker := createMock(t, node)
	defer ctrl.Finish()
	setupWalkAndEndWalk(node, ctrl, walker)
	Visit(walker, node)
}

func createMock(t *testing.T, testNode interface{}) (*gomock.Controller, *MockWalker) {
	ctrl := gomock.NewController(t)
	walker := NewMockWalker(ctrl)
	walker.EXPECT().BeginWalk(gomock.Nil()).Times(1)
	walker.EXPECT().BeginWalk(gomock.Eq(testNode)).Times(1).Return(walker)
	return ctrl, walker
}

type DummyExpr struct {
	ast.BadExpr
	Index int
}

type DummyDecl struct {
	ast.BadDecl
	Index int
}

type DummyStmt struct {
	ast.BadStmt
	Index int
}

var dummyCounter int = 0

func dummyExpr() *DummyExpr {
	dummyCounter++
	return &DummyExpr{Index: dummyCounter}
}

func dummyDecl() *DummyDecl {
	dummyCounter++
	return &DummyDecl{Index: dummyCounter}
}

func dummyStmt() *DummyStmt {
	dummyCounter++
	return &DummyStmt{Index: dummyCounter}
}

func dummyCommentGroup() *ast.CommentGroup {
	dummyCounter++
	return &ast.CommentGroup{List: []*ast.Comment{{Slash: token.Pos(dummyCounter)}}}
}
