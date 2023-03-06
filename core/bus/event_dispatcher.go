package bus

// EventDispatcher publishes events without waiting for a response.
// The ideal use case for the publish-subscribe model is
// when you want to simply notify another service that a certain condition has occurred.
type EventDispatcher interface {
	Publish(message *Message) error
}
