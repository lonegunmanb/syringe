package codegen

import "fmt"

type RegisterCodeWriter interface {
	RegisterCode() string
}

type registerCodeWriter struct {
	typeInfo       TypeInfoWrap
	workingPkgPath string
}

func (r *registerCodeWriter) RegisterCode() string {
	typePkgPath := r.typeInfo.GetPkgPath()
	if !samePackage(typePkgPath, r.workingPkgPath) {
		return fmt.Sprintf("%s.Register_%s(%s)", r.typeInfo.GetPkgNameFromPkgPath(r.typeInfo.GetPkgPath()), r.typeInfo.GetName(), ContainerIdentName)
	}
	return fmt.Sprintf("Register_%s(%s)", r.typeInfo.GetName(), ContainerIdentName)
}
