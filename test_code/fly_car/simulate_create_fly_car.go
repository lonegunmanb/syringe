package fly_car

import (
	"github.com/lonegunmanb/syrinx/ioc"
	"github.com/lonegunmanb/syrinx/test_code/car"
	"github.com/lonegunmanb/syrinx/test_code/flyer"
)

func Create_FlyCar(container ioc.Container) *FlyCar {
	r := new(FlyCar)
	Assemble_FlyCar(r, container)
	return r
}

func Assemble_FlyCar(l *FlyCar, container ioc.Container) {
	l.Car = container.Resolve("github.com/lonegunmanb/syrinx/test_code/car.Car").(*car.Car)
	l.Plane = *container.Resolve("github.com/lonegunmanb/syrinx/test_code/flyer.Plane").(*flyer.Plane)
	l.Decoration = container.Resolve("github.com/lonegunmanb/syrinx/test_code/fly_car.Decoration").(Decoration)
}

func Register_FlyCar(container ioc.Container) {
	container.RegisterFactory((*FlyCar)(nil), func(ioc ioc.Container) interface{} {
		return Create_FlyCar(ioc)
	})
}
