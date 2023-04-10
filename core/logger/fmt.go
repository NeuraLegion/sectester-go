package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/NeuraLegion/sectester-go/core/logger/clock"
	"github.com/NeuraLegion/sectester-go/core/logger/colorize"
	"github.com/NeuraLegion/sectester-go/core/logger/console"
)

// A Fmt represents an active logging object that generates lines of
// output to an os.Stdout by default. A Fmt can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the io.Writer.
type Fmt struct {
	logLevel  LogLevel
	writer    io.Writer
	clock     clock.Provider
	mu        sync.Mutex
	colorizer colorize.Colorize
}

// Default creates a default instance of Fmt that writes output to os.Stdout.
func Default(logLevel LogLevel) *Fmt {
	return New(logLevel, os.Stdout, &clock.SystemProvider{})
}

// New allows to create an instance of Fmt customizing a writer and time provider.
//
//	var buf bytes.Buffer
//	logger.New(logger.Error, &buf, clock.SystemProvider{})
func New(logLevel LogLevel, writer io.Writer, clock clock.Provider) *Fmt {
	return &Fmt{logLevel: logLevel, writer: writer, clock: clock, mu: sync.Mutex{}, colorizer: colorize.New(console.New())}
}

// LogLevel returns a current log level.
func (f *Fmt) LogLevel() LogLevel {
	return f.logLevel
}

// SetLogLevel sets the log level to a desired value.
func (f *Fmt) SetLogLevel(logLevel LogLevel) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.logLevel = logLevel
}

// Error prints a formatted message to the output.
// Arguments are handled in the manner of fmt.Print.
func (f *Fmt) Error(message string, args ...any) {
	f.write(Error, message, args...)
}

// Warn prints a formatted message to the output.
// Arguments are handled in the manner of fmt.Print.
func (f *Fmt) Warn(message string, args ...any) {
	f.write(Warn, message, args...)
}

// Log prints a formatted message to the output.
// Arguments are handled in the manner of fmt.Print.
func (f *Fmt) Log(message string, args ...any) {
	f.write(Notice, message, args...)
}

// Debug prints a formatted message to the output.
// Arguments are handled in the manner of fmt.Print.
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

	return f.colorizer.Colorize(f.getColor(level), header)
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
