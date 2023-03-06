package bus

import (
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Message struct {
	correlationId string
	createdAt     time.Time
	name          string
	payload       any
	ttl           time.Duration
	expectReply   bool
}

type MessageOption func(m *Message) error

func WithPayload[T any](payload T) MessageOption {
	return func(m *Message) error {
		m.payload = payload
		return nil
	}
}

func WithTtl(ttl time.Duration) MessageOption {
	return func(m *Message) error {
		if ttl <= time.Duration(0) {
			return errors.New("ttl must be greater than 0")
		}
		m.ttl = ttl
		return nil
	}
}

func WithExpectReply(expectReply bool) MessageOption {
	return func(m *Message) error {
		m.expectReply = expectReply
		return nil
	}
}

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
		ttl:           time.Minute,
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

func (m *Message) CorrelationId() string {
	return m.correlationId
}

func (m *Message) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Message) Name() string {
	return m.name
}

func (m *Message) Payload() any {
	return m.payload
}

func (m *Message) Ttl() time.Duration {
	return m.ttl
}

func (m *Message) ExpectReply() bool {
	return m.expectReply
}

func (m *Message) Publish(dispatcher EventDispatcher) error {
	return dispatcher.Publish(m)
}

func (m *Message) Execute(dispatcher CommandDispatcher) (any, error) {
	return dispatcher.Execute(m)
}
