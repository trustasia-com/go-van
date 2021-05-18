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
	Addresses []string    // backend endpoint
	TLSConfig *tls.Config // whether use tls
	TTL       time.Duration
}

// Context register with context
func Context(ctx context.Context) Option {
	return func(opts *Options) { opts.Ctx = ctx }
}

// Address registry address to use
func Address(addrs ...string) Option {
	return func(opts *Options) {
		opts.Addresses = append(opts.Addresses, addrs...)
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
