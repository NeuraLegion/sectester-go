package clock

import "time"

type SystemProvider struct{}

func (s *SystemProvider) Now() time.Time {
	return time.Now()
}
