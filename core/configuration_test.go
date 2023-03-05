package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/NeuraLegion/sectester-go/core/credentials"
)

type mockProvider struct {
	mock.Mock
}

func (p *mockProvider) Get() *credentials.Credentials {
	args := p.Called()
	res := args.Get(0)

	if res == nil {
		return nil
	}

	c, ok := res.(*credentials.Credentials)

	if !ok {
		panic(fmt.Sprintf("unable to cast the resulting object: %v", args.Get(0)))
	}

	return c
}

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
	var cred, _ = credentials.New("weobbz5.nexa.vennegtzr2h7urpxgtksetz2kwppdgj0")

	// act
	got, _ := NewConfiguration(hostname, WithCredentials(cred))

	// assert
	assert.Equal(t, got.Credentials(), cred)
}

func TestWithCredentialsProviders(t *testing.T) {
	// arrange
	var hostname = "app.brightsec.com"
	var cred, _ = credentials.New("weobbz5.nexa.vennegtzr2h7urpxgtksetz2kwppdgj0")
	var provider = new(mockProvider)
	var providers = []credentials.Provider{
		provider,
	}
	provider.On("Get").Return(cred)

	// act
	got, _ := NewConfiguration(hostname, WithCredentialsProviders(providers))

	// assert
	assert.Equal(t, got.Credentials(), cred)
}

func TestWithCredentialsProviders_NoCredentials(t *testing.T) {
	// arrange
	var hostname = "app.brightsec.com"
	var providers []credentials.Provider

	// act
	_, err := NewConfiguration(hostname, WithCredentialsProviders(providers))

	// assert
	assert.EqualError(t, err, "please provide either 'credentials' or 'credentialProviders'")
}

func TestWithCredentialsProviders_UnableToFindCredentials(t *testing.T) {
	// arrange
	var hostname = "app.brightsec.com"
	var provider = new(mockProvider)
	var providers = []credentials.Provider{
		provider,
	}
	provider.On("Get").Return(nil)

	// act
	_, err := NewConfiguration(hostname, WithCredentialsProviders(providers))

	// assert
	assert.EqualError(t, err, "could not load cred from any providers")
}

func TestWithCredentialsProviders_MultipleProvidersReturnCredentials(t *testing.T) {
	// arrange
	var hostname = "app.brightsec.com"
	var cred1, _ = credentials.New("weobbz5.nexa.vennegtzr2h7urpxgtksetz2kwppdgj0")
	var cred2, _ = credentials.New("weobbz5.nexa.vennegtzr2h7urpxgtksetz2kwppdgj1")
	var provider1 = new(mockProvider)
	var provider2 = new(mockProvider)
	var providers = []credentials.Provider{
		provider1,
		provider2,
	}
	provider1.On("Get").Return(cred1)
	provider2.On("Get").Return(cred2)

	// act
	got, _ := NewConfiguration(hostname, WithCredentialsProviders(providers))

	// assert
	assert.Equal(t, got.Credentials(), cred1)
}

func TestNewConfiguration_EmptyArrayOfCredentialsProviders(t *testing.T) {
	// arrange
	var hostname = "app.brightsec.com"
	var cred, _ = credentials.New("weobbz5.nexa.vennegtzr2h7urpxgtksetz2kwppdgj0")
	var providers []credentials.Provider

	// act
	got, _ := NewConfiguration(hostname, WithCredentials(cred), WithCredentialsProviders(providers))

	// assert
	assert.Equal(t, got.Credentials(), cred)
}

func TestNewConfiguration_ValidHostname(t *testing.T) {
	// arrange
	type resolvedHost struct {
		Bus string
		Api string
	}

	type testData struct {
		Input    string
		Expected resolvedHost
	}

	var cred, _ = credentials.New("weobbz5.nexa.vennegtzr2h7urpxgtksetz2kwppdgj0")
	var hostnames = []testData{
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
			got, _ := NewConfiguration(data.Input, WithCredentials(cred))

			// assert
			assert.Condition(t, func() (success bool) {
				return got.Bus() == data.Expected.Bus && got.Api() == data.Expected.Api
			})
		})
	}
}
