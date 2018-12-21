package ioc

import (
	"github.com/lonegunmanb/syrinx/test_code/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

const engineInterfaceTypeName = "github.com/lonegunmanb/syrinx/test_code/engine.Engine"

var engineInterfaceType interface{} = (*engine.Engine)(nil)

func TestRegisterInstance(t *testing.T) {
	ioc := newContainer()
	e := engine.NewEngine(100)
	ioc.RegisterInstance(engineInterfaceType, e)
	instances := ioc.instances

	resolvedEngine, ok := instances[engineInterfaceTypeName]
	assert.True(t, ok)
	assert.Equal(t, e, resolvedEngine)
}

func TestRegisterFactory(t *testing.T) {
	c := newContainer()
	e := engine.NewEngine(100)
	factory := func() interface{} {
		return e
	}
	c.RegisterFactory(engineInterfaceType, factory)
	resolvedFactory := c.factories[engineInterfaceTypeName]
	assert.NotNil(t, resolvedFactory)
	assert.Equal(t, e, resolvedFactory())
}

func TestResolveWithInstanceOnly(t *testing.T) {
	c := newContainer()
	e := engine.NewEngine(100)
	c.RegisterInstance(engineInterfaceType, e)
	resolvedEngine := c.Resolve(engineInterfaceTypeName)
	assert.Equal(t, e, resolvedEngine)
}

func TestResolveWithInstanceAndFactoryShouldReturnFromFactory(t *testing.T) {
	ioc := newContainer()
	staticEngine := engine.NewEngine(100)

	ioc.RegisterInstance(engineInterfaceType, staticEngine)
	engineFromFactory := engine.NewEngine(200)
	ioc.RegisterFactory(engineInterfaceType, func() interface{} {
		return engineFromFactory
	})
	resolvedEngine := ioc.Resolve(engineInterfaceTypeName)
	assert.Equal(t, engineFromFactory, resolvedEngine)
}
