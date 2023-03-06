package bus

type CommandDispatcher interface {
	Execute(message *Message) (any, error)
}
