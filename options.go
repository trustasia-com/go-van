// Package van provides ...
package van

import (
	"context"

	"github.com/deepzz0/go-van/pkg/registry"
	"github.com/deepzz0/go-van/pkg/server"
)

// Option one option
type Option func(*options)

// options for micro service
type options struct {
	// listen os signal
	signal bool
	// service id, auto generate
	id string
	// service name
	name string
	// service version
	version string
	// other options for implementations of the interface
	// can be stored in a context
	context context.Context

	// some metadata
	metadata map[string]string
	// registry for discovery
	registry registry.Registry
	// service server list
	servers []server.Server
	// specific servers endpoints
	endpoints []string
}

// WithSignal specific service os signal
func WithSignal(b bool) Option {
	return func(opts *options) { opts.signal = b }
}

// WithName service name
func WithName(name string) Option {
	return func(opts *options) { opts.name = name }
}

// WithVersion service version
func WithVersion(ver string) Option {
	return func(opts *options) { opts.version = ver }
}

// WithContext specifc service context
func WithContext(ctx context.Context) Option {
	return func(opts *options) { opts.context = ctx }
}

// WithMetadata service metadata
func WithMetadata(md map[string]string) Option {
	return func(opts *options) { opts.metadata = md }
}

// WithRegistry sets the registry for the services
func WithRegistry(r registry.Registry) Option {
	return func(opts *options) { opts.registry = r }
}

// WithServer used for service
func WithServer(ss ...server.Server) Option {
	return func(opts *options) {
		opts.servers = append(opts.servers, ss...)
	}
}

// WithEndpoint sets service endpoints
func WithEndpoint(eps ...string) Option {
	return func(opts *options) {
		opts.endpoints = append(opts.endpoints, eps...)
	}
}
