package fly_car

import (
	"github.com/lonegunmanb/syrinx/ioc"
	"github.com/lonegunmanb/syrinx/test_code/car"
	"github.com/lonegunmanb/syrinx/test_code/flyer"
)

func Create_fly_car(container ioc.Container) *FlyCar {
	r := new(FlyCar)
	Assemble_fly_car(r, container)
	return r
}

func Assemble_fly_car(l *FlyCar, container ioc.Container) {
	l.Car = new(car.Car)
	car.Assemble_car(l.Car, container)
	flyer.Assemble_plane(&l.Plane, container)
	l.Decoration = container.Resolve("github.com/lonegunmanb/syrinx/test_code/fly_car.Decoration").(Decoration)
}
