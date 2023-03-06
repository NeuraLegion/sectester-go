package bus

// CommandDispatcher allows to implement the request-response message style is useful
// when you need to exchange messages between various external services.
// Using CommandDispatcher you can easily ensure
// that the service has actually received the message and sent a response back.
type CommandDispatcher interface {
	Execute(message *Message) (any, error)
}
