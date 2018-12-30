package codegen

import "github.com/ahmetb/go-linq"

type M struct {
	Err error
}

func (m *M) Call(f func() error) *M {
	if m.Err == nil {
		m.Err = f()
	}
	return m
}

func Call(f func() error) *M {
	return &M{
		Err: f(),
	}
}

func (m *M) CallEach(iterable interface{}, f func(item interface{}) error) *M {
	if m.Err == nil {
		linq.From(iterable).TakeWhile(func(item interface{}) bool {
			err := f(item)
			if err != nil {
				m.Err = err
				return false
			}
			return true
		}).Last()
	}
	return m
}
