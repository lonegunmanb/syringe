package test_code_test

import (
	"github.com/lonegunmanb/syrinx/ioc"
	"github.com/lonegunmanb/syrinx/test_code/engine"
	"github.com/lonegunmanb/syrinx/test_code/fly_car"
	"github.com/lonegunmanb/syrinx/test_code/flyer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateFlyCarCarByManualSimulate(t *testing.T) {
	container := ioc.NewContainer()
	e := engine.NewEngine(100)
	container.RegisterInstance((*engine.Engine)(nil), e)
	decoration := &fly_car.FancyDecoration{}
	container.RegisterInstance((*fly_car.Decoration)(nil), decoration)
	wing := &flyer.AluminumWing{}
	container.RegisterInstance((*flyer.Wing)(nil), wing)
	c := fly_car.Create_fly_car(container)
	assert.Equal(t, e, c.Engine)
	assert.Equal(t, wing, c.Wing)
	assert.Equal(t, decoration, c.Decoration)
}
