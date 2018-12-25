package ast

type EmbeddedKind string

const EmbeddedByStruct EmbeddedKind = "EmbeddedByStruct"
const EmbeddedByPointer EmbeddedKind = "EmbeddedByPointer"
const EmbeddedByInterface EmbeddedKind = "EmbeddedByInterface"

type EmbeddedType struct {
	FullName string
	PkgPath  string
	Kind     EmbeddedKind
	Tag      string
}
