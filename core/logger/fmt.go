package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/NeuraLegion/sectester-go/core/logger/colorize"
)

type Fmt struct {
	LogLevel LogLevel
}

func (f *Fmt) Error(message string, args ...any) {
	f.write(Error, message, args...)
}

func (f *Fmt) Warn(message string, args ...any) {
	f.write(Warn, message, args...)
}

func (f *Fmt) Log(message string, args ...any) {
	f.write(Notice, message, args...)
}

func (f *Fmt) Debug(message string, args ...any) {
	f.write(Verbose, message, args...)
}

func (f *Fmt) write(level LogLevel, message string, args ...any) {
	if f.LogLevel < level {
		return
	}
	template := fmt.Sprintf("%s - %s", f.buildHeader(level), message)
	var output *os.File
	if f.LogLevel <= Warn {
		output = os.Stderr
	} else {
		output = os.Stdout
	}
	_, _ = fmt.Fprintf(output, template, args...)
}

func (f *Fmt) buildHeader(level LogLevel) string {
	timestamp := time.Now().Format(time.RFC3339)
	header := fmt.Sprintf("[%s] [%s]", timestamp, level.Humanize())

	return f.getForegroundColorAnsiCode(level).String() +
		header +
		f.getForegroundColorAnsiCode(-1).String()
}

func (f *Fmt) getForegroundColorAnsiCode(level LogLevel) colorize.AnsiCodeColor {
	switch level {
	case Error:
		return colorize.Red
	case Warn:
		return colorize.Yellow
	case Notice:
		return colorize.DarkGreen
	case Verbose:
		return colorize.Cyan
	case Silent:
		return colorize.White
	}

	return colorize.DefaultForeground
}
