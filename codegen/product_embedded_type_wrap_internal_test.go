package codegen

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const embeddedAssignTemplate = `
package test
%s
type Struct struct {
%s
}
type Struct2 struct {
}
type Interface interface {
Name() string
}
`

func TestEmbedded(t *testing.T) {
	testEmbeddedType(t, "", "Struct2", "Struct2", "*", "github.com/lonegunmanb/syringe/test.Struct2", "*Struct2")
	testEmbeddedType(t, "", "*Struct2", "Struct2", "", "github.com/lonegunmanb/syringe/test.Struct2", "*Struct2")
	testEmbeddedType(t, "", "Interface", "Interface", "", "github.com/lonegunmanb/syringe/test.Interface", "Interface")
	importDecl := `import "github.com/lonegunmanb/syringe/test_code/car"`
	testEmbeddedType(t, importDecl, "car.Car", "Car", "*", "github.com/lonegunmanb/syringe/test_code/car.Car", "*car.Car")
	testEmbeddedType(t, importDecl, "*car.Car", "Car", "", "github.com/lonegunmanb/syringe/test_code/car.Car", "*car.Car")
	importDecl = `import "github.com/lonegunmanb/syringe/test_code/engine"`
	testEmbeddedType(t, importDecl, "engine.Engine", "Engine", "", "github.com/lonegunmanb/syringe/test_code/engine.Engine", "engine.Engine")
}

type stubTypeCodegen struct{}

func (*stubTypeCodegen) GetPkgPath() string {
	panic("implement me")
}

func (*stubTypeCodegen) GenImportDecls() []string {
	panic("implement me")
}

func (*stubTypeCodegen) GetName() string {
	panic("implement me")
}

func (*stubTypeCodegen) GetPkgName() string {
	panic("implement me")
}

func (*stubTypeCodegen) GetDepPkgPaths(fieldTagFilter string) []string {
	panic("implement me")
}

func (*stubTypeCodegen) GetFieldAssigns() []Assembler {
	panic("implement me")
}

func (*stubTypeCodegen) GetEmbeddedTypeAssigns() []Assembler {
	panic("implement me")
}

func (*stubTypeCodegen) GetPkgNameFromPkgPath(pkgPath string) string {
	return retrievePkgNameFromPkgPath(pkgPath)
}

func testEmbeddedType(t *testing.T, importString string, typeDecl string, assignedField string, star string, key string, convertType string) {
	code := fmt.Sprintf(embeddedAssignTemplate, importString, typeDecl)
	walker := parseCode(t, code)
	embeddedType := &productEmbeddedTypeWrap{walker.GetTypes()[0].GetEmbeddedTypes()[0], &stubTypeCodegen{}}
	expected := fmt.Sprintf(`product.%s = %scontainer.Resolve("%s").(%s)`, assignedField, star, key, convertType)
	assert.Equal(t, expected, embeddedType.AssembleCode())
}

func TestEmbeddedTypeAssignWithCustomIdent(t *testing.T) {
	code := `
package test

type Struct struct {
Struct2
}
type Struct2 struct {
}
`
	walker := parseCode(t, code)
	embeddedType := &productEmbeddedTypeWrap{walker.GetTypes()[0].GetEmbeddedTypes()[0], &stubTypeCodegen{}}
	expected := `product.Struct2 = *c.Resolve("github.com/lonegunmanb/syringe/test.Struct2").(*Struct2)`
	originIdent := ContainerIdentName
	ContainerIdentName = "c"
	defer func() {
		ContainerIdentName = originIdent
	}()
	assert.Equal(t, expected, embeddedType.AssembleCode())
}
