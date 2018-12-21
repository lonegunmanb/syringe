package car

import (
	"github.com/lonegunmanb/syrinx/ioc"
	"github.com/lonegunmanb/syrinx/test_code/engine"
)

func Create_car(container ioc.Container) *Car {
	r := new(Car)
	r.Engine = container.Resolve("github.com/lonegunmanb/syrinx/test_code/engine.Engine").(engine.Engine)
	return r
}
