package codegen

//go:generate mockgen -package=codegen -destination=./mock_type_info.go github.com/lonegunmanb/syrinx/ast TypeInfo
import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenPackageDecl(t *testing.T) {
	writer := &bytes.Buffer{}
	ctrl, typeInfo := prepareMock(t)
	defer ctrl.Finish()
	typeInfo.EXPECT().GetPkgName().Times(1).Return("ast")
	codegen := newCodegen(typeInfo, writer)
	err := codegen.genPkgDecl()
	assert.Nil(t, err)
	code := writer.String()
	assert.Equal(t, "package ast", code)
}

func TestGenImportsDecl(t *testing.T) {
	writer := &bytes.Buffer{}
	ctrl, typeInfo := prepareMock(t)
	defer ctrl.Finish()
	depImports := []string{
		"go/ast",
		"go/token",
		"go/types",
	}
	setupMockToGenImports(typeInfo, depImports)
	codegen := newCodegen(typeInfo, writer)
	err := codegen.genImportDecls()
	assert.Nil(t, err)
	code := writer.String()
	expected := `
import (
"github.com/lonegunmanb/syrinx/ioc"
"go/ast"
"go/token"
"go/types"
)`
	assert.Equal(t, expected, code)
}

func TestShouldNotGenExtraImportsIfDepPathsEmpty(t *testing.T) {
	writer := &bytes.Buffer{}
	ctrl, typeInfo := prepareMock(t)
	defer ctrl.Finish()
	var depImports []string
	setupMockToGenImports(typeInfo, depImports)
	codegen := newCodegen(typeInfo, writer)
	err := codegen.genImportDecls()
	assert.Nil(t, err)
	code := writer.String()
	expected := `
import (
"github.com/lonegunmanb/syrinx/ioc"
)`
	assert.Equal(t, expected, code)
}

func setupMockToGenImports(typeInfo *MockTypeInfo, depImports []string) {
	typeInfo.EXPECT().GetDepPkgPaths().Times(1).Return(depImports)
}

func prepareMock(t *testing.T) (*gomock.Controller, *MockTypeInfo) {
	ctrl := gomock.NewController(t)
	mock := NewMockTypeInfo(ctrl)
	return ctrl, mock
}
