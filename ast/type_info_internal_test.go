package ast

import (
	"github.com/ahmetb/go-linq"
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

func TestGetStructDepImportPkgPaths(t *testing.T) {
	sourceCode := `
package ast
import (
"go/scanner"
"go/token"
)
type Struct struct {
	Err scanner.Error //dep go/scanner
	FileSet token.FileSet //dep go/token
	PErr *scanner.Error //dep go/scanner
	PFileSet *token.FileSet //dep go/token
	SErr []scanner.Error //dep go/scanner
	SFileSet []token.FileSet //dep go/token
	AErr [1]scanner.Error //dep go/scanner
	AFileSet [1]token.FileSet //dep go/token
	NestedStruct struct { //dep pkgPath here
		Name string
	}
	NestedInterface interface { //dep pkgPath here
		Name() string
	}
	Map map[*scanner.Error]*token.FileSet //dep go/scanner, go/token
}
`
	walker := parseCode(t, sourceCode)
	structInfo := walker.Types()[0]
	depPkgPaths := structInfo.DepPkgPaths()
	assert.Equal(t, 3, len(depPkgPaths))
	assert.True(t, linq.From(depPkgPaths).Contains("go/scanner"))
	assert.True(t, linq.From(depPkgPaths).Contains("go/token"))
	assert.True(t, linq.From(depPkgPaths).Contains(pkgPath))
}
