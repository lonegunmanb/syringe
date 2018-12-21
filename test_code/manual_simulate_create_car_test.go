package test_code_test

import (
	"github.com/lonegunmanb/syrinx/ioc"
	"github.com/lonegunmanb/syrinx/test_code/car"
	"github.com/lonegunmanb/syrinx/test_code/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateCarByManualSimulate(t *testing.T) {
	container := ioc.NewContainer()
	e := engine.NewEngine(100)
	container.RegisterInstance((*engine.Engine)(nil), e)
	c := car.Create_car(container)
	assert.Equal(t, e, c.Engine)
}
