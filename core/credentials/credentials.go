package credentials

import (
	"errors"
	"regexp"
)

var tokenValidationRegexp = regexp.MustCompile(`^[A-Za-z0-9+/=]{7}\.nex[apr]\.[A-Za-z0-9+/=]{32}$`)

type Credentials struct {
	token string
}

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
