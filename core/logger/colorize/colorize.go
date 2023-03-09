package colorize

type Colorize interface {
	Colorize(color AnsiCodeColor, s string) string
}
