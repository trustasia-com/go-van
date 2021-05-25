// Package logx provides ...
package logx

import (
	"io"
)

// flag list
const (
	FlagFile = 1 << iota

	stdFlags = FlagFile
)

// Option logger option
type Option func(opts *Options)

// Options logger options
type Options struct {
	prefix string    // service name
	level  Level     // print severity
	writer io.Writer // writer

	flag int // log flag
}

// WithName set service name
func WithName(name string) Option {
	return func(opts *Options) { opts.prefix = name }
}

// WithLevel set severity level
func WithLevel(l Level) Option {
	return func(opts *Options) { opts.level = l }
}

// WithWriter log output
func WithWriter(w io.Writer) Option {
	return func(opts *Options) { opts.writer = w }
}

// WithFlag options flag
func WithFlag(flag int) Option {
	return func(opts *Options) { opts.flag = flag }
}
