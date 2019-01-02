package fly_car

import (
	"github.com/lonegunmanb/syrinx/test_code/car"
	model1 "github.com/lonegunmanb/syrinx/test_code/check_package_name_duplicate_a/model"
	model2 "github.com/lonegunmanb/syrinx/test_code/check_package_name_duplicate_b/model"
	"github.com/lonegunmanb/syrinx/test_code/flyer"
)

type FlyCar struct {
	*car.Car    `inject:""`
	flyer.Plane `inject:""`
	Decoration  Decoration `inject:""`
	R1          *model1.Request
	R2          *model2.Request
}

type Decoration interface {
	LookAndFeel() string
}

type FancyDecoration struct {
}

func (f *FancyDecoration) LookAndFeel() string {
	return "Fancy"
}
