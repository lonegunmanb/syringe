package codegen

import (
	"github.com/ahmetb/go-linq"
	"github.com/lonegunmanb/syrinx/ast"
)

type GenTask interface {
	GetPkgName() string
	GetDepPkgPaths() []string
	GetTypeInfos() []ast.TypeInfo
}

type genTask struct {
	pkgName   string
	typeInfos []ast.TypeInfo
}

func (c *genTask) GetPkgName() string {
	return c.pkgName
}

func (c *genTask) GetDepPkgPaths() []string {
	paths := make([]string, len(c.typeInfos))
	linq.From(c.typeInfos).SelectMany(func(typeInfo interface{}) linq.Query {
		return linq.From(typeInfo.(ast.TypeInfo).GetDepPkgPaths())
	}).Distinct().ToSlice(&paths)
	return paths
}

func (c *genTask) GetTypeInfos() []ast.TypeInfo {
	return c.typeInfos
}

func NewCodegenTask(pkgName string, typeInfos []ast.TypeInfo) GenTask {
	return &genTask{pkgName: pkgName, typeInfos: typeInfos}
}
