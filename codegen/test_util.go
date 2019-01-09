package codegen

import (
	"github.com/lonegunmanb/varys/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func parseCode(t *testing.T, sourceCode string) ast.TypeWalker {
	typeWalker := ast.NewTypeWalker()
	err := typeWalker.Parse("github.com/lonegunmanb/syringe/test", sourceCode)
	assert.Nil(t, err)
	return typeWalker
}
