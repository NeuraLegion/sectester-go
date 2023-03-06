package bus

// EventHandler allows to define an event handler.
// You should implement the EventHandler interface and
// use the IoC container to register a handler using the interface as a provider:
//
//	type IssueDetectedHandler struct {}
//
//	func (h *IssueDetectedHandler) Handle(message Message) (any, error) {
//		// your implementation
//		return nil, nil
//	}
type EventHandler interface {
	Handle(message Message) (any, error)
}
