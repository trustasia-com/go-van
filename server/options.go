// Package server provides ...
package server

import (
	"context"
	"time"
)

// Option server option.
type Option func(o *Options)

// Options registry Options
type Options struct {
	Network string
	Address string
	Timeout time.Duration
	Trace   bool
	Ctx     context.Context
}

// Network server network
func Network(network string) Option {
	return func(opts *Options) { opts.Network = network }
}

// Address server address
func Address(addr string) Option {
	return func(opts *Options) { opts.Address = addr }
}

// Timeout server timeout
func Timeout(timeout time.Duration) Option {
	return func(opts *Options) { opts.Timeout = timeout }
}

// Trace server trace
func Trace(trace bool) Option {
	return func(opts *Options) { opts.Trace = trace }
}
