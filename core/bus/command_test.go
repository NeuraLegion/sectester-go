package bus

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const commandName = "TestCommand"

func TestNewCommand_DefaultOptions(t *testing.T) {
	// act
	var got, err = NewCommand(commandName)

	// assert
	if err != nil {
		t.Fatal(err)
	}

	assert.Condition(t, func() (success bool) {
		return got.Name() == commandName &&
			got.ExpectReply() == true &&
			got.Ttl() == time.Minute &&
			got.CreatedAt().After(time.Time{}) &&
			len(got.CorrelationId()) != 0
	})
}

func TestWithTtl_TtlIsGreaterThan0(t *testing.T) {
	// arrange
	var ttl = time.Duration(10) * time.Second

	// act
	var command, err = NewCommand(commandName, WithTtl(ttl))

	// assert
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, command.Ttl(), ttl)
}

func TestWithTtl_TtlIsLessThanOrEquals0(t *testing.T) {
	// arrange
	var ttl = time.Duration(0)

	// act
	var _, err = NewCommand(commandName, WithTtl(ttl))

	// assert
	assert.EqualError(t, err, "ttl must be greater than 0")
}

func TestWithExpectReply(t *testing.T) {
	// act
	var command, err = NewCommand(commandName, WithExpectReply(false))

	// assert
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, command.ExpectReply(), false)
}
