package clock

import "time"

// SystemProvider implements the Provider interface using the time.Now.
type SystemProvider struct{}

// Now returns the current time.
func (s *SystemProvider) Now() time.Time {
	return time.Now()
}
