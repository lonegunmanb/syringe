package codegen

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
