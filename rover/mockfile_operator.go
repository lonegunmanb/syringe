package rover

import (
	"github.com/stretchr/testify/mock"
	"io"
)

type mockFileOperator struct {
	mock.Mock
}

func (m *mockFileOperator) Open(path string) (io.Writer, error) {
	panic("implement me")
}

func (m *mockFileOperator) Del(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *mockFileOperator) FirstLine(path string) (string, error) {
	args := m.Called(path)
	return args.String(0), args.Error(1)
}
