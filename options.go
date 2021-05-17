// Package van provides ...
package van

import (
	"context"

	"github.com/deepzz0/go-van/registry"
	"github.com/deepzz0/go-van/server"
)

// Option one option
type Option func(*options)

// options for micro service
type options struct {
	id      string
	name    string
	version string
	ctx     context.Context
	signal  bool

	metadata  map[string]string
	registry  registry.Registry
	servers   []server.Server
	endpoints []string
}

func defaultOptions() options {
	return options{
		ctx:    context.Background(),
		signal: true,
	}
}

// Name service name
func Name(name string) Option {
	return func(opts *options) { opts.name = name }
}

// Version service version
func Version(ver string) Option {
	return func(opts *options) { opts.version = ver }
}

// Context specifc service context
func Context(ctx context.Context) Option {
	return func(opts *options) { opts.ctx = ctx }
}

// Signal specific service os signal
func Signal(b bool) Option {
	return func(opts *options) { opts.signal = b }
}

// Metadata service metadata
func Metadata(md map[string]string) Option {
	return func(opts *options) { opts.metadata = md }
}

// Registry sets the registry for the services
func Registry(r registry.Registry) Option {
	return func(opts *options) { opts.registry = r }
}

// Server used for service
func Server(ss ...server.Server) Option {
	return func(opts *options) {
		opts.servers = append(opts.servers, ss...)
	}
}

// Endpoint sets service endpoints
func Endpoint(eps ...string) Option {
	return func(opts *options) {
		opts.endpoints = append(opts.endpoints, eps...)
	}
}
