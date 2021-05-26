// Package server provides ...
package server

import (
	"context"
	"time"

	"github.com/deepzz0/go-van/pkg/registry"
)

// Option server option.
type Option func(opts *Options)

// Options server Options, it's effective server and client
type Options struct {
	// server listen network tcp/udp
	// client dial network
	Network string
	// server run endpoint
	// client connect to endpoint
	Endpoint string
	// connect timeout
	Timeout time.Duration
	// other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	// client: secure with tls
	Secure bool
	// client: discovery registry
	Registry registry.Registry

	// server: recover from panic
	Recover bool
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

// WithContext server context
func WithContext(ctx context.Context) Option {
	return func(opts *Options) { opts.Context = ctx }
}

// WithSecure endpoint secure
func WithSecure(secure bool) Option {
	return func(opts *Options) { opts.Secure = secure }
}

// WithRegistry registry for discovery
func WithRegistry(reg registry.Registry) Option {
	return func(opts *Options) { opts.Registry = reg }
}

// WithRecover server panic recover
func WithRecover(rec bool) Option {
	return func(opts *Options) { opts.Recover = rec }
}
