package console

// Console is a simple implementation of the Console interface for Unix systems.
type Console struct{}

// New returns a new instance of Console.
func New() *Console {
	return &Console{}
}

// IsColored always returns true for the Unix implementation of Console, which assumes support for ANSI escape codes.
func (c *Console) IsColored() bool {
	return true
}
