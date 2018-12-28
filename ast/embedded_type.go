package ast

import "go/types"

type EmbeddedKind string

const EmbeddedByStruct EmbeddedKind = "EmbeddedByStruct"
const EmbeddedByPointer EmbeddedKind = "EmbeddedByPointer"
const EmbeddedByInterface EmbeddedKind = "EmbeddedByInterface"

type EmbeddedType interface {
	GetFullName() string
	GetPkgPath() string
	GetKind() EmbeddedKind
	GetTag() string
	GetType() types.Type
	GetReferenceFrom() TypeInfo
}

type embeddedType struct {
	FullName      string
	PkgPath       string
	Kind          EmbeddedKind
	Tag           string
	Type          types.Type
	ReferenceFrom TypeInfo
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

func (e *embeddedType) GetType() types.Type {
	return e.Type
}

func (e *embeddedType) GetReferenceFrom() TypeInfo {
	return e.ReferenceFrom
}
