package http

import (
	"io"
	"net/http"
)

// Request represents a payload for a command, including URL, method, parameters, body and headers.
type Request struct {
	// Url is the URL of the request.
	Url string
	// Method is the HTTP method of the request, such as GET, POST, PUT, DELETE, etc.
	Method string
	// Params is a map of query parameters for the request.
	Params map[string]string
	// Body is an io.ReadCloser that contains the body of the request.
	Body io.ReadCloser
	// Headers is a http.Header that contains the headers of the request.
	Headers http.Header
}
