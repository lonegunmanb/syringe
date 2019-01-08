package ioc

import (
	"github.com/lonegunmanb/syringe/test_code/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

const engineInterfaceTypeName = "github.com/lonegunmanb/syringe/test_code/engine.Engine"

var engineInterfaceType interface{} = (*engine.Engine)(nil)

func TestRegisterFactory(t *testing.T) {
	c := newContainer()
	e := engine.NewEngine(100)
	factory := func(ioc Container) interface{} {
		return e
	}
	c.RegisterFactory(engineInterfaceType, factory)
	resolvedFactory := c.factories[engineInterfaceTypeName]
	assert.NotNil(t, resolvedFactory)
	assert.Equal(t, e, resolvedFactory(c))
}

type stubEngine struct {
}

func (stubEngine) OutputPower() engine.HorsePower {
	panic("implement me")
}

func TestContainer_RegisterTwoTypes(t *testing.T) {
	c := newContainer()
	e := &stubEngine{}
	c.RegisterFactory((*stubEngine)(nil), func(ioc Container) interface{} {
		return e
	})
	c.RegisterTwoTypes((*engine.Engine)(nil), (*stubEngine)(nil))
	resolvedEngine := c.ResolveByType((*engine.Engine)(nil))
	assert.Equal(t, e, resolvedEngine)
}
