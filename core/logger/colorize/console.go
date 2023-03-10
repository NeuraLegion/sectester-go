package colorize

// Console is an interface for detecting console capabilities.
type Console interface {
	// IsColored returns true if the console supports colored output, and false otherwise.
	IsColored() bool
}
