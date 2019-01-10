package codegen

import (
	"github.com/golang/mock/gomock"
	"github.com/lonegunmanb/varys/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func parseCode(t *testing.T, sourceCode string) ast.TypeWalker {
	pkgPath := "github.com/lonegunmanb/syringe/test"
	typeWalker := createTypeWalker(t, pkgPath)
	defer ast.ClearTypeRegister()
	err := typeWalker.Parse(pkgPath, sourceCode)
	assert.Nil(t, err)
	return typeWalker
}

func createTypeWalker(t *testing.T, pkgPath string) ast.TypeWalker {
	ctrl := gomock.NewController(t)
	mockOsEnv := NewMockGoPathEnv(ctrl)
	mockOsEnv.EXPECT().GetPkgPath(gomock.Any()).Times(1).Return(pkgPath, nil)
	ast.RegisterType((*ast.GoPathEnv)(nil), func() interface{} {
		return mockOsEnv
	})
	walker := ast.NewTypeWalker()
	return walker
}
