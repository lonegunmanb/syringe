package codegen

import "go/types"

var qf = func(p *types.Package) string {
	return p.Name()
}
