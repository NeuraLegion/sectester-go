package colorize

// Colorize defines an interface for applying color to a string.
type Colorize interface {
	// Colorize applies the specified ANSI color code to the provided string.
	// The returned string is the original string with the ANSI escape sequence
	// for the specified color code prepended and the sequence for resetting
	// the color code appended.
	Colorize(color AnsiCodeColor, s string) string
}
