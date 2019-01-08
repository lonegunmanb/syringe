package ast

import (
	"github.com/ahmetb/go-linq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStructName(t *testing.T) {
	sourceCode := `
package ast
type Struct struct{
}
`
	walker := parseCode(t, sourceCode)
	typeInfo := walker.Types()[0]
	assert.Equal(t, "Struct", typeInfo.Name)
}

func TestGetStructPkgPath(t *testing.T) {
	sourceCode := `
package ast
type Struct struct{
}
`
	walker := parseCode(t, sourceCode)
	typeInfo := walker.Types()[0]
	assert.Equal(t, testPkgPath, typeInfo.PkgPath)
	assert.Equal(t, "ast", typeInfo.PkgName)
}

func TestDepSamePkg(t *testing.T) {
	sourceCode := `
package ast
type Struct struct {
s Struct2
}
type Struct2 struct {

}
`
	walker := parseCode(t, sourceCode)
	typeInfo := walker.Types()[0]
	assert.Equal(t, testPkgPath, typeInfo.PkgPath)
}

func TestPackageNameDifferentWithPkgPath(t *testing.T) {
	sourceCode := `
package test
type Struct struct{
}
`
	walker := parseCode(t, sourceCode)
	typeInfo := walker.Types()[0]
	assert.Equal(t, testPkgPath, typeInfo.PkgPath)
	assert.Equal(t, "test", typeInfo.PkgName)
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
	depPkgPaths := structInfo.GetDepPkgPaths("")
	assert.Equal(t, 2, len(depPkgPaths))
	assert.True(t, linq.From(depPkgPaths).Contains("go/scanner"))
	assert.True(t, linq.From(depPkgPaths).Contains("go/token"))
}

func TestIsNotTestFile(t *testing.T) {
	assert.True(t, isTestFile("rover_test.go"))
	assert.True(t, isTestFile("test.go"))
	assert.False(t, isTestFile("rover.go"))
}

func TestIsGoFile(t *testing.T) {
	assert.True(t, isGoSrcFile("src.go"))
	assert.False(t, isGoSrcFile("src.cpp"))
	assert.False(t, isGoSrcFile("go"))
}
