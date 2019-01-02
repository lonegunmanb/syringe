package ioc_gen

import (
	"github.com/lonegunmanb/syrinx/ioc"
	"github.com/lonegunmanb/syrinx/test_code/car"
	p0 "github.com/lonegunmanb/syrinx/test_code/check_package_name_duplicate_a/model"
	p1 "github.com/lonegunmanb/syrinx/test_code/check_package_name_duplicate_b/model"
	"github.com/lonegunmanb/syrinx/test_code/fly_car"
	"github.com/lonegunmanb/syrinx/test_code/flyer"
)

func CreateIoc() ioc.Container {
	container := ioc.NewContainer()
	car.Register_Car(container)
	fly_car.Register_FlyCar(container)
	flyer.Register_Plane(container)
	p0.Register(container)
	p1.Register(container)
	return container
}
