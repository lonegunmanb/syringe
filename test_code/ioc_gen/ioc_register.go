package ioc_gen

import (
	"github.com/lonegunmanb/syrinx/ioc"
	"github.com/lonegunmanb/syrinx/test_code/car"
	"github.com/lonegunmanb/syrinx/test_code/fly_car"
	"github.com/lonegunmanb/syrinx/test_code/flyer"
)

func CreateIoc() ioc.Container {
	container := ioc.NewContainer()
	car.Register_Car(container)
	fly_car.Register_FlyCar(container)
	flyer.Register_Plane(container)
	return container
}
