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
	typeInfo.EXPECT().DepPkgPaths().Times(1).Return(depImports)
	codegen := newCodegen(typeInfo, writer)
	err := codegen.genImportDecls()
	assert.Nil(t, err)
	code := writer.String()
	expected := `
import (
"go/ast"
"go/token"
"go/types"
)`
	assert.Equal(t, expected, code)
}

func prepareMock(t *testing.T) (*gomock.Controller, *MockTypeInfo) {
	ctrl := gomock.NewController(t)
	mock := NewMockTypeInfo(ctrl)
	return ctrl, mock
}
