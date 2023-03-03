package core

import (
	"github.com/NeuraLegion/sectester-go/core/credentials"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfiguration_HostnameIsInvalid(t *testing.T) {
	// arrange
	var hostname = ":test"

	// act
	got, _ := NewConfiguration(hostname)

	// assert
	assert.Nil(t, got)
}

func TestWithCredentials(t *testing.T) {
	// arrange
	var hostname = "app.brightsec.com"
	var cred = &credentials.Credentials{}

	// act
	got, _ := NewConfiguration(hostname, WithCredentials(cred))

	// assert
	assert.Equal(t, got.credentials, cred)
}

func TestNewConfiguration_ValidHostname_ResolvesApiAndBus(t *testing.T) {
	// arrange
	type resolvedHost struct {
		Bus string
		Api string
	}
	type testData struct {
		Input    string
		Expected resolvedHost
	}

	hostnames := []testData{
		{Input: "localhost", Expected: resolvedHost{Bus: "amqp://localhost:5672", Api: "http://localhost:8000"}},
		{Input: "localhost:8080", Expected: resolvedHost{Bus: "amqp://localhost:5672", Api: "http://localhost:8000"}},
		{Input: "http://localhost", Expected: resolvedHost{Bus: "amqp://localhost:5672", Api: "http://localhost:8000"}},
		{Input: "http://localhost:8080", Expected: resolvedHost{Bus: "amqp://localhost:5672", Api: "http://localhost:8000"}},
		{Input: "127.0.0.1", Expected: resolvedHost{Bus: "amqp://127.0.0.1:5672", Api: "http://127.0.0.1:8000"}},
		{Input: "127.0.0.1:8080", Expected: resolvedHost{Bus: "amqp://127.0.0.1:5672", Api: "http://127.0.0.1:8000"}},
		{Input: "http://127.0.0.1", Expected: resolvedHost{Bus: "amqp://127.0.0.1:5672", Api: "http://127.0.0.1:8000"}},
		{Input: "http://127.0.0.1:8080", Expected: resolvedHost{Bus: "amqp://127.0.0.1:5672", Api: "http://127.0.0.1:8000"}},
		{Input: "example.com", Expected: resolvedHost{Bus: "amqps://amq.example.com:5672", Api: "https://example.com"}},
		{Input: "example.com:443", Expected: resolvedHost{Bus: "amqps://amq.example.com:5672", Api: "https://example.com"}},
		{Input: "http://example.com", Expected: resolvedHost{Bus: "amqps://amq.example.com:5672", Api: "https://example.com"}},
		{Input: "http://example.com:443", Expected: resolvedHost{Bus: "amqps://amq.example.com:5672", Api: "https://example.com"}},
	}

	// act
	for _, data := range hostnames {
		t.Run(data.Input, func(t *testing.T) {
			// arc
			got, _ := NewConfiguration(data.Input)

			// assert
			assert.Condition(t, func() (success bool) {
				return got.Bus() == data.Expected.Bus && got.Api() == data.Expected.Api
			})
		})
	}
}
