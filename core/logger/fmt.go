package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/NeuraLegion/sectester-go/core/logger/clock"
	"github.com/NeuraLegion/sectester-go/core/logger/colorize"
)

type Fmt struct {
	logLevel LogLevel
	writer   io.Writer
	clock    clock.Provider
	mu       sync.Mutex
}

func Default(logLevel LogLevel) *Fmt {
	return New(logLevel, os.Stdout, &clock.SystemProvider{})
}

func New(logLevel LogLevel, writer io.Writer, clock clock.Provider) *Fmt {
	return &Fmt{logLevel: logLevel, writer: writer, clock: clock}
}

func (f *Fmt) LogLevel() LogLevel {
	return f.logLevel
}

func (f *Fmt) SetLogLevel(logLevel LogLevel) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logLevel = logLevel
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
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.logLevel < level {
		return
	}
	template := fmt.Sprintf("%s - %s", f.buildHeader(level), message)
	_, _ = fmt.Fprintf(f.writer, template, args...)
}

func (f *Fmt) buildHeader(level LogLevel) string {
	timestamp := time.Now().Format(time.RFC3339)
	header := fmt.Sprintf("[%s] [%s]", timestamp, level.Humanize())

	return f.getColor(level).String() +
		header +
		f.getColor(-1).String()
}

func (f *Fmt) getColor(level LogLevel) colorize.AnsiCodeColor {
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
