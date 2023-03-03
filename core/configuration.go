package core

import (
	"fmt"
	"net/url"
	"regexp"
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
	Bus               string `exhaustruct:"optional"`
	Api               string `exhaustruct:"optional"`
	Name              string
	Version           string
	RepeaterVersion   string
	loopbackAddresses []string
}

func NewConfiguration(hostname string) (*Configuration, error) {
	c := &Configuration{
		Name:              Name,
		Version:           Version,
		RepeaterVersion:   RepeaterVersion,
		loopbackAddresses: []string{"localhost", "127.0.0.1"},
	}
	err := c.resolveUrls(hostname)
	if err != nil {
		return nil, err
	}
	return c, nil
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
			c.Bus = fmt.Sprintf("amqp://%s:5672", host)
			c.Api = fmt.Sprintf("http://%s:8000", host)
			return nil
		}
	}
	c.Bus = fmt.Sprintf("amqps://amq.%s:5672", host)
	c.Api = fmt.Sprintf("https://%s", host)
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
