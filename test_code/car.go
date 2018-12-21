package test_code

type HorsePower = int

type Engine interface {
	OutputPower() HorsePower
}

type pistonEngin struct {
	horsePower HorsePower
}

func (e *pistonEngin) OutputPower() HorsePower {
	return e.horsePower
}

func NewEngine(horsePower HorsePower) Engine {
	return &pistonEngin{horsePower: horsePower}
}
