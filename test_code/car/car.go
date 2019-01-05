package car

import "github.com/lonegunmanb/syringe/test_code/engine"

type Car struct {
	Engine engine.Engine `inject:""`
}
