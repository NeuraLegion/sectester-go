package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider_Get_EnvVariableIsNotProvided(t *testing.T) {
	// arrange
	provider := new(Provider)

	// act
	result := provider.Get()

	// assert
	assert.Nil(t, result)
}

func TestProvider_Get_ReturnCredentials(t *testing.T) {
	// arrange
	const token = "0zmcwpe.nexr.0vlon8mp7lvxzjuvgjy88olrhadhiukk"
	oldValue := os.Getenv(BrightToken)
	provider := new(Provider)
	_ = os.Setenv(BrightToken, token)

	// act
	result := provider.Get()
	_ = os.Setenv(BrightToken, oldValue)

	// assert
	assert.Condition(t, func() (success bool) {
		return result.Token() == token
	})
}
