package logger

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type LogLevel int

type Logger interface {
	LogLevel() LogLevel
	Error(message string, args ...any)
	Warn(message string, args ...any)
	Log(message string, args ...any)
	Debug(message string, args ...any)
}

const (
	Silent LogLevel = iota
	Verbose
	Notice
	Warn
	Error
)

func (s LogLevel) String() string {
	return humanizedLevel(s.Index())
}

func (s LogLevel) Index() int {
	return int(s)
}

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
