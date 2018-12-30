package codegen_test

import (
	"errors"
	"github.com/lonegunmanb/syrinx/codegen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type invokeInterface interface {
	Invoke() error
}

type mockInvoke struct {
	mock.Mock
}

func (m *mockInvoke) Invoke() error {
	args := m.Called()
	return args.Error(0)
}

func TestCallEachTest(t *testing.T) {
	mockInvoke1 := new(mockInvoke)
	mockInvoke1.On("Invoke").Times(1).Return(nil)
	mockInvoke2 := new(mockInvoke)
	err := errors.New("expected")
	mockInvoke2.On("Invoke").Times(1).Return(err)
	mockInvoke3 := new(mockInvoke)
	mockInvoke2.On("Invoke").Times(0).Return(err)
	invokes := []invokeInterface{mockInvoke1, mockInvoke2, mockInvoke3}
	m := &codegen.M{
		Err: nil,
	}
	e := m.CallEach(invokes, func(i interface{}) error {
		return i.(invokeInterface).Invoke()
	}).Err
	assert.Equal(t, err, e)
	mockInvoke1.AssertExpectations(t)
	mockInvoke2.AssertExpectations(t)
	mockInvoke3.AssertExpectations(t)
}
