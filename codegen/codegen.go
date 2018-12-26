package codegen

import (
	"github.com/lonegunmanb/syrinx/ast"
	"io"
)

type codegen struct {
	typeInfo ast.TypeInfo
	writer   io.Writer
}
