package rabbit

import (
	"github.com/NeuraLegion/sectester-go/bus/dispatcher/rabbit/connection"
)

type Options struct {
	connection.Options
	AppQueue      string
	Exchange      string
	ClientQueue   string
	PrefetchCount int
}
