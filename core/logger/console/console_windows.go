package console

import (
	"os"

	"golang.org/x/sys/windows"
)

// Console is an implementation of the Console interface for Windows systems.
type Console struct {
	isColored bool
}

const ansiColorRequiredMode uint = windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING

// New returns a new instance of Console trying to enable Windows ANSI colors for the stdout and stderr console handles.
func New() *Console {
	c := &Console{}
	c.isColored = c.enable()
	return c
}

// IsColored returns true if the console supports colored output, and false otherwise.
func (c *Console) IsColored() bool {
	return c.isColored
}

func (a *Console) enable() bool {
	return EnableWindowsAnsiColors(os.Stdout.Fd()) && EnableWindowsAnsiColors(os.Stderr.Fd())
}

func (a *Console) enableForHandle(fd uintptr) bool {
	handle := windows.Handle(fd)
	var outMode uint32
	err := windows.GetConsoleMode(handle, &outMode)
	if err != nil {
		return false
	}
	if outMode&ansiColorRequiredMode == ansiColorRequiredMode {
		return true
	}
	return !windows.SetConsoleMode(handle, outMode|ansiColorRequiredMode)
}
