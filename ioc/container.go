package ioc

import (
	"fmt"
	"reflect"
)

type Container interface {
	RegisterFactory(interfaceType interface{}, factory func(ioc Container) interface{})
	Resolve(name string) interface{}
	ResolveByType(interfaceType interface{}) interface{}
	Has(interfaceType interface{}) bool
	GetOrRegister(interfaceType interface{}, factory func(ioc Container) interface{}) interface{}
	UnRegister(name string)
	UnRegisterByType(interfaceType interface{})
	Clear()
}

type container struct {
	factories map[string]func(ioc Container) interface{}
}

func (c *container) Clear() {
	c.factories = make(map[string]func(Container) interface{})
}

func (c *container) UnRegister(name string) {
	delete(c.factories, name)
}

func (c *container) UnRegisterByType(interfaceType interface{}) {
	c.UnRegister(getTypeName(interfaceType))
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

func (c *container) GetOrRegister(interfaceType interface{}, factory func(ioc Container) interface{}) interface{} {
	if !c.Has(interfaceType) {
		c.RegisterFactory(interfaceType, factory)
	}
	return c.ResolveByType(interfaceType)
}

func (c *container) ResolveByType(interfaceType interface{}) interface{} {
	return c.Resolve(getTypeName(interfaceType))
}

func (c *container) Has(interfaceType interface{}) bool {
	_, ok := c.factories[getTypeName(interfaceType)]
	return ok
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
