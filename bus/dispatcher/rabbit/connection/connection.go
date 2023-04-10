package connection

import "github.com/rabbitmq/amqp091-go"

type Connection interface {
	NotifyClose(ch chan *amqp091.Error) chan *amqp091.Error
	NotifyBlocked(ch chan amqp091.Blocking) chan amqp091.Blocking
	IsClosed() bool
	Channel() (*amqp091.Channel, error)
}
