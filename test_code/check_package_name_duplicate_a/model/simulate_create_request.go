package model

import (
	"github.com/lonegunmanb/syrinx/ioc"
)

func Create_Request(container ioc.Container) *Request {
	product := new(Request)
	Assemble_Request(product, container)
	return product
}

//noinspection ALL
func Assemble_Request(product *Request, container ioc.Container) {
}
func Register_Request(container ioc.Container) {
	container.RegisterFactory((*Request)(nil), func(ioc ioc.Container) interface{} {
		return Create_Request(ioc)
	})
}
