package bus

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Message interface {
	CorrelationId() string
	CreatedAt() time.Time
	Name() string
}

type BaseMessage struct {
	correlationId string
	createdAt     time.Time
	name          string
}

func NewMessage(name string) (*BaseMessage, error) {
	id, err := uuid.NewV4()

	if err != nil {
		return nil, err
	}

	return &BaseMessage{name: name, correlationId: id.String(), createdAt: time.Now()}, nil
}

func (b *BaseMessage) CorrelationId() string {
	return b.correlationId
}

func (b *BaseMessage) CreatedAt() time.Time {
	return b.createdAt
}

func (b *BaseMessage) Name() string {
	return b.name
}
