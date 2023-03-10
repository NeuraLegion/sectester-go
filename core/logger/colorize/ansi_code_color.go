package colorize

// AnsiCodeColor represents an ANSI escape code for a specific text color.
type AnsiCodeColor string

// Predefined ANSI escape codes for common text colors.
const (
	DefaultForeground AnsiCodeColor = "\x1B[39m\x1B[22m"
	Red               AnsiCodeColor = "\x1B[1m\x1B[31m"
	DarkRed           AnsiCodeColor = "\x1B[31m"
	Yellow            AnsiCodeColor = "\x1B[1m\x1B[33m"
	DarkGreen         AnsiCodeColor = "\x1B[32m"
	White             AnsiCodeColor = "\x1B[1m\x1B[37m"
	Cyan              AnsiCodeColor = "\x1B[1m\x1B[36m"
)

// String returns the string representation of the ANSI escape code.
func (c AnsiCodeColor) String() string {
	return string(c)
}
