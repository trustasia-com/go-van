// Package server provides ...
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/deepzz0/go-van/pkg/registry"

	"google.golang.org/grpc"
)

// ServerOption server option
type ServerOption func(opts *ServerOptions)

// ServerOptions server Options
type ServerOptions struct {
	// server listen network tcp/udp
	Network string
	// server run address
	Address string
	// handler for http server
	Handler http.Handler
	// Options for gRPC server
	Options []grpc.ServerOption

	// server: recover from panic
	Recover bool
}

// WithNetwork server network
func WithNetwork(network string) ServerOption {
	return func(opts *ServerOptions) { opts.Network = network }
}

// WithAddress server endpoint
func WithAddress(addr string) ServerOption {
	return func(opts *ServerOptions) { opts.Address = addr }
}

// WithHandler server handler
func WithHandler(h http.Handler) ServerOption {
	return func(opts *ServerOptions) { opts.Handler = h }
}

// WithOptions gRPC server option
func WithOptions(sopts ...grpc.ServerOption) ServerOption {
	return func(opts *ServerOptions) {
		opts.Options = append(opts.Options, sopts...)
	}
}

// WithRecover server panic recover
func WithRecover(rec bool) ServerOption {
	return func(opts *ServerOptions) { opts.Recover = rec }
}

// DialOption client dial option
type DialOption func(opts *DialOptions)

// DialOptions client dial Options
type DialOptions struct {
	// client connect to endpoint
	Endpoint string
	// connect timeout
	Timeout time.Duration
	// user-agent
	UserAgent string
	// other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	// client: secure with tls
	Secure bool
	// client: discovery registry
	Registry registry.Registry
}

// WithEndpoint connect to server endpoint
func WithEndpoint(addr string) DialOption {
	return func(opts *DialOptions) { opts.Endpoint = addr }
}

// WithTimeout dial timeout
func WithTimeout(timeout time.Duration) DialOption {
	return func(opts *DialOptions) { opts.Timeout = timeout }
}

// WithContext client context
func WithContext(ctx context.Context) DialOption {
	return func(opts *DialOptions) { opts.Context = ctx }
}

// WithUserAgent client user-agent
func WithUserAgent(ua string) DialOption {
	return func(opts *DialOptions) { opts.UserAgent = ua }
}

// WithSecure endpoint secure
func WithSecure(secure bool) DialOption {
	return func(opts *DialOptions) { opts.Secure = secure }
}

// WithRegistry registry for discovery
func WithRegistry(reg registry.Registry) DialOption {
	return func(opts *DialOptions) { opts.Registry = reg }
}
