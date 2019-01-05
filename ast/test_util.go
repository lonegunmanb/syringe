package ast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const testPkgPath = "github.com/lonegunmanb/syringe/ast"

func parseCode(t *testing.T, sourceCode string) *typeWalker {
	typeWalker := NewTypeWalker().(*typeWalker)

	err := typeWalker.Parse(testPkgPath, sourceCode)
	assert.Nil(t, err)
	return typeWalker
}
