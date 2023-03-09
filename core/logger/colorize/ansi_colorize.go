package colorize

type AnsiColorize struct {
	enabled bool
}

func New(console Console) *AnsiColorize {
	return &AnsiColorize{
		enabled: console.IsColored(),
	}
}

func (a *AnsiColorize) Colorize(color AnsiCodeColor, s string) string {
	if !a.enabled {
		return s
	} else {
		return color.String() + s + DefaultForeground.String()
	}
}
