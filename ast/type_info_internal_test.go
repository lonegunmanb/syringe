package ast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetStructPkgPath(t *testing.T) {
	sourceCode := `
package ast
type Struct struct{
}
`
	walker := parseCode(t, sourceCode)
	assert.Equal(t, pkgPath, walker.Types()[0].PkgPath)
}
