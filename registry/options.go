// Package registry provides ...
package registry

import (
	"context"
	"crypto/tls"
	"time"
)

// Option regsitry option
type Option func(o *Options)

// Options registry Options
type Options struct {
	Ctx       context.Context
	Addrs     []string    // backend endpoint
	TLSConfig *tls.Config // whether use tls
	TTL       time.Duration

	// service info
	name      string
	version   string
	metadata  map[string]string
	endpoints []string
}

// Context register with context
func Context(ctx context.Context) Option {
	return func(opts *Options) { opts.Ctx = ctx }
}

// Addr registry address to use
func Addr(addrs ...string) Option {
	return func(opts *Options) {
		opts.Addrs = append(opts.Addrs, addrs...)
	}
}

// TLSConfig registry secure tlc config
func TLSConfig(tls *tls.Config) Option {
	return func(opts *Options) { opts.TLSConfig = tls }
}

// TTL register ttl
func TTL(ttl time.Duration) Option {
	return func(opts *Options) { opts.TTL = ttl }
}

// Name register with name
func Name(name string) Option {
	return func(opts *Options) { opts.name = name }
}

// Version register with version
func Version(ver string) Option {
	return func(opts *Options) { opts.version = ver }
}

// Metadata register with metadata
func Metadata(md map[string]string) Option {
	return func(opts *Options) { opts.metadata = md }
}

// Endpoint register with endpoint
func Endpoint(eps ...string) Option {
	return func(opts *Options) {
		opts.endpoints = append(opts.endpoints, eps...)
	}
}
