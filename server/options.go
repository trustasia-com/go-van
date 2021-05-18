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
	// server listen network tcp/udp
	// client dial network
	Network string
	// server run endpoint
	// client connect to endpoint
	Endpoint string
	// connect timeout
	Timeout time.Duration

	// open trace middleware
	Trace bool
	// other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// WithNetwork server network
func WithNetwork(network string) Option {
	return func(opts *Options) { opts.Network = network }
}

// WithEndpoint server endpoint
func WithEndpoint(addr string) Option {
	return func(opts *Options) { opts.Endpoint = addr }
}

// WithTimeout server timeout
func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) { opts.Timeout = timeout }
}

// WithTrace server trace
func WithTrace(trace bool) Option {
	return func(opts *Options) { opts.Trace = trace }
}

// WithContext server context
func WithContext(ctx context.Context) Option {
	return func(opts *Options) { opts.Context = ctx }
}
