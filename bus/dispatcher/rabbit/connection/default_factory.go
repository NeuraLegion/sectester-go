package connection

import (
	"github.com/rabbitmq/amqp091-go"
)

type DefaultFactory struct{}

func (d *DefaultFactory) Create(options *Options) (Connection, error) {
	auth := &amqp091.PlainAuth{Password: options.Password, Username: options.Username}

	//nolint:exhaustruct // redundant and optional fields
	return amqp091.DialConfig(options.Url, amqp091.Config{
		Heartbeat: options.HeartbeatInterval,
		SASL:      []amqp091.Authentication{auth},
	})
}
