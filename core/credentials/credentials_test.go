package credentials

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew_InvalidToken(t *testing.T) {
	// arrange
	const token string = "qwerty"

	// act
	_, err := New(token)

	// assert
	assert.Error(t, err, "unable to recognize the API key")
}

func TestNew_TokenNotDefined(t *testing.T) {
	// arrange
	const token string = ""

	// act
	_, err := New(token)

	// assert
	assert.Error(t, err, "provide an API key")
}

func TestCredentials_Token(t *testing.T) {
	// arrange
	tokens := []string{
		"weobbz5.nexa.vennegtzr2h7urpxgtksetz2kwppdgj0",
		"w0iikzf.nexp.aeish9lhiag7ledmsdwpwcbilagupc3r",
		"0zmcwpe.nexr.0vlon8mp7lvxzjuvgjy88olrhadhiukk",
	}

	for _, token := range tokens {
		t.Run(token, func(t *testing.T) {
			// act
			got, _ := New(token)

			// assert
			assert.Equal(t, got.Token(), token)
		})
	}
}
