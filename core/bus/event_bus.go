package bus

// EventBus allows to implement an event bus pattern.
type EventBus interface {
	EventDispatcher
	CommandDispatcher

	Register(name string, handler EventHandler) error
	Unregister(name string, handler EventHandler) error
}
