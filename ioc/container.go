package ioc

import (
	"fmt"
	"reflect"
)

type Container interface {
	RegisterFactory(interfaceType interface{}, factory func(ioc Container) interface{})
	Resolve(name string) interface{}
}

type container struct {
	factories map[string]func(ioc Container) interface{}
}

func (c *container) RegisterFactory(interfaceType interface{}, factory func(ioc Container) interface{}) {
	typeName := getTypeName(interfaceType)
	c.factories[typeName] = factory
}

func (c *container) Resolve(s string) interface{} {
	factory, ok := c.factories[s]
	if ok {
		return factory(c)
	}
	return nil
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
		factories: make(map[string]func(ioc Container) interface{}),
	}
}
