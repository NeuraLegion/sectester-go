package credentials

import (
	"errors"
	"regexp"
)

var tokenValidationRegexp = regexp.MustCompile(`^[A-Za-z0-9+/=]{7}\.nex[apr]\.[A-Za-z0-9+/=]{32}$`)

// A Credentials allows to set credentials to access the application.
// More info about [setting up an API key].
//
// [setting up an API key]: https://docs.brightsec.com/docs/manage-your-personal-account#manage-your-personal-api-keys-authentication-tokens
//
//nolint:lll // linter does not respect a long URL in the docs
type Credentials struct {
	token string
}

// A New creates a new instance of Credentials.
//
// var cred, _ = credentials.New("your API key")
// var config = core.NewConfiguration("app.brightsec.com", core.WithCredentials(cred).
func New(token string) (*Credentials, error) {
	err := validate(token)

	if err != nil {
		return nil, err
	}

	return &Credentials{token: token}, nil
}

func (c *Credentials) Token() string {
	return c.token
}

func validate(token string) error {
	if token == "" {
		return errors.New("provide an API key")
	}

	if !tokenValidationRegexp.MatchString(token) {
		return errors.New("unable to recognize the API key")
	}

	return nil
}
