package car

import (
	"github.com/lonegunmanb/syrinx/ioc"
	"github.com/lonegunmanb/syrinx/test_code/engine"
)

//noinspection GoUnusedExportedFunction
func Create_car(container ioc.Container) *Car {
	r := new(Car)
	Assemble_car(r, container)
	return r
}

func Assemble_car(c *Car, container ioc.Container) {
	c.Engine = container.Resolve("github.com/lonegunmanb/syrinx/test_code/engine.Engine").(engine.Engine)
}
