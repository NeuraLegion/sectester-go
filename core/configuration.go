package core

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/NeuraLegion/sectester-go/core/credentials"
)

const (
	Name            = "sectester-go"
	Version         = "0.0.1"
	RepeaterVersion = "9.3.0"
)

var (
	schemaRegex                = regexp.MustCompile(`^.+://`)
	hostnameNormalizationRegex = regexp.MustCompile(`^((?:\w+:)?//)|^(//)?`)
)

type Configuration struct {
	name              string
	version           string
	repeaterVersion   string
	loopbackAddresses []string
	bus               string                   `exhaustruct:"optional"`
	api               string                   `exhaustruct:"optional"`
	credentials       *credentials.Credentials `exhaustruct:"optional"`
}

type ConfigurationOption func(f *Configuration)

func WithCredentials(credentials *credentials.Credentials) ConfigurationOption {
	return func(f *Configuration) {
		f.credentials = credentials
	}
}

func NewConfiguration(hostname string, opts ...ConfigurationOption) (*Configuration, error) {
	c := &Configuration{
		name:              Name,
		version:           Version,
		repeaterVersion:   RepeaterVersion,
		loopbackAddresses: []string{"localhost", "127.0.0.1"},
	}
	err := c.resolveUrls(hostname)
	if err != nil {
		return nil, err
	}
	for _, applyOpt := range opts {
		applyOpt(c)
	}
	return c, nil
}

func (c *Configuration) Credentials() *credentials.Credentials {
	return c.credentials
}

func (c *Configuration) Api() string {
	return c.api
}

func (c *Configuration) Bus() string {
	return c.bus
}

func (c *Configuration) Name() string {
	return c.name
}

func (c *Configuration) Version() string {
	return c.version
}

func (c *Configuration) RepeaterVersion() string {
	return c.repeaterVersion
}

func (c *Configuration) normalizeHostname(hostname string) (string, error) {
	uri, err := url.Parse(c.addSchema(hostname))

	if err != nil {
		return "", err
	}

	return uri.Hostname(), nil
}

func (c *Configuration) resolveUrls(hostname string) error {
	host, err := c.normalizeHostname(hostname)
	if err != nil {
		return err
	}
	for _, a := range c.loopbackAddresses {
		if a == host {
			c.bus = fmt.Sprintf("amqp://%s:5672", host)
			c.api = fmt.Sprintf("http://%s:8000", host)
			return nil
		}
	}
	c.bus = fmt.Sprintf("amqps://amq.%s:5672", host)
	c.api = fmt.Sprintf("https://%s", host)
	return nil
}

func (c *Configuration) addSchema(hostname string) string {
	if !schemaRegex.MatchString(hostname) {
		return hostnameNormalizationRegex.ReplaceAllString(
			hostname,
			"https://",
		)
	}
	return hostname
}
