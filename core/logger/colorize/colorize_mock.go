// Code generated by mockery v2.20.0. DO NOT EDIT.

package colorize

import mock "github.com/stretchr/testify/mock"

// MockColorize is an autogenerated mock type for the Colorize type
type MockColorize struct {
	mock.Mock
}

type MockColorize_Expecter struct {
	mock *mock.Mock
}

func (_m *MockColorize) EXPECT() *MockColorize_Expecter {
	return &MockColorize_Expecter{mock: &_m.Mock}
}

// Colorize provides a mock function with given fields: color, s
func (_m *MockColorize) Colorize(color AnsiCodeColor, s string) string {
	ret := _m.Called(color, s)

	var r0 string
	if rf, ok := ret.Get(0).(func(AnsiCodeColor, string) string); ok {
		r0 = rf(color, s)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockColorize_Colorize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Colorize'
type MockColorize_Colorize_Call struct {
	*mock.Call
}

// Colorize is a helper method to define mock.On call
//   - color AnsiCodeColor
//   - s string
func (_e *MockColorize_Expecter) Colorize(color interface{}, s interface{}) *MockColorize_Colorize_Call {
	return &MockColorize_Colorize_Call{Call: _e.mock.On("Colorize", color, s)}
}

func (_c *MockColorize_Colorize_Call) Run(run func(color AnsiCodeColor, s string)) *MockColorize_Colorize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(AnsiCodeColor), args[1].(string))
	})
	return _c
}

func (_c *MockColorize_Colorize_Call) Return(_a0 string) *MockColorize_Colorize_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockColorize_Colorize_Call) RunAndReturn(run func(AnsiCodeColor, string) string) *MockColorize_Colorize_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockColorize interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockColorize creates a new instance of MockColorize. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockColorize(t mockConstructorTestingTNewMockColorize) *MockColorize {
	mock := &MockColorize{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
