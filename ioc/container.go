package ioc

import (
	"fmt"
	"reflect"
)

type Container interface {
	RegisterInstance(interfaceType interface{}, instance interface{})
	Resolve(name string) interface{}
}

type container struct {
	instances map[string]interface{}
	factories map[string]func() interface{}
}

func (c *container) RegisterInstance(interfaceType interface{}, instance interface{}) {
	typeName := getTypeName(interfaceType)
	c.instances[typeName] = instance
}

func (c *container) RegisterFactory(interfaceType interface{}, factory func() interface{}) {
	typeName := getTypeName(interfaceType)
	c.factories[typeName] = factory
}

func (c *container) Resolve(s string) interface{} {
	factory, ok := c.factories[s]
	if ok {
		return factory()
	}
	return c.instances[s]
}

func getTypeName(interfaceType interface{}) string {
	t := reflect.TypeOf(interfaceType).Elem()
	pkgPath := t.PkgPath()
	reflectedTypeName := t.Name()
	if pkgPath != "" {
		reflectedTypeName = fmt.Sprintf("%s.%s", pkgPath, reflectedTypeName)
	}
	return reflectedTypeName
}

//noinspection GoUnusedExportedFunction
func NewContainer() Container {
	return newContainer()
}

func newContainer() *container {
	return &container{
		instances: make(map[string]interface{}),
		factories: make(map[string]func() interface{}),
	}
}
