package bus

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewEvent_DefaultOptions(t *testing.T) {
	// arrange
	const eventName = "TestEvent"

	// act
	var got, err = NewEvent(eventName)

	// assert
	if err != nil {
		t.Fatal(err)
	}

	assert.Condition(t, func() (success bool) {
		return got.Name() == eventName &&
			got.CreatedAt().After(time.Time{}) &&
			len(got.CorrelationId()) != 0
	})
}
