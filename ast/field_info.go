package ast

import "go/types"

type FieldInfo struct {
	Name string
	Type types.Type
	Tag  string
}
