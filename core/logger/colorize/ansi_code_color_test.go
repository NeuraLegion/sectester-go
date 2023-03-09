package colorize

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnsiCodeColor_String(t *testing.T) {
	// arrange
	color := Red
	// act
	result := color.String()
	// assert
	assert.Equal(t, result, "\u001B[1m\u001B[31m")
}
