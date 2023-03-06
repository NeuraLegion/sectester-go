package bus

type Event struct {
	Message
}

func NewEvent(name string) (*Event, error) {
	message, err := NewMessage(name)
	if err != nil {
		return nil, err
	}

	return &Event{Message: message}, nil
}
