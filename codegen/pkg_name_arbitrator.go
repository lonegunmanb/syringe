package codegen

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/golang-collections/collections/set"
	"github.com/lonegunmanb/varys/ast"
	"strings"
)

type DepMode = string

var ProductCodegenMode = "ProductCodegen"
var RegisterCodegenMode = "RegisterCodegen"

type PkgNameArbitrator interface {
	GenImportDecls() []string
	GetPkgNameFromPkgPath(pkgPath string) string
}

type packageNameCounter struct {
	count       int
	nameExisted *set.Set
}

func newNameCounter() *packageNameCounter {
	return &packageNameCounter{
		count:       0,
		nameExisted: set.New(),
	}
}

func (n *packageNameCounter) generateName() string {
	for {
		nextName := fmt.Sprintf("p%d", n.count)
		n.count++
		if !n.nameExisted.Has(nextName) {
			n.nameExisted.Insert(nextName)
			return nextName
		}
	}
}

func (n *packageNameCounter) uniqueName(name string) string {
	if n.nameExisted.Has(name) {
		return n.generateName()
	}
	n.nameExisted.Insert(name)
	return name
}

type pkgNameArbitrator struct {
	typeInfos            []ast.TypeInfo
	depPkgPaths          []string
	depPkgPathPkgNameMap map[string]string
	currentPkgPath       string
	mode                 DepMode
}

func newPkgNameArbitrator(typeInfos []ast.TypeInfo, currentPkgPath string, mode DepMode) PkgNameArbitrator {
	return &pkgNameArbitrator{typeInfos: typeInfos, currentPkgPath: currentPkgPath, mode: mode}
}

func (c *pkgNameArbitrator) GenImportDecls() []string {
	paths := c.getDepPkgPaths()
	results := make([]string, 0, len(c.depPkgPathPkgNameMap))
	//we iterate paths so generated import decls' order is as same as fields' order
	for _, pkgPath := range paths {
		if pkgPath == "" {
			panic("no package path found, maybe forgot to set walker's physical path?")
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

func (c *pkgNameArbitrator) GetPkgNameFromPkgPath(pkgPath string) string {
	name, ok := c.depPkgPathPkgNameMap[pkgPath]
	if !ok {
		name = retrievePkgNameFromPkgPath(pkgPath)
	}
	return name
}

func (c *pkgNameArbitrator) getDepPkgPaths() []string {
	if c.depPkgPaths != nil {
		return c.depPkgPaths
	}
	c.depPkgPaths = c.initDepPkgPaths()
	c.depPkgPathPkgNameMap = c.initDepPkgPathPkgNameMap()
	return c.depPkgPaths
}

func (c *pkgNameArbitrator) initDepPkgPaths() []string {
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
		return path.(string) != c.currentPkgPath
	}).ToSlice(&paths)
	return paths
}

func (c *pkgNameArbitrator) initDepPkgPathPkgNameMap() map[string]string {
	pkgNamePkgPathMap := make(map[string][]string)
	packageNameCounter := newNameCounter()
	for _, path := range c.depPkgPaths {
		pkgName := retrievePkgNameFromPkgPath(path)
		paths := pkgNamePkgPathMap[pkgName]
		pkgNamePkgPathMap[pkgName] = append(paths, path)
	}
	pkgPathPkgNameMap := make(map[string]string)
	for pkgName, paths := range pkgNamePkgPathMap {
		if len(paths) == 1 && pkgName != "ioc" {
			pkgPathPkgNameMap[paths[0]] = packageNameCounter.uniqueName(pkgName)
		} else {
			for _, path := range paths {
				pkgPathPkgNameMap[path] = packageNameCounter.generateName()
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
