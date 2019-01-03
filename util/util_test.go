package util_test

import (
	"errors"
	"github.com/lonegunmanb/syrinx/util"
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
	//mockInvoke3.On("Invoke").Times(0).Return(err)
	invokes := []invokeInterface{mockInvoke1, mockInvoke2, mockInvoke3}

	e := util.CallEach(invokes, func(i interface{}) error {
		return i.(invokeInterface).Invoke()
	}).Err
	assert.Equal(t, err, e)
	mockInvoke1.AssertExpectations(t)
	mockInvoke2.AssertExpectations(t)
	mockInvoke3.AssertNotCalled(t, "Invoke")
}

func TestCallSingleRet(t *testing.T) {
	expected := "expected"
	err := util.CallSingleRet(func() (interface{}, error) {
		return expected, nil
	}).CallSingleRet(func(input interface{}) (interface{}, error) {
		assert.Equal(t, expected, input)
		return input, nil
	}).Err
	assert.Nil(t, err)
}

func TestCallSingleRetOccurError(t *testing.T) {
	mockInvoke1 := new(mockInvoke)
	expectedError := errors.New("expected")
	mockInvoke1.On("Invoke").Times(1).Return(expectedError)
	mockInvoke2 := new(mockInvoke)
	err := util.CallSingleRet(func() (interface{}, error) {
		return nil, mockInvoke1.Invoke()
	}).CallSingleRet(func(input interface{}) (interface{}, error) {
		return nil, mockInvoke2.Invoke()
	}).Err
	assert.Equal(t, expectedError, err)
	mockInvoke2.AssertNotCalled(t, "Invoke")
}
