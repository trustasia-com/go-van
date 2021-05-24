// Package logx provides ...
package logx

import (
	"io"
	"sync"
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

// Logging logging setup.
type Logging struct {
	lock   sync.Mutex
	prefix string
	level  Level

	writer []io.WriteCloser
}

// Info logs to INFO log.
func (log *Logging) Info(args ...interface{}) {

}

// Info logs to INFO log.
func (log *Logging) Infof(format string, args ...interface{}) {

}

// Warning logs to WARNING log.
func (log *Logging) Warning(args ...interface{}) {

}

// Warning logs to WARNING log.
func (log *Logging) Warningf(format string, args ...interface{}) {

}

// Error logs to ERROR log.
func (log *Logging) Error(args ...interface{}) {

}

// Error logs to ERROR log.
func (log *Logging) Errorf(format string, args ...interface{}) {

}

// Fatal logs to ERROR log. with os.Exit(1).
func (log *Logging) Fatal(args ...interface{}) {

}

// Fatal logs to ERROR log. with os.Exit(1).
func (log *Logging) Fatalf(format string, args ...interface{}) {

}

// V reports whether verbosity level log is at least the requested verbose level.
func (log *Logging) V(l Level) bool {
	return l < log.level
}
