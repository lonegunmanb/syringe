package ioc

import (
	"github.com/lonegunmanb/syrinx/test_code"
	"github.com/stretchr/testify/assert"
	"testing"
)

const engineInterfaceTypeName = "github.com/lonegunmanb/syrinx/test_code.Engine"

var engineInterfaceType interface{} = (*test_code.Engine)(nil)

func TestRegisterInstance(t *testing.T) {
	ioc := newContainer()
	engine := test_code.NewEngine(100)
	ioc.RegisterInstance(engineInterfaceType, engine)
	instances := ioc.instances

	resolvedEngine, ok := instances[engineInterfaceTypeName]
	assert.True(t, ok)
	assert.Equal(t, engine, resolvedEngine)
}

func TestRegisterFactory(t *testing.T) {
	ioc := newContainer()
	engine := test_code.NewEngine(100)
	factory := func() interface{} {
		return engine
	}
	ioc.RegisterFactory(engineInterfaceType, factory)
	resolvedFactory := ioc.factories[engineInterfaceTypeName]
	assert.NotNil(t, resolvedFactory)
	assert.Equal(t, engine, resolvedFactory())
}

func TestResolveWithInstanceOnly(t *testing.T) {
	ioc := newContainer()
	engine := test_code.NewEngine(100)
	ioc.RegisterInstance(engineInterfaceType, engine)
	resolvedEngine := ioc.Resolve(engineInterfaceTypeName)
	assert.Equal(t, engine, resolvedEngine)
}

func TestResolveWithInstanceAndFactoryShouldReturnFromFactory(t *testing.T) {
	ioc := newContainer()
	staticEngine := test_code.NewEngine(100)

	ioc.RegisterInstance(engineInterfaceType, staticEngine)
	engineFromFactory := test_code.NewEngine(200)
	ioc.RegisterFactory(engineInterfaceType, func() interface{} {
		return engineFromFactory
	})
	resolvedEngine := ioc.Resolve(engineInterfaceTypeName)
	assert.Equal(t, engineFromFactory, resolvedEngine)
}
