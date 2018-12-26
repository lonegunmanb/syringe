package ast

type EmbeddedKind string

const EmbeddedByStruct EmbeddedKind = "EmbeddedByStruct"
const EmbeddedByPointer EmbeddedKind = "EmbeddedByPointer"
const EmbeddedByInterface EmbeddedKind = "EmbeddedByInterface"

type EmbeddedType interface {
	GetName() string
	GetFullName() string
	GetPkgPath() string
	GetKind() EmbeddedKind
	GetTag() string
}

type embeddedType struct {
	Name     string
	FullName string
	PkgPath  string
	Kind     EmbeddedKind
	Tag      string
}

func (e *embeddedType) GetName() string {
	return e.Name
}

func (e *embeddedType) GetFullName() string {
	return e.FullName
}

func (e *embeddedType) GetPkgPath() string {
	return e.PkgPath
}

func (e *embeddedType) GetKind() EmbeddedKind {
	return e.Kind
}

func (e *embeddedType) GetTag() string {
	return e.Tag
}
