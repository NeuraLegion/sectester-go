package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NeuraLegion/sectester-go/bus/dispatcher/rabbit/connection"
	"github.com/NeuraLegion/sectester-go/core/bus"
	"github.com/NeuraLegion/sectester-go/core/logger"
	"github.com/rabbitmq/amqp091-go"
)

type consumedMessage struct {
	correlationId string
	name          string
	payload       any
	replyTo       string
	timestamp     time.Time
}

type replyMessage struct {
	routingKey    string
	exchange      string
	correlationId string
	name          string
	payload       any
	replyTo       string
	timestamp     time.Time
}

type Rabbit struct {
	options  *Options
	manager  connection.Manager
	logger   logger.Logger
	handlers map[string][]bus.EventHandler
	channel  *amqp091.Channel
}

func NewRabbit(options *Options, manager connection.Manager, logger logger.Logger) (*Rabbit, error) {
	eb := &Rabbit{
		options:  options,
		manager:  manager,
		logger:   logger,
		handlers: map[string][]bus.EventHandler{},
		channel:  nil,
	}
	channel, err := eb.createConsumerChannel()
	if err != nil {
		return nil, err
	}
	eb.channel = channel
	return eb, nil
}

func (r *Rabbit) Execute(message *bus.Message) (any, error) {
	// TODO implement me
	panic("implement me")
}

func (r *Rabbit) Publish(message *bus.Message) error {
	r.manager.TryConnect()

	return r.sendMessageViaNewChannel(&replyMessage{
		payload:       message.Payload(),
		name:          message.Name(),
		routingKey:    message.Name(),
		exchange:      r.options.Exchange,
		correlationId: message.CorrelationId(),
		timestamp:     message.CreatedAt(),
		replyTo:       "",
	})
}

func (r *Rabbit) Register(name string, handler bus.EventHandler) error {
	handlers := r.handlers[name]
	if handlers == nil {
		handlers = []bus.EventHandler{}
		if err := r.bindQueue(name); err != nil {
			return err
		}
	}
	r.handlers[name] = append(handlers, handler)
	return nil
}

func (r *Rabbit) Unregister(name string, handler bus.EventHandler) error {
	handlers := r.handlers[name]
	if handlers == nil {
		return fmt.Errorf(
			"no subscriptions found. Please register a handler for the %s event in the event bus",
			name,
		)
	}

	idx := r.findHandler(handlers, handler)
	r.handlers[name] = append(handlers[:idx], handlers[idx+1:]...)

	if len(r.handlers) == 0 {
		if err := r.unBindQueue(name); err != nil {
			return err
		}
	}

	return nil
}

func (r *Rabbit) findHandler(collection []bus.EventHandler, el bus.EventHandler) int {
	for i, x := range collection {
		if x == el {
			return i
		}
	}
	return -1
}

func (r *Rabbit) bindQueue(name string) error {
	r.manager.TryConnect()
	channel, err := r.manager.CreateChannel()
	if err != nil {
		return err
	}
	return channel.QueueBind(r.options.ClientQueue,
		name,
		r.options.Exchange,
		false,
		nil)
}

func (r *Rabbit) unBindQueue(name string) error {
	r.manager.TryConnect()
	channel, err := r.manager.CreateChannel()
	if err != nil {
		return err
	}
	return channel.QueueUnbind(r.options.ClientQueue,
		name,
		r.options.Exchange,
		nil)
}

