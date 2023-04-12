// Code generated by mockery v2.20.0. DO NOT EDIT.

package logger

import mock "github.com/stretchr/testify/mock"

// MockLogger is an autogenerated mock type for the Logger type
type MockLogger struct {
	mock.Mock
}

type MockLogger_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLogger) EXPECT() *MockLogger_Expecter {
	return &MockLogger_Expecter{mock: &_m.Mock}
}

// Debug provides a mock function with given fields: message, args
func (_m *MockLogger) Debug(message string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// MockLogger_Debug_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Debug'
type MockLogger_Debug_Call struct {
	*mock.Call
}

// Debug is a helper method to define mock.On call
//   - message string
//   - args ...interface{}
func (_e *MockLogger_Expecter) Debug(message interface{}, args ...interface{}) *MockLogger_Debug_Call {
	return &MockLogger_Debug_Call{Call: _e.mock.On("Debug",
		append([]interface{}{message}, args...)...)}
}

func (_c *MockLogger_Debug_Call) Run(run func(message string, args ...interface{})) *MockLogger_Debug_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockLogger_Debug_Call) Return() *MockLogger_Debug_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockLogger_Debug_Call) RunAndReturn(run func(string, ...interface{})) *MockLogger_Debug_Call {
	_c.Call.Return(run)
	return _c
}

// Error provides a mock function with given fields: message, args
func (_m *MockLogger) Error(message string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// MockLogger_Error_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Error'
type MockLogger_Error_Call struct {
	*mock.Call
}

// Error is a helper method to define mock.On call
//   - message string
//   - args ...interface{}
func (_e *MockLogger_Expecter) Error(message interface{}, args ...interface{}) *MockLogger_Error_Call {
	return &MockLogger_Error_Call{Call: _e.mock.On("Error",
		append([]interface{}{message}, args...)...)}
}

func (_c *MockLogger_Error_Call) Run(run func(message string, args ...interface{})) *MockLogger_Error_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockLogger_Error_Call) Return() *MockLogger_Error_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockLogger_Error_Call) RunAndReturn(run func(string, ...interface{})) *MockLogger_Error_Call {
	_c.Call.Return(run)
	return _c
}

// Log provides a mock function with given fields: message, args
func (_m *MockLogger) Log(message string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// MockLogger_Log_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Log'
type MockLogger_Log_Call struct {
	*mock.Call
}

// Log is a helper method to define mock.On call
//   - message string
//   - args ...interface{}
func (_e *MockLogger_Expecter) Log(message interface{}, args ...interface{}) *MockLogger_Log_Call {
	return &MockLogger_Log_Call{Call: _e.mock.On("Log",
		append([]interface{}{message}, args...)...)}
}

func (_c *MockLogger_Log_Call) Run(run func(message string, args ...interface{})) *MockLogger_Log_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockLogger_Log_Call) Return() *MockLogger_Log_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockLogger_Log_Call) RunAndReturn(run func(string, ...interface{})) *MockLogger_Log_Call {
	_c.Call.Return(run)
	return _c
}

// LogLevel provides a mock function with given fields:
func (_m *MockLogger) LogLevel() LogLevel {
	ret := _m.Called()

	var r0 LogLevel
	if rf, ok := ret.Get(0).(func() LogLevel); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(LogLevel)
	}

	return r0
}

// MockLogger_LogLevel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LogLevel'
type MockLogger_LogLevel_Call struct {
	*mock.Call
}

// LogLevel is a helper method to define mock.On call
func (_e *MockLogger_Expecter) LogLevel() *MockLogger_LogLevel_Call {
	return &MockLogger_LogLevel_Call{Call: _e.mock.On("LogLevel")}
}

func (_c *MockLogger_LogLevel_Call) Run(run func()) *MockLogger_LogLevel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockLogger_LogLevel_Call) Return(_a0 LogLevel) *MockLogger_LogLevel_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLogger_LogLevel_Call) RunAndReturn(run func() LogLevel) *MockLogger_LogLevel_Call {
	_c.Call.Return(run)
	return _c
}

// Warn provides a mock function with given fields: message, args
func (_m *MockLogger) Warn(message string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// MockLogger_Warn_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Warn'
type MockLogger_Warn_Call struct {
	*mock.Call
}

// Warn is a helper method to define mock.On call
//   - message string
//   - args ...interface{}
func (_e *MockLogger_Expecter) Warn(message interface{}, args ...interface{}) *MockLogger_Warn_Call {
	return &MockLogger_Warn_Call{Call: _e.mock.On("Warn",
		append([]interface{}{message}, args...)...)}
}

func (_c *MockLogger_Warn_Call) Run(run func(message string, args ...interface{})) *MockLogger_Warn_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockLogger_Warn_Call) Return() *MockLogger_Warn_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockLogger_Warn_Call) RunAndReturn(run func(string, ...interface{})) *MockLogger_Warn_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockLogger interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockLogger creates a new instance of MockLogger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockLogger(t mockConstructorTestingTNewMockLogger) *MockLogger {
	mock := &MockLogger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}