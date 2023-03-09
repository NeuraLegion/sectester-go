package logger

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type LogLevel int

const (
	Silent LogLevel = iota
	Error
	Warn
	Notice
	Verbose
)

var humanizedLevels = []string{
	"silent", "error", "warn", "notice", "verbose",
}

func (s LogLevel) String() string {
	return humanizedLevels[s]
}

func maxLength() int {
	copied := make([]string, len(humanizedLevels))
	copy(copied, humanizedLevels)
	sort.Slice(copied, func(i, j int) bool {
		return len(copied[i]) >= len(copied[j])
	})
	return len(copied[0])
}

func (s LogLevel) Humanize() string {
	template := fmt.Sprintf("%-"+strconv.Itoa(maxLength())+"s", s.String())
	return strings.ToUpper(template)
}

type Logger interface {
	LogLevel() LogLevel
	Error(message string, args ...any)
	Warn(message string, args ...any)
	Log(message string, args ...any)
	Debug(message string, args ...any)
}
