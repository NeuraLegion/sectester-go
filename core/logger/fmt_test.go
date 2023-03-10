package logger

import (
	"bytes"
	"fmt"
	"github.com/NeuraLegion/sectester-go/core/logger/colorize"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type mockProvider struct {
	mock.Mock
}

func (p *mockProvider) Now() time.Time {
	args := p.Called()

	return args.Get(0).(time.Time)
}

func createLogPattern(level string, timestamp time.Time, color colorize.AnsiCodeColor) string {
	return fmt.Sprintf(
		"%s[%s] [%s]%s - message",
		color,
		timestamp.Format(time.RFC3339),
		level,
		colorize.DefaultForeground,
	)
}

func TestFmt_Error(t *testing.T) {
	timestamp := time.Now()
	provider := new(mockProvider)
	provider.On("Now").Return(timestamp)
	expected := createLogPattern("ERROR  ", timestamp, colorize.Red)
	for _, data := range []struct {
		input    LogLevel
		expected string
	}{
		{Error, expected},
		{Warn, ""},
		{Notice, ""},
		{Verbose, ""},
		{Silent, ""},
	} {
		t.Run(data.input.String(), func(t *testing.T) {
			// arrange
			var buf bytes.Buffer
			logger := New(data.input, &buf, provider)

			// arc
			logger.Error("message")

			// assert
			assert.Equal(t, buf.String(), data.expected)
		})
	}
}

func TestFmt_Warn(t *testing.T) {
	timestamp := time.Now()
	provider := new(mockProvider)
	provider.On("Now").Return(timestamp)
	expected := createLogPattern("WARN   ", timestamp, colorize.Yellow)
	for _, data := range []struct {
		input    LogLevel
		expected string
	}{
		{Error, expected},
		{Warn, expected},
		{Verbose, ""},
		{Notice, ""},
		{Silent, ""},
	} {
		t.Run(data.input.String(), func(t *testing.T) {
			// arrange
			var buf bytes.Buffer
			logger := New(data.input, &buf, provider)

			// arc
			logger.Warn("message")

			// assert
			assert.Equal(t, buf.String(), data.expected)
		})
	}
}

func TestFmt_Log(t *testing.T) {
	timestamp := time.Now()
	provider := new(mockProvider)
	provider.On("Now").Return(timestamp)
	expected := createLogPattern("NOTICE ", timestamp, colorize.DarkGreen)
	for _, data := range []struct {
		input    LogLevel
		expected string
	}{
		{Error, expected},
		{Warn, expected},
		{Notice, expected},
		{Verbose, ""},
		{Silent, ""},
	} {
		t.Run(data.input.String(), func(t *testing.T) {
			// arrange
			var buf bytes.Buffer
			logger := New(data.input, &buf, provider)

			// arc
			logger.Log("message")

			// assert
			assert.Equal(t, buf.String(), data.expected)
		})
	}
}

func TestFmt_Debug(t *testing.T) {
	timestamp := time.Now()
	provider := new(mockProvider)
	provider.On("Now").Return(timestamp)
	expected := createLogPattern("VERBOSE", timestamp, colorize.Cyan)
	for _, data := range []struct {
		input    LogLevel
		expected string
	}{
		{Error, expected},
		{Warn, expected},
		{Notice, expected},
		{Verbose, expected},
		{Silent, ""},
	} {
		t.Run(data.input.String(), func(t *testing.T) {
			// arrange
			var buf bytes.Buffer
			logger := New(data.input, &buf, provider)

			// arc
			logger.Debug("message")

			// assert
			assert.Equal(t, buf.String(), data.expected)
		})
	}
}

func TestFmt_SetLogLevel(t *testing.T) {
	// arrange
	logger := Default(Error)

	// act
	logger.SetLogLevel(Warn)

	// assert
	assert.Equal(t, logger.LogLevel(), Warn)
}
