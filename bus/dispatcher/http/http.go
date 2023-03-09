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
	AuthScheme               = "api-key %s"
	CorrelationIdHeaderField = "x-correlation-id"
	Date                     = "date"
	AuthorizationHeaderField = "authorization"
)

type executionResult struct {
	result any
	err    error
}

type Http struct {
	client  *http.Client
	options Options
}

func (h *Http) Execute(message *bus.Message) (any, error) {
	ch := make(chan *executionResult)

	go func(message *bus.Message) {
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
	}(message)

	select {
	case r := <-ch:
		if r.err != nil {
			return nil, r.err
		}

		return r.result, nil
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
	options.Headers.Add(Date, message.CreatedAt().Format(time.RFC3339))
	return &options
}

func (h *Http) parseResponse(res *http.Response) (any, error) {
	err := h.throwIfUnsuccessful(res)
	if err != nil {
		return nil, err
	}
	body, err := h.readBody(res)
	if err != nil {
		return nil, err
	}
	return h.parseBody(body)
}

func (h *Http) readBody(res *http.Response) ([]byte, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
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
	body, err := h.readBody(res)
	if err != nil {
		return err
	}
	return errors.New(string(body))
}

func (h *Http) sendRequest(message *bus.Message) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), message.Ttl())
	defer cancel()

	options := h.normalizeRequestOptions(message)
	u := h.buildUrl(options)

	r, _ := http.NewRequestWithContext(ctx, options.Method, u.String(), options.Body)

	h.addHeaders(r, options)

	return h.client.Do(r)
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
