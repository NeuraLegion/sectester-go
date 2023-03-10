package core

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/NeuraLegion/sectester-go/core/credentials"
	"github.com/NeuraLegion/sectester-go/core/credentials/env"
	"github.com/NeuraLegion/sectester-go/core/logger"
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

// Configuration allows you to configure the SDK including credentials from different sources.
// The default configuration is as follows:
//
//	&Configuration { credentialProviders: []credentials.Provider{ new(env.Provider), logLevel: logger.Error } }.
type Configuration struct {
	name                string
	version             string
	repeaterVersion     string
	loopbackAddresses   []string
	bus                 string                   `exhaustruct:"optional"`
	api                 string                   `exhaustruct:"optional"`
	credentials         *credentials.Credentials `exhaustruct:"optional"`
	credentialProviders []credentials.Provider
	logLevel            logger.LogLevel
}

type ConfigurationOption func(f *Configuration)

// WithCredentials sets credentials.Credentials to access the application.
//
//	var config = core.NewConfiguration("app.brightsec.com", core.WithCredentials(/* your credentials */).
func WithCredentials(credentials *credentials.Credentials) ConfigurationOption {
	return func(c *Configuration) {
		c.credentials = credentials
	}
}

// WithCredentialsProviders allows you to provide a list of credentials.Provider to load credentials.Credentials
// in runtime.
//
//	var config = core.NewConfiguration("app.brightsec.com", core.WithCredentialsProviders(/* your providers */).
func WithCredentialsProviders(providers []credentials.Provider) ConfigurationOption {
	return func(c *Configuration) {
		c.credentialProviders = providers
	}
}

// WithLogLevel allows you to say what level of logs to report. Any logs of a higher level than the setting are shown.
func WithLogLevel(level logger.LogLevel) ConfigurationOption {
	return func(c *Configuration) {
		c.logLevel = level
	}
}

// NewConfiguration creates a new instance of Configuration.
// Requires the application name (domain name), that is used to establish connection with.
//
//	var config = NewConfiguration("app.neuralegion.com").
func NewConfiguration(hostname string, opts ...ConfigurationOption) (*Configuration, error) {
	c := &Configuration{
		name:                Name,
		version:             Version,
		repeaterVersion:     RepeaterVersion,
		loopbackAddresses:   []string{"localhost", "127.0.0.1"},
		credentialProviders: []credentials.Provider{new(env.Provider)},
		logLevel:            logger.Error,
	}
	err := c.resolveUrls(hostname)
	if err != nil {
		return nil, err
	}
	for _, applyOpt := range opts {
		applyOpt(c)
	}
	err = c.loadCredentials()
	if err != nil {
		return nil, err
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

func (c *Configuration) LogLevel() logger.LogLevel {
	return c.logLevel
}

func (c *Configuration) loadCredentials() error {
	if c.credentials == nil && (c.credentialProviders == nil || len(c.credentialProviders) == 0) {
		return errors.New("please provide either 'credentials' or 'credentialProviders'")
	}

	if c.credentials != nil {
		return nil
	}

	cred, err := c.discoverCredentials()

	if err != nil {
		return err
	}

	c.credentials = cred

	return nil
}

func (c *Configuration) discoverCredentials() (*credentials.Credentials, error) {
	var cred *credentials.Credentials

	for _, provider := range c.credentialProviders {
		cred = provider.Get()

		if cred != nil {
			break
		}
	}

	if cred == nil {
		return nil, errors.New("could not load cred from any providers")
	}

	return cred, nil
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
