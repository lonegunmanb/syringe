package fly_car

import (
	"github.com/lonegunmanb/syrinx/ioc"
	"github.com/lonegunmanb/syrinx/test_code/car"
	"github.com/lonegunmanb/syrinx/test_code/flyer"
)

func Create_fly_car_FlyCar(container ioc.Container) *FlyCar {
	product := new(FlyCar)
	Assemble_fly_car_FlyCar(product, container)
	return product
}
func Assemble_fly_car_FlyCar(product *FlyCar, container ioc.Container) {
	product.Car = container.Resolve("github.com/lonegunmanb/syrinx/test_code/car.Car").(*car.Car)
	product.Plane = *container.Resolve("github.com/lonegunmanb/syrinx/test_code/flyer.Plane").(*flyer.Plane)
	product.Decoration = container.Resolve("github.com/lonegunmanb/syrinx/test_code/fly_car.Decoration").(Decoration)
	//product.r1 = container.Resolve("github.com/lonegunmanb/syrinx/test_code/check_package_name_duplicate_a/model.Request").(*p1.Request)
	//product.r2 = container.Resolve("github.com/lonegunmanb/syrinx/test_code/check_package_name_duplicate_b/model.Request").(*p2.Request)
}
func Register_FlyCar(container ioc.Container) {
	container.RegisterFactory((*FlyCar)(nil), func(ioc ioc.Container) interface{} {
		return Create_fly_car_FlyCar(ioc)
	})
}
