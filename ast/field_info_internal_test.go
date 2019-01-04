package ast

import (
	"github.com/ahmetb/go-linq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetNamedTypeStructFieldPkgPath(t *testing.T) {
	sourceCode := `
package ast
import (
"go/scanner"
"go/token"
)
type Struct struct {
	Err scanner.Error
	FileSet token.FileSet
}
`
	walker := parseCode(t, sourceCode)
	structInfo := walker.Types()[0]
	field1 := structInfo.Fields[0]
	field2 := structInfo.Fields[1]
	fieldPackages := field1.GetDepPkgPaths()
	assert.Equal(t, 1, len(fieldPackages))
	assert.Equal(t, "go/scanner", fieldPackages[0])
	fieldPackages = field2.GetDepPkgPaths()
	assert.Equal(t, 1, len(fieldPackages))
	assert.Equal(t, "go/token", fieldPackages[0])
}

func TestGetPointerToNamedTypeStructFieldPkgPath(t *testing.T) {
	sourceCode := `
package ast
import (
"go/scanner"
)
type Struct struct {
	Err *scanner.Error
}
`
	walker := parseCode(t, sourceCode)
	assertPkgPath(t, walker, "go/scanner")
}

func TestGetSliceOfNamedTypeStructFieldPkgPath(t *testing.T) {
	sourceCode := `
package ast
import (
"go/types"
)
type Struct struct {
	_type types.Type
}
`
	walker := parseCode(t, sourceCode)
	assertPkgPath(t, walker, "go/types")
}

func TestGetArrayOfNamedTypeStructFieldPkgPath(t *testing.T) {
	sourceCode := `
package ast
import (
"go/scanner"
)
type Struct struct {
	Err [1]scanner.Error
}
`
	walker := parseCode(t, sourceCode)
	assertPkgPath(t, walker, "go/scanner")
}

func TestGetSliceOfPointerOfNamedTypeStructFieldPkgPath(t *testing.T) {
	sourceCode := `
package ast
import (
"go/scanner"
)
type Struct struct {
	Err []*scanner.Error
}
`
	walker := parseCode(t, sourceCode)
	assertPkgPath(t, walker, "go/scanner")
}

func TestGetMapFieldPkgPath(t *testing.T) {
	sourceCode := `
package ast
import (
"go/scanner"
"go/token"
)
type Struct struct {
	Data map[*token.FileSet]*scanner.Error
}
`
	walker := parseCode(t, sourceCode)
	structInfo := walker.Types()[0]
	field1 := structInfo.Fields[0]
	packagePaths := field1.GetDepPkgPaths()
	assert.Equal(t, 2, len(packagePaths))
	assert.Equal(t, "go/token", packagePaths[0])
	assert.Equal(t, "go/scanner", packagePaths[1])
}

func TestGetMapFieldPkgPath2(t *testing.T) {
	sourceCode := `
package ast
import (
"go/scanner"
"go/token"
"go/types"
)
type Struct struct {
	Data map[*token.FileSet]map[*scanner.Error]*types.Type
}
`
	walker := parseCode(t, sourceCode)
	structInfo := walker.Types()[0]
	field1 := structInfo.Fields[0]
	packagePaths := field1.GetDepPkgPaths()
	assert.Equal(t, 3, len(packagePaths))
	assert.Equal(t, "go/token", packagePaths[0])
	assert.Equal(t, "go/scanner", packagePaths[1])
	assert.Equal(t, "go/types", packagePaths[2])
}

func TestNestedStructFieldPkgPath(t *testing.T) {
	sourceCode := `
package ast

type Struct struct {
	Field struct{
		Name string
	}
}
`
	walker := parseCode(t, sourceCode)
	assertPkgPath(t, walker, testPkgPath)
}

func TestNestedInterfaceFieldPkgPath(t *testing.T) {
	sourceCode := `
package ast

type Struct struct {
	Field interface {
		Name() string
	}
}
`
	walker := parseCode(t, sourceCode)
	assertPkgPath(t, walker, testPkgPath)
}

func TestBasicFieldPkgPath(t *testing.T) {
	sourceCode := `
package ast

type Struct struct {
	Name string
}
`
	walker := parseCode(t, sourceCode)
	structInfo := walker.Types()[0]
	field1 := structInfo.Fields[0]
	assert.Equal(t, 0, len(field1.GetDepPkgPaths()))
}

func TestChanFieldPkgPath(t *testing.T) {
	sourceCode := `
package ast
import "go/token" 
type Struct struct {
	FileSetChan chan *token.FileSet
}
`
	walker := parseCode(t, sourceCode)
	assertPkgPath(t, walker, "go/token")
}

func TestFuncFieldPkgPath(t *testing.T) {
	sourceCode := `
package ast
import (
"go/scanner"
"go/token"
"go/types"
)
type Struct struct {
	Func func(fileSet *token.FileSet, e *scanner.Error) (types.Type, error)
}
`
	walker := parseCode(t, sourceCode)
	structInfo := walker.Types()[0]
	depPkgPaths := structInfo.GetDepPkgPaths("")
	assert.Equal(t, 3, len(depPkgPaths))
	assert.True(t, linq.From(depPkgPaths).Contains("go/scanner"))
	assert.True(t, linq.From(depPkgPaths).Contains("go/token"))
	assert.True(t, linq.From(depPkgPaths).Contains("go/types"))
}

func TestFuncEmbeddedTypePkgPath(t *testing.T) {
	sourceCode := `
package ast
import (
"go/types"
)
type Struct struct {
	types.Named
}
`
	walker := parseCode(t, sourceCode)
	structInfo := walker.Types()[0]
	depPkgPaths := structInfo.GetDepPkgPaths("")
	assert.Equal(t, 1, len(depPkgPaths))
	dep := depPkgPaths[0]
	assert.Equal(t, "go/types", dep)
}

func assertPkgPath(t *testing.T, walker *typeWalker, pkgPath string) {
	structInfo := walker.Types()[0]
	field1 := structInfo.Fields[0]
	packagePaths := field1.GetDepPkgPaths()
	assert.Equal(t, 1, len(packagePaths))
	assert.Equal(t, pkgPath, packagePaths[0])
}
