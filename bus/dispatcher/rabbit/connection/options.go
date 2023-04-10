package connection

import "time"

type Options struct {
	Url               string
	ConnectTimeout    time.Duration
	HeartbeatInterval time.Duration
	ReconnectTime     time.Duration
	Username          string
	Password          string
}
