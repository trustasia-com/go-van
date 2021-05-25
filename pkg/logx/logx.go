// Package logx provides ...
package logx

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var std = NewLogging()

// NewLogging log printer
func NewLogging(opts ...Option) *Logging {
	options := Options{
		level:  LevelInfo,
		writer: os.Stderr,
		flag:   log.LstdFlags | log.Lshortfile,
	}
	// apply opts
	for _, o := range opts {
		o(&options)
	}
	// new logging
	return &Logging{options: options}
}

// Logging logging setup.
type Logging struct {
	entryPool  sync.Pool
	bufferPool sync.Pool

	options Options
}

// newEntry return obj of entry
func (log *Logging) newEntry() *Entry {
	entry, ok := log.entryPool.Get().(*Entry)
	if ok {
		return entry
	}
	return &Entry{
		logging: log,
		Data:    make(map[string]interface{}),
	}
}

// releaseEntry release entry to pool
func (log *Logging) releaseEntry(entry *Entry) {
	entry.Data = map[string]interface{}{}
	log.entryPool.Put(entry)
}

// output output log
func (log *Logging) output(calldepth int, l Level, msg string) {
	if !log.V(l) {
		return
	}
	entry := log.newEntry()
	defer log.releaseEntry(entry)

	entry.Level = l
	entry.Time = time.Now()
	entry.Message = msg
	entry.Output(calldepth)
}

// Info logs to INFO log.
func (log *Logging) Info(args ...interface{}) {
	log.output(2, LevelInfo, fmt.Sprintln(args...))
}

// Info logs to INFO log.
func (log *Logging) Infof(format string, args ...interface{}) {
	log.output(2, LevelInfo, fmt.Sprintf(format, args...))
}

// Warning logs to WARNING log.
func (log *Logging) Warning(args ...interface{}) {
	log.output(2, LevelWarning, fmt.Sprintln(args...))
}

// Warning logs to WARNING log.
func (log *Logging) Warningf(format string, args ...interface{}) {
	log.output(2, LevelWarning, fmt.Sprintf(format, args...))
}

// Error logs to ERROR log.
func (log *Logging) Error(args ...interface{}) {
	log.output(2, LevelError, fmt.Sprintln(args...))
}

// Error logs to ERROR log.
func (log *Logging) Errorf(format string, args ...interface{}) {
	log.output(2, LevelError, fmt.Sprintf(format, args...))
}

// Fatal logs to ERROR log. with os.Exit(1).
func (log *Logging) Fatal(args ...interface{}) {
	log.output(2, LevelFatal, fmt.Sprintln(args...))
	os.Exit(1)
}

// Fatal logs to ERROR log. with os.Exit(1).
func (log *Logging) Fatalf(format string, args ...interface{}) {
	log.output(2, LevelFatal, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// V reports whether verbosity level log is at least the requested verbose level.
func (log *Logging) V(l Level) bool {
	return l <= log.options.level
}

// Info logs to INFO log.
func Info(args ...interface{}) {
	std.Info(args...)
}

// Info logs to INFO log.
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

// Warning logs to WARNING log.
func Warning(args ...interface{}) {
	std.Warning(args...)
}

// Warning logs to WARNING log.
func Warningf(format string, args ...interface{}) {
	std.Warningf(format, args...)
}

// Error logs to ERROR log.
func Error(args ...interface{}) {
	std.Error(args...)
}

// Error logs to ERROR log.
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

// Fatal logs to ERROR log. with os.Exit(1).
func Fatal(args ...interface{}) {
	std.Fatal(args...)
}

// Fatal logs to ERROR log. with os.Exit(1).
func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}
