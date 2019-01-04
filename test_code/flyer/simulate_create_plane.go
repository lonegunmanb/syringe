package flyer

//func Create_Plane(container ioc.Container) *Plane {
//	plane := new(Plane)
//	Assemble_Plane(plane, container)
//	return plane
//}
//
//func Assemble_Plane(p *Plane, container ioc.Container) {
//	p.Wing = container.Resolve("github.com/lonegunmanb/syrinx/test_code/flyer.Wing").(Wing)
//}
//
//func Register_Plane(container ioc.Container) {
//	container.RegisterFactory((*Plane)(nil), func(ioc ioc.Container) interface{} {
//		return Create_Plane(ioc)
//	})
//}
