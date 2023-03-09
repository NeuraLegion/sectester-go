package http

import (
	"io"
	"net/http"
)

type Request struct {
	Url     string
	Method  string
	Params  map[string]string
	Body    io.ReadCloser
	Headers http.Header
}
