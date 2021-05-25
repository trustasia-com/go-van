// Package logx provides ...
package logx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/trace"
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

func NewEntry(log *Logging) *Entry {
	return &Entry{
		logging: log,
		Data:    make(map[string]interface{}, 6),
	}
}

// Entry log entry
type Entry struct {
	Level   Level
	Time    time.Time
	Data    map[string]interface{}
	Message string

	logging *Logging
	context context.Context
}

// WithData custom data
func (e *Entry) WithData(data map[string]interface{}) *Entry {
	for k, v := range data {
		e.Data[k] = v
	}
	return e
}

// WithContext context
func (e *Entry) WithContext(ctx context.Context) *Entry {
	e.context = ctx
	return e
}

// Info logs to INFO log.
func (e *Entry) Info(args ...interface{}) {
	e.Level = LevelInfo
	e.Message = fmt.Sprintln(args...)
	e.Output(2)
}

// Info logs to INFO log.
func (e *Entry) Infof(format string, args ...interface{}) {
	e.Level = LevelInfo
	e.Message = fmt.Sprintf(format, args...)
	e.Output(2)
}

// Warning logs to WARNING log.
func (e *Entry) Warning(args ...interface{}) {
	e.Level = LevelWarning
	e.Message = fmt.Sprintln(args...)
	e.Output(2)
}

// Warning logs to WARNING log.
func (e *Entry) Warningf(format string, args ...interface{}) {
	e.Level = LevelWarning
	e.Message = fmt.Sprintf(format, args...)
	e.Output(2)
}

// Error logs to ERROR log.
func (e *Entry) Error(args ...interface{}) {
	e.Level = LevelError
	e.Message = fmt.Sprintln(args...)
	e.Output(2)
}

// Error logs to ERROR log.
func (e *Entry) Errorf(format string, args ...interface{}) {
	e.Level = LevelError
	e.Message = fmt.Sprintf(format, args...)
	e.Output(2)
}

// Fatal logs to ERROR log. with os.Exit(1).
func (e *Entry) Fatal(args ...interface{}) {
	e.Level = LevelFatal
	e.Message = fmt.Sprintln(args...)
	e.Output(2)
	os.Exit(1)
}

// Fatal logs to ERROR log. with os.Exit(1).
func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.Level = LevelFatal
	e.Message = fmt.Sprintln(args...)
	e.Output(2)
	os.Exit(1)
}

// Output print log
func (e *Entry) Output(calldepth int) {
	buf := e.logging.bufferPool.Get().(*bytes.Buffer)
	defer func() {
		e.logging.releaseEntry(e)
		e.logging.releaseBuffer(buf)
	}()

	// serialize
	data := make(map[string]interface{}, len(e.Data)+5)
	for k, v := range e.Data {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	data["level"] = e.Level.String()
	if e.Time.IsZero() {
		e.Time = time.Now()
	}
	data["time"] = e.Time.Format(time.RFC3339)
	data["msg"] = e.Message
	// file line
	if e.logging.options.flag&FlagFile > 0 {
		_, file, line, ok := runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				file = file[i+1:]
				break
			}
		}
		data["file"] = fmt.Sprintf("%s:%d", file, line)
	}
	// service name
	if e.logging.options.service != "" {
		data["service"] = e.logging.options.service
	}
	// opentelemetry tracing
	if e.context != nil {
		spanCtx := trace.SpanContextFromContext(e.context)
		if spanCtx.IsValid() {
			data["trace_id"] = spanCtx.SpanID().String
		}
	}

	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(data); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		return
	}
	e.logging.mu.Lock()
	defer e.logging.mu.Unlock()
	_, err := e.logging.options.writer.Write(buf.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}
}
