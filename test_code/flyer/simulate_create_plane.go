package flyer

import "github.com/lonegunmanb/syrinx/ioc"

func Assemble_plane(p *Plane, container ioc.Container) {
	p.Wing = container.Resolve("github.com/lonegunmanb/syrinx/test_code/flyer.Wing").(Wing)
}
