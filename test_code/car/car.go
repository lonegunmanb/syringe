package car

import "github.com/lonegunmanb/syrinx/test_code/engine"

type Car struct {
	Engine engine.Engine `inject:""`
}