func (r *Rabbit) createConsumerChannel() (*amqp091.Channel, error) {
	r.manager.TryConnect()
	channel, err := r.manager.CreateChannel()
	if err != nil {
		return nil, err
	}
	go func(ch chan *amqp091.Error) {
		<-ch
		channel, _ = r.createConsumerChannel()
		r.channel = channel
	}(channel.NotifyClose(make(chan *amqp091.Error, 1)))
	err = r.bindQueueToExchange(channel)
	if err != nil {
		return nil, err
	}
	err = r.startBasicConsume(channel)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (r *Rabbit) bindQueueToExchange(channel *amqp091.Channel) error {
	err := channel.ExchangeDeclare(
		r.options.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	_, err = channel.QueueDeclare(
		r.options.ClientQueue,
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = channel.Qos(r.options.PrefetchCount, 0, false)
	if err != nil {
		return err
	}
	return nil
}

func (r *Rabbit) startBasicConsume(channel *amqp091.Channel) error {
	consume, err := channel.Consume(
		r.options.ClientQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	go func(consume <-chan amqp091.Delivery) {
		for message := range consume {
			r.receivingMessage(&message)
		}
	}(consume)
	return nil
}

func (r *Rabbit) receivingMessage(message *amqp091.Delivery) {
	if message.Redelivered {
		return
	}
	name := getMessageName(message)
	body, err := r.deserialize(message)
	if err != nil {
		return
	}
	cm := r.buildConsumedMessage(message, name, body)
	r.logger.Debug(
		"Received a event (%s) with following payload: %j", cm.name,
		cm.payload,
	)
	handlers, err := r.getHandlers(cm.name)
	if err != nil {
		return
	}
	for _, handler := range handlers {
		if err = r.handleEvent(cm, handler); err != nil {
			r.logger.Error(
				"Error while processing a message (%s) due to error occurred. Event: %s",
				cm.correlationId,
				cm.payload,
			)
		}
	}
}

func (r *Rabbit) buildConsumedMessage(message *amqp091.Delivery, name string, body any) *consumedMessage {
	return &consumedMessage{
		correlationId: message.CorrelationId,
		name:          name,
		payload:       body,
		replyTo:       message.ReplyTo,
		timestamp:     message.Timestamp,
	}
}

func (r *Rabbit) deserialize(message *amqp091.Delivery) (map[string]any, error) {
	var body map[string]any
	err := json.Unmarshal(message.Body, &body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (r *Rabbit) handleEvent(cm *consumedMessage, handler bus.EventHandler) error {
	m := bus.NewRawMessage(cm.name, cm.correlationId, cm.timestamp, cm.payload)
	result, err := handler.Handle(*m)
	if err != nil {
		return err
	}
	if result != nil && len(cm.replyTo) != 0 {
		if err = r.sendReplyOnEvent(cm, result); err != nil {
			return err
		}
	}
	return nil
}

func (r *Rabbit) sendReplyOnEvent(cm *consumedMessage, result any) error {
	r.logger.Debug(
		"Sending a reply (%s) back with following payload: %j",
		cm.name,
		result,
	)
	//nolint:exhaustruct // redundant and optional fields
	return r.sendMessageViaNewChannel(&replyMessage{
		routingKey:    cm.replyTo,
		correlationId: cm.correlationId,
		payload:       result,
		timestamp:     time.Now(),
	})
}

func (r *Rabbit) sendMessageViaNewChannel(rm *replyMessage) error {
	channel, err := r.manager.CreateChannel()
	if err != nil {
		return err
	}
	return r.sendMessage(channel, rm)
}

func (r *Rabbit) sendMessage(channel *amqp091.Channel, rm *replyMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r.logger.Debug("Send a message with following parameters: %j", rm)
	body, err := json.Marshal(rm.payload)
	if err != nil {
		return err
	}
	return channel.PublishWithContext(
		ctx,
		rm.exchange,
		rm.routingKey,
		true,
		false,
		//nolint:exhaustruct // redundant and optional fields
		amqp091.Publishing{
			CorrelationId: rm.correlationId,
			Type:          rm.name,
			ReplyTo:       rm.replyTo,
			Timestamp:     rm.timestamp,
			DeliveryMode:  2,
			ContentType:   "application/json",
			Body:          body,
		},
	)
}

func (r *Rabbit) getHandlers(name string) ([]bus.EventHandler, error) {
	handlers := r.handlers[name]
	if handlers == nil {
		return nil, fmt.Errorf("no subscriptions found. Please register a handler for the %s event in the event bus", name)
	}
	if len(handlers) == 0 {
		return nil, fmt.Errorf("event handler not found. Please register a handler for the following events: %s", name)
	}
	return handlers, nil
}

func getMessageName(message *amqp091.Delivery) string {
	if len(message.Type) == 0 {
		return message.RoutingKey
	}
	return message.Type
}
