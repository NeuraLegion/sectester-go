package bus

type EventDispatcher interface {
	Publish(message *Message) error
}
