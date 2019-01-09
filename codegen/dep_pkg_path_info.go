package codegen

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/golang-collections/collections/set"
	"github.com/lonegunmanb/varys/ast"
	"strings"
)

type DepMode = string

var ProductCodegenMode DepMode = "ProductCodegen"
var RegisterCodegenMode DepMode = "RegisterCodegen"

type DepPkgPathInfo interface {
	GetDepPkgPaths() []string
	GenImportDecls() []string
	GetPkgNameFromPkgPath(pkgPath string) string
}

type packageNamer struct {
	count       int
	nameExisted *set.Set
}

func newNamer() *packageNamer {
	return &packageNamer{
		count:       0,
		nameExisted: set.New(),
	}
}

func (n *packageNamer) nextName() string {
	for {
		nextName := fmt.Sprintf("p%d", n.count)
		n.count++
		if !n.nameExisted.Has(nextName) {
			n.nameExisted.Insert(nextName)
			return nextName
		}
	}
}

func (n *packageNamer) uniqueName(name string) string {
	if n.nameExisted.Has(name) {
		return n.nextName()
	}
	n.nameExisted.Insert(name)
	return name
}

type depPkgPathInfo struct {
	typeInfos            []ast.TypeInfo
	depPkgPaths          []string
	depPkgPathPkgNameMap map[string]string
	pkgPath              string
	mode                 DepMode
}

func newDepPkgPathInfo(typeInfos []ast.TypeInfo, pkgPath string, mode DepMode) DepPkgPathInfo {
	return &depPkgPathInfo{typeInfos: typeInfos, pkgPath: pkgPath, mode: mode}
}

func (c *depPkgPathInfo) GenImportDecls() []string {
	paths := c.GetDepPkgPaths()
	results := make([]string, 0, len(c.depPkgPathPkgNameMap))
	//we iterate paths so generated import decls' order is as same as fields' order
	for _, pkgPath := range paths {
		if pkgPath == "" {
			continue
		}
		pkgName := c.depPkgPathPkgNameMap[pkgPath]
		if pkgName != retrievePkgNameFromPkgPath(pkgPath) {
			results = append(results, fmt.Sprintf(`%s "%s"`, pkgName, pkgPath))
		} else {
			results = append(results, fmt.Sprintf(`"%s"`, pkgPath))
		}
	}
	return results
}

func (c *depPkgPathInfo) GetPkgNameFromPkgPath(pkgPath string) string {
	name, ok := c.depPkgPathPkgNameMap[pkgPath]
	if !ok {
		name = retrievePkgNameFromPkgPath(pkgPath)
	}
	return name
}

func (c *depPkgPathInfo) GetDepPkgPaths() []string {
	if c.depPkgPaths != nil {
		return c.depPkgPaths
	}
	c.depPkgPaths = c.initDepPkgPaths()
	c.depPkgPathPkgNameMap = c.initDepPkgPathPkgNameMap()
	return c.depPkgPaths
}

func (c *depPkgPathInfo) GetTypeInfos() []ast.TypeInfo {
	return c.typeInfos
}

func (c *depPkgPathInfo) initDepPkgPaths() []string {
	paths := make([]string, len(c.typeInfos))
	query := linq.From(c.typeInfos).Select(func(typeInfo interface{}) interface{} {
		return typeInfo.(ast.TypeInfo).GetPkgPath()
	})
	if c.mode == ProductCodegenMode {
		query = query.Concat(
			linq.From(c.typeInfos).SelectMany(func(typeInfo interface{}) linq.Query {
				return linq.From(typeInfo.(ast.TypeInfo).GetDepPkgPaths("inject"))
			}))
	}
	query.Distinct().Where(func(path interface{}) bool {
		return path.(string) != c.pkgPath
	}).ToSlice(&paths)
	return paths
}

func (c *depPkgPathInfo) initDepPkgPathPkgNameMap() map[string]string {
	pkgNamePkgPathMap := make(map[string][]string)
	packageNamer := newNamer()
	for _, path := range c.depPkgPaths {
		pkgName := retrievePkgNameFromPkgPath(path)
		paths := pkgNamePkgPathMap[pkgName]
		pkgNamePkgPathMap[pkgName] = append(paths, path)
	}
	pkgPathPkgNameMap := make(map[string]string)
	for pkgName, paths := range pkgNamePkgPathMap {
		if len(paths) == 1 && pkgName != "ioc" {
			pkgPathPkgNameMap[paths[0]] = packageNamer.uniqueName(pkgName)
		} else {
			for _, path := range paths {
				pkgPathPkgNameMap[path] = packageNamer.nextName()
			}
		}
	}
	return pkgPathPkgNameMap
}

func retrievePkgNameFromPkgPath(pkgPath string) string {
	if !strings.Contains(pkgPath, "/") {
		return pkgPath
	}
	s := strings.Split(pkgPath, "/")
	return s[len(s)-1]
}
