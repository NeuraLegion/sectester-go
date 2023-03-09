package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogLevel_String(t *testing.T) {
	// arrange
	level := Error
	// act
	result := level.String()
	// assert
	assert.Equal(t, result, "error")
}

func TestLogLevel_Humanize(t *testing.T) {
	// arrange
	type testData struct {
		input    LogLevel
		expected string
	}
	levels := []testData{
		{input: Error, expected: "ERROR  "},
		{input: Warn, expected: "WARN   "},
		{input: Notice, expected: "NOTICE "},
		{input: Verbose, expected: "VERBOSE"},
	}
	// act
	for _, data := range levels {
		t.Run(data.input.String(), func(t *testing.T) {
			// arc
			got := data.input.Humanize()

			// assert
			assert.Equal(t, got, data.expected)
		})
	}
}
