package console

import (
	"os"

	"golang.org/x/sys/windows"
)

type Console struct {
	isColored bool
}

const ansiColorRequiredMode uint = windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING

func New() *Console {
	c := &Console{}
	c.isColored = len(os.Getenv("NO_COLOR")) == 0 && c.enable()
	return c
}

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
