package codegen

type M struct {
	err error
}

func (m *M) Call(f func() error) *M {
	if m.err == nil {
		m.err = f()
	}
	return m
}

func Call(f func() error) *M {
	return &M{
		err: f(),
	}
}
