package colorize

import "os"

// AnsiColorize is a utility for colorizing strings using ANSI escape codes.
type AnsiColorize struct {
	enabled bool
}

// New returns a new instance of AnsiColorize that is enabled if the console
// supports colored output and the "NO_COLOR" environment variable is not set.
func New(console Console) *AnsiColorize {
	return &AnsiColorize{
		enabled: len(os.Getenv("NO_COLOR")) == 0 && console.IsColored(),
	}
}

// Colorize returns the input string wrapped with ANSI escape codes for the
// specified color if the colored output is enabled. Otherwise, the original
// string is returned.
func (a *AnsiColorize) Colorize(color AnsiCodeColor, s string) string {
	if !a.enabled {
		return s
	} else {
		return color.String() + s + DefaultForeground.String()
	}
}
