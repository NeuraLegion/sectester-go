package logger

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// LogLevel defines level of logs to report.
type LogLevel int

// Logger interface defines the logger that the SDK will use.
type Logger interface {
	LogLevel() LogLevel
	Error(message string, args ...any)
	Warn(message string, args ...any)
	Log(message string, args ...any)
	Debug(message string, args ...any)
}

// Predefined log levels for common scenarios.
const (
	Silent LogLevel = iota
	Verbose
	Notice
	Warn
	Error
)

// String returns the string representation of the level.
func (s LogLevel) String() string {
	return humanizedLevel(s.Index())
}

// Index returns an integer representation of level.
func (s LogLevel) Index() int {
	return int(s)
}

// Humanize returns a formatted humanized representation of level, e.g. 'NOTICE '.
func (s LogLevel) Humanize() string {
	template := fmt.Sprintf("%-"+strconv.Itoa(maxLength())+"s", s.String())
	return strings.ToUpper(template)
}

func humanizedLevels() []string {
	return []string{
		"silent", "verbose", "notice", "warn", "error",
	}
}

func humanizedLevel(idx int) string {
	return humanizedLevels()[idx]
}

func maxLength() int {
	levels := humanizedLevels()
	copied := make([]string, len(levels))
	copy(copied, levels)
	sort.Slice(copied, func(i, j int) bool {
		return len(copied[i]) >= len(copied[j])
	})
	return len(copied[0])
}
