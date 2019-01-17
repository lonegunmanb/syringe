package util_test

import (
	"errors"
	. "github.com/lonegunmanb/syringe/util"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"testing"
)

type invokeInterface interface {
	Invoke() error
}

type mockInvoke struct {
	mock.Mock
	t *testing.T
}

func (m *mockInvoke) Invoke() error {
	args := m.Called()
	return args.Error(0)
}

func newMockInvoke(t *testing.T) *mockInvoke {
	return &mockInvoke{t: t}
}

func TestCallEachTest(t *testing.T) {
	Given("three func calls as slice which the second call will return error", t, func() {
		firstCall := newMockInvoke(t)
		firstCall.On("Invoke").Times(1).Return(nil)
		secondCall := newMockInvoke(t)
		err := errors.New("expected")
		secondCall.On("Invoke").Times(1).Return(err)
		lastCall := newMockInvoke(t)
		invokes := []invokeInterface{firstCall, secondCall, lastCall}
		When("use CallEach", func() {
			e := CallEach(invokes, func(i interface{}) error {
				return i.(invokeInterface).Invoke()
			}).Err
			Then("first two call occurred and the third never happened", func() {
				So(e, ShouldEqual, err)
				And(firstCall, shouldHappened)
				And(secondCall, shouldHappened)
				And(lastCall, neverHappened)
			})
		})
	})
}

func TestCallSingleRet(t *testing.T) {
	Given("two input arguments", t, func() {
		firstCallRet := "firstCallRet"
		secondCallRet := "secondCallRet"
		When("use CallSingleRet to call two functions which return them", func() {
			var argumentThatSecondCallReceived interface{}
			call := CallSingleRet(func() (interface{}, error) {
				return firstCallRet, nil
			}).CallSingleRet(func(input interface{}) (interface{}, error) {
				argumentThatSecondCallReceived = input
				return secondCallRet, nil
			})
			finalRet := call.Ret
			err := call.Err
			Then("call should run as firstCallRet", func() {
				So(err, ShouldBeNil)
				And(argumentThatSecondCallReceived, ShouldEqual, firstCallRet)
				And(finalRet, ShouldEqual, secondCallRet)
			})
		})
	})
}

func TestCallSingleRetOccurError(t *testing.T) {
	Given("two calls which first one will return error", t, func() {
		firstCall := newMockInvoke(t)
		expectedError := errors.New("expected")
		firstCall.On("Invoke").Times(1).Return(expectedError)
		lastCall := newMockInvoke(t)
		When("use CallSingleRet", func() {
			err := CallSingleRet(func() (interface{}, error) {
				return nil, firstCall.Invoke()
			}).CallSingleRet(func(input interface{}) (interface{}, error) {
				return nil, lastCall.Invoke()
			}).Err
			Then("second call should not happened", func() {
				So(err, ShouldEqual, expectedError)
				And(firstCall, shouldHappened)
				And(lastCall, neverHappened)
			})
		})
	})
}

//noinspection GoUnusedParameter
func neverHappened(actual interface{}, expected ...interface{}) string {
	mockInvoke := actual.(*mockInvoke)
	mockInvoke.AssertNotCalled(mockInvoke.t, "Invoke")
	return ""
}

//noinspection GoUnusedParameter
func shouldHappened(actual interface{}, expected ...interface{}) string {
	mockInvoke := actual.(*mockInvoke)
	mockInvoke.AssertExpectations(mockInvoke.t)
	return ""
}
