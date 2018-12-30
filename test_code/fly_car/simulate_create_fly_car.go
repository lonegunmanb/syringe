package fly_car

import (
	"github.com/lonegunmanb/syrinx/ioc"
	"github.com/lonegunmanb/syrinx/test_code/car"
	"github.com/lonegunmanb/syrinx/test_code/flyer"
)

func Create_FlyCar(container ioc.Container) *FlyCar {
	product := new(FlyCar)
	Assemble_FlyCar(product, container)
	return product
}
func Assemble_FlyCar(product *FlyCar, container ioc.Container) {
	product.Car = container.Resolve("github.com/lonegunmanb/syrinx/test_code/car.Car").(*car.Car)
	product.Plane = *container.Resolve("github.com/lonegunmanb/syrinx/test_code/flyer.Plane").(*flyer.Plane)
	product.Decoration = container.Resolve("github.com/lonegunmanb/syrinx/test_code/fly_car.Decoration").(Decoration)
	//product.r1 = container.Resolve("github.com/lonegunmanb/syrinx/test_code/check_package_name_duplicate_a/model.Request").(*syrinx_p1.Request)
	//product.r2 = container.Resolve("github.com/lonegunmanb/syrinx/test_code/check_package_name_duplicate_b/model.Request").(*syrinx_p2.Request)
}
func Register_FlyCar(container ioc.Container) {
	container.RegisterFactory((*FlyCar)(nil), func(ioc ioc.Container) interface{} {
		return Create_FlyCar(ioc)
	})
}
