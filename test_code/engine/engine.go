package engine

type HorsePower = int

type Engine interface {
	OutputPower() HorsePower
}

type pistonEngine struct {
	horsePower HorsePower
}

func (e *pistonEngine) OutputPower() HorsePower {
	return e.horsePower
}

func NewEngine(horsePower HorsePower) Engine {
	return &pistonEngine{horsePower: horsePower}
}
