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
