// Package logx provides ...
package logx

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sync"
)

var std = NewLogging()

// NewLogging log logger
func NewLogging(opts ...Option) *Logging {
	options := Options{
		level:  LevelInfo,
		writer: os.Stderr,
		flag:   0, // default flag
	}
	// apply opts
	for _, o := range opts {
		o(&options)
	}

	// new logging
	logging := &Logging{options: options}
	// object pool
	logging.entryPool.New = logging.newEntry
	logging.bufferPool.New = logging.newBuffer
	return logging
}

// Logging logging setup.
type Logging struct {
	mu         sync.Mutex
	entryPool  sync.Pool
	bufferPool sync.Pool

	options Options
}

// newEntry new entry
func (log *Logging) newEntry() any {
	return &Entry{
		logging: log,
		Data:    make(map[string]any, 4),
	}
}

// releaseEntry release entry
func (log *Logging) releaseEntry(e *Entry) {
	e.Data = map[string]any{}
	log.entryPool.Put(e)
}

// newBuffer new buffer
func (log *Logging) newBuffer() any {
	return new(bytes.Buffer)
}

// releaseBuffer release buffer
func (log *Logging) releaseBuffer(buf *bytes.Buffer) {
	buf.Reset()
	log.bufferPool.Put(buf)
}

// output output log
func (log *Logging) output(l Level, msg string) {
	if !log.V(l) {
		return
	}
	// get & put
	entry, _ := log.entryPool.Get().(*Entry)
	entry.Level = l
	entry.Message = msg

	calldepth := 4
	entry.Output(calldepth)
}

// Debug logs to DEBUG log.
func (log *Logging) Debug(args ...any) {
	log.output(LevelDebug, fmt.Sprintln(args...))
}

// Debugf logs to DEBUG log.
func (log *Logging) Debugf(format string, args ...any) {
	log.output(LevelDebug, fmt.Sprintf(format, args...))
}

// Info logs to INFO log.
func (log *Logging) Info(args ...any) {
	log.output(LevelInfo, fmt.Sprintln(args...))
}

// Infof logs to INFO log.
func (log *Logging) Infof(format string, args ...any) {
	log.output(LevelInfo, fmt.Sprintf(format, args...))
}

// Warning logs to WARNING log.
func (log *Logging) Warning(args ...any) {
	log.output(LevelWarning, fmt.Sprintln(args...))
}

// Warningf logs to WARNING log.
func (log *Logging) Warningf(format string, args ...any) {
	log.output(LevelWarning, fmt.Sprintf(format, args...))
}

// Error logs to ERROR log.
func (log *Logging) Error(args ...any) {
	log.output(LevelError, fmt.Sprintln(args...))
}

// Errorf logs to ERROR log.
func (log *Logging) Errorf(format string, args ...any) {
	log.output(LevelError, fmt.Sprintf(format, args...))
}

// Fatal logs to ERROR log. with os.Exit(1).
func (log *Logging) Fatal(args ...any) {
	log.output(LevelFatal, fmt.Sprintln(args...))
	os.Exit(1)
}

// Fatalf logs to ERROR log. with os.Exit(1).
func (log *Logging) Fatalf(format string, args ...any) {
	log.output(LevelFatal, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// SetLevel change logging options level
func (log *Logging) SetLevel(lv Level) {
	if lv >= LevelDebug && lv <= LevelFatal {
		log.options.level = lv
	}
}

// V reports whether verbosity level log is at least the requested verbose level.
func (log *Logging) V(l Level) bool {
	return log.options.level <= l
}

// Debug logs to DEBUG log.
func Debug(args ...any) {
	std.Debug(args...)
}

// Debugf logs to DEBUG log.
func Debugf(format string, args ...any) {
	std.Debugf(format, args...)
}

// Info logs to INFO log.
func Info(args ...any) {
	std.Info(args...)
}

// Infof logs to INFO log.
func Infof(format string, args ...any) {
	std.Infof(format, args...)
}

// Warning logs to WARNING log.
func Warning(args ...any) {
	std.Warning(args...)
}

// Warningf logs to WARNING log.
func Warningf(format string, args ...any) {
	std.Warningf(format, args...)
}

// Error logs to ERROR log.
func Error(args ...any) {
	std.Error(args...)
}

// Errorf logs to ERROR log.
func Errorf(format string, args ...any) {
	std.Errorf(format, args...)
}

// Fatal logs to ERROR log. with os.Exit(1).
func Fatal(args ...any) {
	std.Fatal(args...)
}

// Fatalf logs to ERROR log. with os.Exit(1).
func Fatalf(format string, args ...any) {
	std.Fatalf(format, args...)
}

// SetLevel change logging options level
func SetLevel(lv Level) {
	std.SetLevel(lv)
}

// WithData custom data
func WithData(data map[string]any) *Entry {
	entry, _ := std.entryPool.Get().(*Entry)
	entry.Data = data
	return entry
}

// WithContext context
func WithContext(ctx context.Context) *Entry {
	entry, _ := std.entryPool.Get().(*Entry)
	entry.context = ctx
	return entry
}
