package codegen

import "fmt"

type Register interface {
	RegisterCode() string
}

type register struct {
	typeInfo        TypeInfoWrap
	registeringPath string
}

func (r *register) RegisterCode() string {
	pkgPath := r.typeInfo.GetPkgPath()
	if !samePackage(pkgPath, r.registeringPath) {
		return fmt.Sprintf("%s.Register_%s(container)", r.typeInfo.GetPkgNameFromPkgPath(r.typeInfo.GetPkgPath()), r.typeInfo.GetName())
	}
	return fmt.Sprintf("Register_%s(container)", r.typeInfo.GetName())
}
