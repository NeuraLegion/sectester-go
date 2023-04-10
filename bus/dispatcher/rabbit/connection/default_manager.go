package connection

import (
	"errors"
	"sync"

	"github.com/NeuraLegion/sectester-go/core/logger"
	"github.com/rabbitmq/amqp091-go"
)

type DefaultManager struct {
	mu                sync.Mutex
	connection        Connection
	options           Options
	logger            logger.Logger
	connectionFactory Factory
}

func NewManager(options Options, logger logger.Logger, connectionFactory Factory) *DefaultManager {
	return &DefaultManager{
		sync.Mutex{},
		nil,
		options,
		logger,
		connectionFactory,
	}
}

func (a *DefaultManager) IsConnected() bool {
	return a.connection != nil && !a.connection.IsClosed()
}

func (a *DefaultManager) Connect() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	connection, err := a.connectionFactory.Create(&a.options)
	if err != nil {
		return err
	}
	a.setConnection(connection)
	a.logger.Debug("Event bus connected to %s", a.options.Url)
	return nil
}

func (a *DefaultManager) TryConnect() {
	if a.IsConnected() {
		return
	}

	err := a.Connect()

	if err != nil {
		return
	}
}

func (a *DefaultManager) CreateChannel() (*amqp091.Channel, error) {
	if !a.IsConnected() {
		return nil, errors.New("please make sure that client established a connection with host")
	}
	channel, err := a.connection.Channel()
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (a *DefaultManager) setConnection(connection Connection) {
	a.connection = connection

	closing := make(chan *amqp091.Error, 1)
	blocking := make(chan amqp091.Blocking, 1)

	connection.NotifyClose(closing)
	connection.NotifyBlocked(blocking)

	go a.reconnecting(closing, blocking)
}

func (a *DefaultManager) reconnecting(closing chan *amqp091.Error, blocking chan amqp091.Blocking) {
	defer close(closing)
	defer close(blocking)

	select {
	case <-closing:
		a.logger.Warn("A Event Bus connection shutdown. Trying to re-connect.")
		_ = a.Connect()
	case <-blocking:
		a.logger.Warn("A Event Bus connection blocked. Trying to re-connect.")
		_ = a.Connect()
	}
}
