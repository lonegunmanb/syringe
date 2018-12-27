package ast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const pkgPath = "github.com/lonegunmanb/syrinx/ast"

func parseCode(t *testing.T, sourceCode string) *typeWalker {
	typeWalker := NewTypeWalker().(*typeWalker)

	err := typeWalker.Parse(pkgPath, sourceCode)
	assert.Nil(t, err)
	return typeWalker
}
