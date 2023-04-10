package connection

import "github.com/rabbitmq/amqp091-go"

type Manager interface {
	IsConnected() bool
	Connect() error
	TryConnect()
	CreateChannel() (*amqp091.Channel, error)
}
