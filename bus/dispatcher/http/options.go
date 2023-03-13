package http

import (
	"net/url"
	"time"
)

// Options allows to customize Http command dispatcher.
type Options struct {
	// BaseUrl for your application instance, e.g. https://app.brightsec.com
	BaseUrl *url.URL
	// Token is API key to access the API. Find out how to obtain [personal] and
	// [organization] API keys in the knowledgebase.
	//
	// [personal]: https://docs.brightsec.com/docs/manage-your-personal-account#manage-your-personal-api-keys-authentication-tokens
	// [organization]: https://docs.brightsec.com/docs/manage-your-organization#manage-organization-apicli-authentication-tokens
	//
	//nolint:lll // linter does not respect a long URL in the docs
	Token string
	// Timeout is time to wait for a server to send response headers (and start the
	// response body) before aborting the request. Default 10s
	Timeout time.Duration
}
