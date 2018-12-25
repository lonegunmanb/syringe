package codegen

import (
	"bytes"
	"github.com/lonegunmanb/syrinx/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenPackageDecl(t *testing.T) {
	writer := &bytes.Buffer{}
	typeInfo := &ast.TypeInfo{
		PkgName: "ast",
	}
	codegen := newCodegen(typeInfo, writer)
	err := codegen.genPkgDecl()
	assert.Nil(t, err)
	code := writer.String()
	assert.Equal(t, "package ast", code)
}
