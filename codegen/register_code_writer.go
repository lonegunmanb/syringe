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
	t := r.typeInfo
	if !samePackage(typePkgPath, r.workingPkgPath) {
		return fmt.Sprintf("%s.Register_%s(%s)",
			t.GetPkgNameFromPkgPath(t.GetPkgPath()), t.GetName(), ContainerIdentName)
	}
	return fmt.Sprintf("Register_%s(%s)", t.GetName(), ContainerIdentName)
}
