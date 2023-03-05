package env

import (
	"os"

	"github.com/NeuraLegion/sectester-go/core/credentials"
)

const BrightToken = "BRIGHT_TOKEN"

// A Provider reads credentials from the following environment variable: 'BRIGHT_TOKEN'
//
// var p = new(env.Provider)
// var config = core.NewConfiguration("app.brightsec.com", core.WithCredentialsProviders([]credentials.Provider { p })).
type Provider struct{}

// Get returns a instance of credentials.Credentials.
// If the 'BRIGHT_TOKEN' environment variable is not set or contains a falsy value, it will return undefined.
func (p Provider) Get() *credentials.Credentials {
	token := os.Getenv(BrightToken)
	c, err := credentials.New(token)

	if err != nil {
		return nil
	}

	return c
}
