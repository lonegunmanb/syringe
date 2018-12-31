package ioc_gen

import (
	"github.com/lonegunmanb/syrinx/ioc"
	"github.com/lonegunmanb/syrinx/test_code/car"
	"github.com/lonegunmanb/syrinx/test_code/fly_car"
	"github.com/lonegunmanb/syrinx/test_code/flyer"
)

func Resolve(key string, container ioc.Container) interface{} {
	switch key {
	case "github.com/lonegunmanb/syrinx/car.Car":
		{
			return car.Create_Car(container)
		}
	case "github.com/lonegunmanb/syrinx/flyer.Plane":
		{
			return flyer.Create_Plane(container)
		}
	case "github.com/lonegunmanb/syrinx/fly_car.FlyCar":
		{
			return fly_car.Create_fly_car_FlyCar(container)
		}
	default:
		{
			return nil
		}
	}
}
