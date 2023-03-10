package clock

import "time"

// Provider is an interface for retrieving the current time.
type Provider interface {
	Now() time.Time
}
