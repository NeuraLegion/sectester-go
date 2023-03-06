package bus

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

const expectedName = "TestCommand"

type mockDispatcher struct {
	mock.Mock
}

func (m *mockDispatcher) Publish(message *Message) error {
	args := m.Called(message)

	return args.Error(0)
}

func (m *mockDispatcher) Execute(message *Message) (any, error) {
	args := m.Called(message)

	return args.Get(0), args.Error(1)
}

func TestNewMessage_DefaultOptions(t *testing.T) {
	// act
	var got, err = NewMessage(expectedName)

	// assert
	if err != nil {
		t.Fatal(err)
	}

	assert.Condition(t, func() (success bool) {
		return got.Name() == expectedName &&
			got.CreatedAt().After(time.Time{}) &&
			got.ExpectReply() == true &&
			got.Ttl() == time.Minute &&
			len(got.CorrelationId()) != 0
	})
}

func TestWithTtl_TtlIsGreaterThan0(t *testing.T) {
	// arrange
	var ttl = time.Duration(10) * time.Second

	// act
	var got, err = NewMessage(expectedName, WithTtl(ttl))

	// assert
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, got.Ttl(), ttl)
}

func TestWithTtl_TtlIsLessThanOrEquals0(t *testing.T) {
	// arrange
	var ttl = time.Duration(0)

	// act
	var _, err = NewMessage(expectedName, WithTtl(ttl))

	// assert
	assert.EqualError(t, err, "ttl must be greater than 0")
}

func TestWithExpectReply(t *testing.T) {
	// act
	var got, err = NewMessage(expectedName, WithExpectReply(false))

	// assert
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, got.ExpectReply(), false)
}

func TestMessage_Publish(t *testing.T) {
	// arrange
	m, err := NewMessage(expectedName, WithPayload("test"))
	if err != nil {
		t.Fatal(err)
	}
	dispatcher := new(mockDispatcher)
	dispatcher.On("Publish", m).Return(nil)

	// act
	err = m.Publish(dispatcher)

	// assert
	if err != nil {
		t.Fatal(err)
	}
	dispatcher.AssertExpectations(t)
}

func TestMessage_Execute_ReturnsResult(t *testing.T) {
	// arrange
	m, err := NewMessage(expectedName, WithPayload("test"))
	if err != nil {
		t.Fatal(err)
	}
	dispatcher := new(mockDispatcher)
	expected := "result"
	dispatcher.On("Execute", m).Return(expected, nil)

	// act
	result, err := m.Execute(dispatcher)

	// assert
	if err != nil {
		t.Fatal(err)
	}
	dispatcher.AssertExpectations(t)
	assert.Equal(t, result, expected)
}

func TestMessage_Execute(t *testing.T) {
	// arrange
	m, err := NewMessage(expectedName, WithPayload("test"))
	if err != nil {
		t.Fatal(err)
	}
	dispatcher := new(mockDispatcher)
	dispatcher.On("Execute", m).Return(nil, nil)

	// act
	_, err = m.Execute(dispatcher)

	// assert
	if err != nil {
		t.Fatal(err)
	}
	dispatcher.AssertExpectations(t)
}

func TestMessage_Execute_ReturnsError(t *testing.T) {
	// arrange
	m, err := NewMessage(expectedName, WithPayload("test"))
	if err != nil {
		t.Fatal(err)
	}
	dispatcher := new(mockDispatcher)
	dispatcher.On("Execute", m).Return(nil, errors.New("something went wrong"))

	// act
	_, err = m.Execute(dispatcher)

	// assert
	dispatcher.AssertExpectations(t)
	assert.EqualError(t, err, "something went wrong")
}
