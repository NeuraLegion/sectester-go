package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/NeuraLegion/sectester-go/core/bus"
)

const token = "weobbz5.nexa.vennegtzr2h7urpxgtksetz2kwppdgj0"

func TestHttp_Execute_SendsAuthorizationHeader(t *testing.T) {
	// arrange
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, fmt.Sprintf("api-key %s", token), r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	m, err := bus.NewMessage("test", bus.WithPayload(Request{}))
	if err != nil {
		t.Fatal(err)
	}
	h := &Http{
		client:  &http.Client{},
		options: Options{Token: token, BaseUrl: u},
	}

	// act
	result, err := h.Execute(m)

	// assert
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestHttp_Execute_SendsCorrelationIdHeader(t *testing.T) {
	// arrange
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotEmpty(t, r.Header.Get("x-correlation-id"))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	m, err := bus.NewMessage("test", bus.WithPayload(Request{}))
	if err != nil {
		t.Fatal(err)
	}
	h := &Http{
		client:  &http.Client{},
		options: Options{Token: token, BaseUrl: u},
	}

	// act
	result, err := h.Execute(m)

	// assert
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestHttp_Execute_SendsDateHeader(t *testing.T) {
	// arrange
	m, err := bus.NewMessage("test", bus.WithPayload(Request{}))
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Header.Get("date"), m.CreatedAt().Format(time.RFC3339))
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	h := &Http{
		client:  &http.Client{},
		options: Options{Token: token, BaseUrl: u},
	}

	// act
	result, err := h.Execute(m)

	// assert
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestHttp_Execute_SendsWithQueryString(t *testing.T) {
	// arrange
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "baz=xyzzy", r.URL.RawQuery)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	m, err := bus.NewMessage("test", bus.WithPayload(Request{
		Params: map[string]string{"baz": "xyzzy"},
	}))
	if err != nil {
		t.Fatal(err)
	}
	h := &Http{
		client:  &http.Client{},
		options: Options{Token: token, BaseUrl: u},
	}

	// act
	result, err := h.Execute(m)

	// assert
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestHttp_Execute_AppendsQueryToExistingInUri(t *testing.T) {
	// arrange
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "baz=xyzzy&foo=bar", r.URL.RawQuery)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	query := u.Query()
	query.Add("foo", "bar")
	u.RawQuery = query.Encode()
	m, err := bus.NewMessage("test", bus.WithPayload(Request{
		Params: map[string]string{"baz": "xyzzy"},
	}))
	if err != nil {
		t.Fatal(err)
	}
	h := &Http{
		client:  &http.Client{},
		options: Options{Token: token, BaseUrl: u},
	}

	// act
	result, err := h.Execute(m)

	// assert
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestHttp_Execute_SendsBody(t *testing.T) {
	// arrange
	data := map[string]any{"foo": "bar"}
	expected, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, body, expected)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	m, err := bus.NewMessage(
		"test",
		bus.WithPayload(
			Request{Body: io.NopCloser(bytes.NewReader(expected))},
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	h := &Http{
		client:  &http.Client{},
		options: Options{Token: token, BaseUrl: u},
	}

	// act
	result, err := h.Execute(m)

	// assert
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestHttp_Execute_ReturnsReply(t *testing.T) {
	// arrange
	expected := map[string]any{"foo": "bar"}
	body, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(body)
	}))
	defer ts.Close()
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	m, err := bus.NewMessage("test", bus.WithPayload(Request{}))
	if err != nil {
		t.Fatal(err)
	}
	h := &Http{
		client:  &http.Client{},
		options: Options{Token: token, BaseUrl: u},
	}

	// act
	result, err := h.Execute(m)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, result, expected)
}

func TestHttp_ExecuteWithoutExpectingReply_ReturnsNil(t *testing.T) {
	// arrange
	expected := map[string]any{"foo": "bar"}
	body, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(body)
	}))
	defer ts.Close()
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	m, err := bus.NewMessage("test", bus.WithPayload(Request{}), bus.WithExpectReply(false))
	if err != nil {
		t.Fatal(err)
	}
	h := &Http{
		client:  &http.Client{},
		options: Options{Token: token, BaseUrl: u},
	}

	// act
	result, err := h.Execute(m)

	// assert
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestHttp_Execute_ReturnsErrorByTimeout(t *testing.T) {
	// arrange
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 5)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	m, err := bus.NewMessage("test", bus.WithPayload(Request{}), bus.WithTtl(time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}
	h := &Http{
		client:  &http.Client{},
		options: Options{Token: token, BaseUrl: u},
	}

	// act
	result, err := h.Execute(m)

	// assert
	assert.ErrorContains(t, err, fmt.Sprintf("no response for %s", m.Ttl().String()))
	assert.Nil(t, result)
}

func TestHttp_Execute_ReturnsError(t *testing.T) {
	// arrange
	expected := "something went wrong"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(expected))
	}))
	defer ts.Close()
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	m, err := bus.NewMessage("test", bus.WithPayload(Request{}))
	if err != nil {
		t.Fatal(err)
	}
	h := &Http{
		client:  &http.Client{},
		options: Options{Token: token, BaseUrl: u},
	}

	// act
	result, err := h.Execute(m)

	// assert
	assert.ErrorContains(t, err, expected)
	assert.Nil(t, result)
}

func TestHttp_Execute_ReturnsDefaultError(t *testing.T) {
	// arrange
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	m, err := bus.NewMessage("test", bus.WithPayload(Request{}))
	if err != nil {
		t.Fatal(err)
	}
	h := &Http{
		client:  &http.Client{},
		options: Options{Token: token, BaseUrl: u},
	}

	// act
	result, err := h.Execute(m)

	// assert
	assert.ErrorContains(t, err, "request failed with status code 400")
	assert.Nil(t, result)
}
