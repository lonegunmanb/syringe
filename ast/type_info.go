package ast

import (
	"go/types"
	"reflect"
)

type TypeInfo struct {
	Name    string
	PkgPath string
	Fields  []*FieldInfo
	Type    types.Type
	Kind    reflect.Kind
}
