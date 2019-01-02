package codegen

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenRegisterPackageDecl(t *testing.T) {
	writer := &bytes.Buffer{}
	sut := &registerCodegen{pkgName: "ast", writer: writer}
	err := sut.genPkgDecl()
	assert.Nil(t, err)
	assert.Equal(t, expectedPackageDecl, writer.String())
}

func TestGenRegisterImportDecl(t *testing.T) {
	writer := &bytes.Buffer{}
	ctrl := gomock.NewController(t)
	mockDepPkgPathInfo := NewMockDepPkgPathInfo(ctrl)
	mockDepPkgPathInfo.EXPECT().GenImportDecls().Times(1).Return([]string{
		`"go/ast"`,
		`"go/token"`,
		`"go/types"`,
	})
	sut := &registerCodegen{writer: writer, depPkgPathInfo: mockDepPkgPathInfo}
	err := sut.genImportDecls()
	assert.Nil(t, err)
	assert.Equal(t, expectedImportDecl, writer.String())
}
