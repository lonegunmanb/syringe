package car

//import (
//	"github.com/lonegunmanb/syrinx/ioc"
//	"github.com/lonegunmanb/syrinx/test_code/engine"
//)
//
////noinspection GoUnusedExportedFunction
//func Create_Car(container ioc.Container) *Car {
//	r := new(Car)
//	Assemble_Car(r, container)
//	return r
//}
//
//func Assemble_Car(c *Car, container ioc.Container) {
//	c.Engine = container.Resolve("github.com/lonegunmanb/syrinx/test_code/engine.Engine").(engine.Engine)
//}
//
//func Register_Car(container ioc.Container) {
//	container.RegisterFactory((*Car)(nil), func(ioc ioc.Container) interface{} {
//		return Create_Car(ioc)
//	})
//}
