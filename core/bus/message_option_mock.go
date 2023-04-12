// Code generated by mockery v2.20.0. DO NOT EDIT.

package bus

import mock "github.com/stretchr/testify/mock"

// MockMessageOption is an autogenerated mock type for the MessageOption type
type MockMessageOption struct {
	mock.Mock
}

type MockMessageOption_Expecter struct {
	mock *mock.Mock
}

func (_m *MockMessageOption) EXPECT() *MockMessageOption_Expecter {
	return &MockMessageOption_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: m
func (_m *MockMessageOption) Execute(m *Message) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(*Message) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMessageOption_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockMessageOption_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - m *Message
func (_e *MockMessageOption_Expecter) Execute(m interface{}) *MockMessageOption_Execute_Call {
	return &MockMessageOption_Execute_Call{Call: _e.mock.On("Execute", m)}
}

func (_c *MockMessageOption_Execute_Call) Run(run func(m *Message)) *MockMessageOption_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*Message))
	})
	return _c
}

func (_c *MockMessageOption_Execute_Call) Return(_a0 error) *MockMessageOption_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessageOption_Execute_Call) RunAndReturn(run func(*Message) error) *MockMessageOption_Execute_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockMessageOption interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockMessageOption creates a new instance of MockMessageOption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockMessageOption(t mockConstructorTestingTNewMockMessageOption) *MockMessageOption {
	mock := &MockMessageOption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}