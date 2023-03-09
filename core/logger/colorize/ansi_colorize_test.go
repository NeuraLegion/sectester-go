package colorize

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockConsole struct {
	mock.Mock
}

func (p *mockConsole) IsColored() bool {
	args := p.Called()
	return args.Bool(0)
}

func TestAnsiColorize_Colorize_ReturnsInputWithColors(t *testing.T) {
	// arrange
	input := "input"
	console := new(mockConsole)
	console.On("IsColored").Return(true)
	colorize := New(console)

	// act
	result := colorize.Colorize(Red, input)

	// assert
	assert.Equal(t, result, Red.String()+input+DefaultForeground.String())
}

func TestAnsiColorize_Colorize_ReturnsInput(t *testing.T) {
	// arrange
	input := "input"
	console := new(mockConsole)
	console.On("IsColored").Return(false)
	colorize := New(console)

	// act
	result := colorize.Colorize(Red, input)

	// assert
	assert.Equal(t, result, input)
}
