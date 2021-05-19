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
	// connect to backend store, maybe is a cluster
	Enpoints []string
	// whether use secure tls
	TLSConfig *tls.Config
	// time to live of heartbeat
	TTL time.Duration
	// other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// WithEndpoint registry address to use
func WithEndpoint(endpoints ...string) Option {
	return func(opts *Options) {
		opts.Enpoints = append(opts.Enpoints, endpoints...)
	}
}

// WithTLS registry secure tlc config
func WithTLS(tls *tls.Config) Option {
	return func(opts *Options) { opts.TLSConfig = tls }
}

// WithTTL register ttl
func WithTTL(ttl time.Duration) Option {
	return func(opts *Options) { opts.TTL = ttl }
}

// WithContext register with context
func WithContext(ctx context.Context) Option {
	return func(opts *Options) { opts.Context = ctx }
}
