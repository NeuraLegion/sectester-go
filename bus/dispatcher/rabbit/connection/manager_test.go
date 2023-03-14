package connection

import (
	"errors"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/NeuraLegion/sectester-go/core/logger"
)

type mockConnection struct {
	mock.Mock
}

func (m *mockConnection) NotifyClose(ch chan *amqp091.Error) chan *amqp091.Error {
	return ch
}

func (m *mockConnection) NotifyBlocked(ch chan amqp091.Blocking) chan amqp091.Blocking {
	return ch
}

func (m *mockConnection) IsClosed() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *mockConnection) Channel() (*amqp091.Channel, error) {
	args := m.Called()
	res, _ := args.Get(0).(*amqp091.Channel)
	return res, args.Error(1)
}

type mockFactory struct {
	mock.Mock
}

func (m *mockFactory) Create(opts *Options) (Connection, error) {
	args := m.Called(opts)
	res, _ := args.Get(0).(Connection)
	return res, args.Error(1)
}

type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) LogLevel() logger.LogLevel {
	args := m.Called()
	return args.Get(0).(logger.LogLevel)
}

func (m *mockLogger) Error(string, ...any) {
	// noop
}

func (m *mockLogger) Warn(string, ...any) {
	// noop
}

func (m *mockLogger) Log(string, ...any) {
	// noop
}

func (m *mockLogger) Debug(string, ...any) {
	// noop
}

var options = Options{
	Url:               "amqp://localhost:5672",
	Username:          "user",
	Password:          "pa$$word",
	HeartbeatInterval: time.Duration(1),
}

func setupManager() (*mockLogger, *mockFactory, *mockConnection, *DefaultManager) {
	mockedLogger := new(mockLogger)
	mockedFactory := new(mockFactory)
	mockedConnection := new(mockConnection)
	mgr := NewManager(
		options,
		mockedLogger,
		mockedFactory,
	)
	return mockedLogger, mockedFactory, mockedConnection, mgr
}

func TestDefaultManager_IsConnected_ReturnsFalse(t *testing.T) {
	// arrange
	_, f, _, mgr := setupManager()
	// act
	got := mgr.IsConnected()
	// assert
	assert.False(t, got)
	f.AssertExpectations(t)
}

func TestDefaultManager_Connect(t *testing.T) {
	// arrange
	_, f, c, mgr := setupManager()
	f.On("Create", &options).Return(c, nil)
	// act
	err := mgr.Connect()

	// Assert
	assert.NoError(t, err)
	f.AssertExpectations(t)
}

func TestDefaultManager_Connect_ReturnsError(t *testing.T) {
	// arrange
	_, f, _, mgr := setupManager()
	f.On("Create", &options).Return(nil, errors.New("something went wrong"))
	// act
	err := mgr.Connect()

	// Assert
	assert.ErrorContains(t, err, "something went wrong")
	f.AssertExpectations(t)
}

func TestDefaultManager_Connect_ConnectionClosed_Reconnects(t *testing.T) {
	// arrange
	_, f, c, mgr := setupManager()
	ch1 := make(chan *amqp091.Error, 1)
	ch2 := make(chan amqp091.Blocking, 1)
	f.On("Create", &options).Return(c, nil)
	go func() {
		ch1 <- &amqp091.Error{
			Code:    -1,
			Reason:  "something went wrong",
			Server:  false,
			Recover: false,
		}
	}()
	// act
	mgr.reconnecting(ch1, ch2)
	// assert
	f.AssertExpectations(t)
	f.AssertNumberOfCalls(t, "Create", 1)
}

func TestDefaultManager_Connect_ConnectionBlocked_Reconnects(t *testing.T) {
	// arrange
	_, f, c, mgr := setupManager()
	ch1 := make(chan *amqp091.Error, 1)
	ch2 := make(chan amqp091.Blocking, 1)
	f.On("Create", &options).Return(c, nil)
	go func() {
		ch2 <- amqp091.Blocking{
			Active: false,
			Reason: "something went wrong",
		}
	}()
	// act
	mgr.reconnecting(ch1, ch2)
	// assert
	f.AssertExpectations(t)
	f.AssertNumberOfCalls(t, "Create", 1)
}

func TestDefaultManager_TryConnect(t *testing.T) {
	// arrange
	_, f, c, mgr := setupManager()
	f.On("Create", &options).Return(c, nil)
	// act
	mgr.TryConnect()
	// Assert
	f.AssertExpectations(t)
	f.AssertNumberOfCalls(t, "Create", 1)
}

func TestDefaultManager_TryConnect_Connected(t *testing.T) {
	// arrange
	_, f, c, mgr := setupManager()
	c.On("IsClosed").Return(false)
	f.On("Create", &options).Return(c, nil)
	mgr.TryConnect()
	// act
	mgr.TryConnect()
	// Assert
	f.AssertExpectations(t)
	f.AssertNumberOfCalls(t, "Create", 1)
}

func TestDefaultManager_TryConnect_ConnectionFailed(t *testing.T) {
	// arrange
	_, f, c, mgr := setupManager()
	c.On("IsClosed").Return(false)
	f.On("Create", &options).Return(c, nil).Return(nil, errors.New("something went wrong"))
	// act
	mgr.TryConnect()
	// Assert
	f.AssertExpectations(t)
	f.AssertNumberOfCalls(t, "Create", 1)
}

func TestDefaultManager_CreateChannel_NotConnected_ReturnsError(t *testing.T) {
	// arrange
	l := new(mockLogger)
	f := new(mockFactory)
	c := new(mockConnection)
	c.On("IsClosed").Return(false)
	mgr := NewManager(
		options,
		l,
		f,
	)
	// act
	channel, err := mgr.CreateChannel()
	// Assert
	assert.Error(t, err)
	assert.Nil(t, channel)
	f.AssertExpectations(t)
}

func TestDefaultManager_CreateChannel_FailedToOpenChannel_ReturnsError(t *testing.T) {
	// arrange
	_, f, c, mgr := setupManager()
	c.On("IsClosed").Return(false)
	c.On("Channel").Return(nil, errors.New("something went wrong"))
	f.On("Create", &options).Return(c, nil).Return(c, nil)
	mgr.TryConnect()
	// act
	channel, err := mgr.CreateChannel()
	// Assert
	assert.ErrorContains(t, err, "something went wrong")
	assert.Nil(t, channel)
	f.AssertExpectations(t)
}

func TestDefaultManager_CreateChannel_Connected_CreatesChannel(t *testing.T) {
	// arrange
	_, f, c, mgr := setupManager()
	c.On("IsClosed").Return(false)
	c.On("Channel").Return(&amqp091.Channel{}, nil)
	f.On("Create", &options).Return(c, nil).Return(c, nil)
	mgr.TryConnect()
	// act
	channel, err := mgr.CreateChannel()
	// Assert
	assert.NoError(t, err)
	assert.IsType(t, new(amqp091.Channel), channel)
	f.AssertExpectations(t)
}
