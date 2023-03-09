package http

import (
	"net/url"
	"time"
)

type Options struct {
	BaseUrl *url.URL
	Token   string
	Timeout time.Duration
}
