package fly_car

//func Create_FlyCar(container ioc.Container) *FlyCar {
//	product := new(FlyCar)
//	Assemble_FlyCar(product, container)
//	return product
//}
//func Assemble_FlyCar(product *FlyCar, container ioc.Container) {
//	product.Car = container.Resolve("github.com/lonegunmanb/syrinx/test_code/car.Car").(*car.Car)
//	product.Plane = *container.Resolve("github.com/lonegunmanb/syrinx/test_code/flyer.Plane").(*flyer.Plane)
//	product.Decoration = container.Resolve("github.com/lonegunmanb/syrinx/test_code/fly_car.Decoration").(Decoration)
//	product.R1 = container.Resolve("github.com/lonegunmanb/syrinx/test_code/check_package_name_duplicate_a/model.Request").(*p0.Request)
//	product.R2 = container.Resolve("github.com/lonegunmanb/syrinx/test_code/check_package_name_duplicate_b/model.Request").(*p1.Request)
//}
//func Register_FlyCar(container ioc.Container) {
//	container.RegisterFactory((*FlyCar)(nil), func(ioc ioc.Container) interface{} {
//		return Create_FlyCar(ioc)
//	})
//}
