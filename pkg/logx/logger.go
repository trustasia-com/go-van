// Package logx provides ...
package logx

import (
	"context"
	"time"
)

// Logger represents a logger.
type Logger interface {
	// Info logs to INFO log.
	Info(args ...interface{})
	// Info logs to INFO log.
	Infof(format string, args ...interface{})
	// Warning logs to WARNING log.
	Warning(args ...interface{})
	// Warning logs to WARNING log.
	Warningf(format string, args ...interface{})
	// Error logs to ERROR log.
	Error(args ...interface{})
	// Error logs to ERROR log.
	Errorf(format string, args ...interface{})
	// Fatal logs to ERROR log. with os.Exit(1).
	Fatal(args ...interface{})
	// Fatal logs to ERROR log. with os.Exit(1).
	Fatalf(format string, args ...interface{})
	// V reports whether verbosity level log is at least the requested verbose level.
	V(level int) bool
}

// list log level
const (
	LevelInfo Level = iota
	LevelWarning
	LevelError
	LevelFatal

	numSeverity = 4
)

// level string
var levelName = []string{
	"INFO",
	"WARNING",
	"ERROR",
	"FATAL",
}

// Level logger level
type Level int

func (l Level) String() string { return levelName[l] }

// Entry log entry
type Entry struct {
	Level   Level                  `json:"level"`
	Time    time.Time              `json:"time"`
	File    string                 `json:"file"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`

	logging *Logging        `json:"-"`
	context context.Context `json:"-"`
}

// Output print log
func (e *Entry) Output(calldepth int) {
	defer func() {

	}()

}
