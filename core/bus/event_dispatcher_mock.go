// Code generated by mockery v2.20.0. DO NOT EDIT.

package bus

import mock "github.com/stretchr/testify/mock"

// MockEventDispatcher is an autogenerated mock type for the EventDispatcher type
type MockEventDispatcher struct {
	mock.Mock
}

type MockEventDispatcher_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEventDispatcher) EXPECT() *MockEventDispatcher_Expecter {
	return &MockEventDispatcher_Expecter{mock: &_m.Mock}
}

// Publish provides a mock function with given fields: message
func (_m *MockEventDispatcher) Publish(message *Message) error {
	ret := _m.Called(message)

	var r0 error
	if rf, ok := ret.Get(0).(func(*Message) error); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockEventDispatcher_Publish_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Publish'
type MockEventDispatcher_Publish_Call struct {
	*mock.Call
}

// Publish is a helper method to define mock.On call
//   - message *Message
func (_e *MockEventDispatcher_Expecter) Publish(message interface{}) *MockEventDispatcher_Publish_Call {
	return &MockEventDispatcher_Publish_Call{Call: _e.mock.On("Publish", message)}
}

func (_c *MockEventDispatcher_Publish_Call) Run(run func(message *Message)) *MockEventDispatcher_Publish_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*Message))
	})
	return _c
}

func (_c *MockEventDispatcher_Publish_Call) Return(_a0 error) *MockEventDispatcher_Publish_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEventDispatcher_Publish_Call) RunAndReturn(run func(*Message) error) *MockEventDispatcher_Publish_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockEventDispatcher interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockEventDispatcher creates a new instance of MockEventDispatcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockEventDispatcher(t mockConstructorTestingTNewMockEventDispatcher) *MockEventDispatcher {
	mock := &MockEventDispatcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}