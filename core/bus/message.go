package bus

import (
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
)

// Message is used for syncing state between SDK, application and/or external services.
// This functionality is done by sending messages outside using a concrete implementation of Dispatcher.
//
// Depending on the type of derived class from the Message, it might be addressed to only one consumer
// or have typically multiple consumers as well.
// When a message is sent to multiple consumers, the appropriate event handler in each consumer handles the message.
//
// The Message is a data-holding class, but it implements a [Visitor pattern].
// to allow clients to perform operations on it using a visitor without modifying the source.
//
// [Visitor pattern]: https://en.wikipedia.org/wiki/Visitor_pattern#:~:text=In%20object%2Doriented%20programming%20and,structures%20without%20modifying%20the%20structures.
//
//nolint:lll // linter does not respect a long URL in the docs
type Message struct {
	correlationId string
	createdAt     time.Time
	name          string
	payload       any
	ttl           time.Duration
	expectReply   bool
}

type MessageOption func(m *Message) error

// WithPayload sets the payload to be sent to the application.
func WithPayload[T any](payload T) MessageOption {
	return func(m *Message) error {
		m.payload = payload
		return nil
	}
}

// WithTtl allows to change a default Message.Ttl.
func WithTtl(ttl time.Duration) MessageOption {
	return func(m *Message) error {
		if ttl <= time.Duration(0) {
			return errors.New("ttl must be greater than 0")
		}
		m.ttl = ttl
		return nil
	}
}

// WithExpectReply sets Message.ExpectReply to a desired value.
func WithExpectReply(expectReply bool) MessageOption {
	return func(m *Message) error {
		m.expectReply = expectReply
		return nil
	}
}

// NewMessage creates an instance of Message.
// Using this constructor, you can create an instance as follows:
//
//	var event, err = NewMessage("IssueDetected")
//
// For instance, you can dispatch a message in a way that is more approach you
// or convenient from the client's perspective.
//
//	// using a visitor pattern
//	event.Publish(dispatcher)
//	// or directly
//	dispatcher.Publish(event)
//
// The same is applicable for the request-response style.
// You just need to use the CommandDispatcher instead of EventDispatcher.
//
// To adjust its behavior you can use extra options, as shown below:
//
//	var command, err = NewMessage("Ping", WithTtl(time.Second * 60), WithExpectReply(false))
func NewMessage(name string, opts ...MessageOption) (*Message, error) {
	id, err := uuid.NewV4()

	if err != nil {
		return nil, err
	}

	m := &Message{
		correlationId: id.String(),
		createdAt:     time.Now(),
		name:          name,
		payload:       nil,
		ttl:           time.Second * 10,
		expectReply:   true,
	}

	for _, applyOpt := range opts {
		err = applyOpt(m)
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

// CorrelationId is used to ensure atomicity while working with EventBus. By default, random UUID.
func (m *Message) CorrelationId() string {
	return m.correlationId
}

// CreatedAt is exact date and time the event was created.
func (m *Message) CreatedAt() time.Time {
	return m.createdAt
}

// Name of a message.
func (m *Message) Name() string {
	return m.name
}

// Payload that we want to transmit to the remote service.
func (m *Message) Payload() any {
	return m.payload
}

// Ttl is a period of time that command should be handled before being discarded. By default, 10s.
func (m *Message) Ttl() time.Duration {
	return m.ttl
}

// ExpectReply indicates whether to wait for a reply. By default, `true`.
func (m *Message) ExpectReply() bool {
	return m.expectReply
}

// Publish publishes a message without waiting for a response.
// The ideal use case for the publish-subscribe model is when you want to simply notify another service
// that a certain condition has occurred.
func (m *Message) Publish(dispatcher EventDispatcher) error {
	return dispatcher.Publish(m)
}

// Execute executes a command according to the request-response message (aka 'Command') style is useful
// when you need to exchange messages between various external services.
// Using Execute you can easily ensure that the service has actually received the message and sent a response back.
func (m *Message) Execute(dispatcher CommandDispatcher) (any, error) {
	return dispatcher.Execute(m)
}
