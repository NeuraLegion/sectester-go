package colorize

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnsiColorize_Colorize_ReturnsInputWithColors(t *testing.T) {
	// arrange
	input := "input"
	console := NewMockConsole(t)
	console.EXPECT().IsColored().Return(true)
	colorize := New(console)

	// act
	got := colorize.Colorize(Red, input)

	// assert
	assert.Equal(t, got, Red.String()+input+DefaultForeground.String())
}

func TestAnsiColorize_Colorize_ReturnsInput(t *testing.T) {
	// arrange
	input := "input"
	console := NewMockConsole(t)
	console.EXPECT().IsColored().Return(false)
	colorize := New(console)

	// act
	got := colorize.Colorize(Red, input)

	// assert
	assert.Equal(t, got, input)
}
