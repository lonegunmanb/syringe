package util

import "github.com/ahmetb/go-linq"

type Action struct {
	Err error
}

func (m *Action) Call(f func() error) *Action {
	if m.Err == nil {
		m.Err = f()
	}
	return m
}

func Call(f func() error) *Action {
	return &Action{
		Err: f(),
	}
}

func (m *Action) CallSingleRet(f func() (interface{}, error)) *SingleFunc {
	if m.Err == nil {
		return CallSingleRet(f)
	}
	return &SingleFunc{
		Err: m.Err,
	}
}

func CallSingleRet(f func() (interface{}, error)) *SingleFunc {
	ret, err := f()
	return &SingleFunc{
		Ret: ret,
		Err: err,
	}
}

func CallEach(iterable interface{}, f func(item interface{}) error) *Action {
	m := &Action{}
	return m.CallEach(iterable, f)
}

func (m *Action) CallEach(iterable interface{}, f func(item interface{}) error) *Action {
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

type SingleFunc struct {
	Err error
	Ret interface{}
}

func (s *SingleFunc) CallSingleRet(f func(interface{}) (interface{}, error)) *SingleFunc {
	if s.Err == nil {
		s.Ret, s.Err = f(s.Ret)
	}
	return s
}

func (s *SingleFunc) Call(f func(interface{}) error) *Action {
	var err error
	if s.Err == nil {
		err = f(s.Ret)
	} else {
		err = s.Err
	}
	return &Action{
		Err: err,
	}
}
