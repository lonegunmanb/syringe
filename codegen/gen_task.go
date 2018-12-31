package codegen

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/lonegunmanb/syrinx/ast"
	"strings"
)

type GenTask interface {
	GetPkgName() string
	GetDepPkgPaths() []string
	GenImportDecls() []string
	GetTypeInfos() []ast.TypeInfo
	GetPkgNameFromPkgPath(pkgPath string) string
}

type genTask struct {
	pkgName              string
	typeInfos            []ast.TypeInfo
	depPkgPaths          []string
	depPkgPathPkgNameMap map[string]string
}

func (c *genTask) GenImportDecls() []string {
	paths := c.GetDepPkgPaths()
	results := make([]string, 0, len(c.depPkgPathPkgNameMap))
	//we iterate paths so generated import decls' order is as same as fields' order
	for _, pkgPath := range paths {
		pkgName := c.depPkgPathPkgNameMap[pkgPath]
		if pkgName != getPkgNameFromPkgPath(pkgPath) {
			results = append(results, fmt.Sprintf(`%s "%s"`, pkgName, pkgPath))
		} else {
			results = append(results, fmt.Sprintf(`"%s"`, pkgPath))
		}
	}
	return results
}

func (c *genTask) GetPkgNameFromPkgPath(pkgPath string) string {
	name, ok := c.depPkgPathPkgNameMap[pkgPath]
	if !ok {
		name = getPkgNameFromPkgPath(pkgPath)
	}
	return name
}

func (c *genTask) GetPkgName() string {
	return c.pkgName
}

func (c *genTask) GetDepPkgPaths() []string {
	if c.depPkgPaths != nil {
		return c.depPkgPaths
	}
	c.depPkgPaths = c.initDepPkgPaths()
	c.depPkgPathPkgNameMap = c.initDepPkgPathPkgNameMap()
	return c.depPkgPaths
}

func (c *genTask) GetTypeInfos() []ast.TypeInfo {
	return c.typeInfos
}

func NewCodegenTask(pkgName string, typeInfos []ast.TypeInfo) GenTask {
	return &genTask{pkgName: pkgName, typeInfos: typeInfos}
}

func (c *genTask) initDepPkgPaths() []string {
	paths := make([]string, len(c.typeInfos))
	linq.From(c.typeInfos).SelectMany(func(typeInfo interface{}) linq.Query {
		return linq.From(typeInfo.(ast.TypeInfo).GetDepPkgPaths())
	}).Distinct().ToSlice(&paths)
	return paths
}

func (c *genTask) initDepPkgPathPkgNameMap() map[string]string {
	pkgNamePkgPathMap := make(map[string][]string)
	for _, path := range c.depPkgPaths {
		pkgName := getPkgNameFromPkgPath(path)
		paths := pkgNamePkgPathMap[pkgName]
		pkgNamePkgPathMap[pkgName] = append(paths, path)
	}
	count := 0
	pkgPathPkgNameMap := make(map[string]string)
	for pkgName, paths := range pkgNamePkgPathMap {
		if len(paths) == 1 && pkgName != "ioc" {
			pkgPathPkgNameMap[paths[0]] = pkgName
		} else {
			for _, path := range paths {
				pkgPathPkgNameMap[path] = fmt.Sprintf("p%d", count)
				count++
			}
		}
	}
	return pkgPathPkgNameMap
}

func getPkgNameFromPkgPath(pkgPath string) string {
	if !strings.Contains(pkgPath, "/") {
		return pkgPath
	}
	s := strings.Split(pkgPath, "/")
	return s[len(s)-1]
}
