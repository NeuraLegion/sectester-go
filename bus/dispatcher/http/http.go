package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/NeuraLegion/sectester-go/core/bus"
)

const (
	// AuthScheme is the authentication scheme used for the HTTP request.
	AuthScheme = "api-key %s"
	// CorrelationIdHeaderField is the name of the header field used for the correlation ID of the HTTP request.
	CorrelationIdHeaderField = "x-correlation-id"
	// DateHeaderField is the name of the header field used for the date of the HTTP request.
	DateHeaderField = "date"
	// AuthorizationHeaderField is the name of the header field used for the authorization token of the HTTP request.
	AuthorizationHeaderField = "authorization"
)

type executionResult struct {
	result any
	err    error
}

// Http represents an HTTP implementation of the bus.CommandDispatcher interface.
// Http is an alternative way to execute the commands over HTTP. To start, you should create
// a Http instance by passing the following options to the constructor:
//
//	dispatcher: = &Http{
//		client: &http.Client{},
//		options: Options{Token: config.Credentials().Token(), BaseUrl: config.Api()},
//	}
type Http struct {
	client  *http.Client
	options Options
}

// Execute sends an HTTP request based on the provided bus.Message and returns the response, unless otherwise is stated.
// Before, you have to create an instance of bus.Message with an instance of Request as payload:
//
//	p := map[string]any{"foo": "bar"}
//	d, err := json.Marshal(p)
//	if err != nil {
//		t.Fatal(err)
//	}
//	m, err := bus.NewMessage("CreateRepeater", bus.WithPayload(
//		Request{
//			Url: "/api/v1/repeaters",
//			Method: http.MethodPost,
//			Body: io.NopCloser(bytes.NewReader(d))
//		}
//	))
//
// Once it is done, you can perform a request using Execute as follows:
//
//	response := dispatcher.Execute(m)
func (h *Http) Execute(message *bus.Message) (any, error) {
	ch := make(chan *executionResult)

	go func(message *bus.Message, ch chan<- *executionResult) {
		defer close(ch)
		res, err := h.sendRequest(message)
		if err != nil {
			ch <- &executionResult{result: nil, err: err}
			return
		}
		defer res.Body.Close()
		if !message.ExpectReply() {
			ch <- &executionResult{result: nil, err: nil}
			return
		}
		result, err := h.parseResponse(res)
		ch <- &executionResult{result: result, err: err}
	}(message, ch)

	select {
	case r := <-ch:
		return r.result, r.err
	case <-time.After(message.Ttl()):
		return nil, fmt.Errorf("no response for %s", message.Ttl().String())
	}
}

func (h *Http) normalizeRequestOptions(message *bus.Message) *Request {
	options := message.Payload().(Request) //nolint:errcheck // false positive finding
	if options.Headers == nil {
		options.Headers = http.Header{}
	}
	options.Headers.Add(CorrelationIdHeaderField, message.CorrelationId())
	options.Headers.Add(DateHeaderField, message.CreatedAt().Format(time.RFC3339))
	return &options
}

func (h *Http) parseResponse(res *http.Response) (any, error) {
	err := h.throwIfUnsuccessful(res)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return h.parseBody(body)
}

func (h *Http) parseBody(body []byte) (any, error) {
	if len(body) == 0 {
		return nil, nil
	}
	var result any
	err := json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (h *Http) throwIfUnsuccessful(res *http.Response) error {
	if res.StatusCode < http.StatusBadRequest {
		return nil
	}
	if !h.canObtainErrorMessage(res) {
		return fmt.Errorf("request failed with status code %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return errors.New(string(body))
}

func (h *Http) sendRequest(message *bus.Message) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), message.Ttl())
	defer cancel()

	options := h.normalizeRequestOptions(message)
	uri := h.buildUrl(options)

	req, _ := http.NewRequestWithContext(ctx, options.Method, uri.String(), options.Body)

	h.addHeaders(req, options)

	return h.client.Do(req)
}

func (h *Http) buildUrl(options *Request) *url.URL {
	u := h.options.BaseUrl.JoinPath(options.Url)
	query := u.Query()
	for key, value := range options.Params {
		query.Set(key, value)
	}
	u.RawQuery = query.Encode()
	return u
}

func (h *Http) addHeaders(request *http.Request, options *Request) {
	for key, value := range options.Headers {
		for _, val := range value {
			request.Header.Add(key, val)
		}
	}

	request.Header.Set(AuthorizationHeaderField, fmt.Sprintf(AuthScheme, h.options.Token))
}

func (h *Http) canObtainErrorMessage(response *http.Response) bool {
	contentType := response.Header.Get("content-type")
	allowedContentTypes := &[]string{
		"text/html", "text/plain",
	}

	for _, v := range *allowedContentTypes {
		if strings.HasPrefix(contentType, v) {
			return response.ContentLength > 0
		}
	}

	return false
}
