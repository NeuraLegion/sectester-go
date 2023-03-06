package bus

// EventBus allows to implement a event bus pattern.
type EventBus interface {
	EventDispatcher
	CommandDispatcher

	register(messageName string, handler EventHandler) error
	unregister(messageName string, handler EventHandler) error
}
