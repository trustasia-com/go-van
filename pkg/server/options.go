// Package server provides ...
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/trustasia-com/go-van/pkg/registry"
	"github.com/trustasia-com/go-van/pkg/telemetry"

	"google.golang.org/grpc"
)

// FlagOption to flag with 0/1
type FlagOption int

// flag list
const (
	// server: recover from panic
	FlagRecover FlagOption = 1 << iota
	// client: tracing with telemetry
	FlagTracing
	// client: secure for tls
	FlagInsecure
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
	// other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	// server flag
	Flag FlagOption

	// telemetry options
	Telemetry []telemetry.Option
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

// WithTelemetry opentelemetry
func WithTelemetry(topts ...telemetry.Option) ServerOption {
	return func(opts *ServerOptions) {
		opts.Telemetry = append(opts.Telemetry, topts...)
	}
}

// WithSrvFlag server flag
func WithSrvFlag(flag FlagOption) ServerOption {
	return func(opts *ServerOptions) { opts.Flag |= flag }
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

	// discovery registry to server
	Registry registry.Registry
	// client flag
	Flag FlagOption
}

// WithEndpoint connect to server endpoint
func WithEndpoint(addr string) DialOption {
	return func(opts *DialOptions) { opts.Endpoint = addr }
}

// WithTimeout dial timeout
func WithTimeout(timeout time.Duration) DialOption {
	return func(opts *DialOptions) { opts.Timeout = timeout }
}

// WithUserAgent client user-agent
func WithUserAgent(ua string) DialOption {
	return func(opts *DialOptions) { opts.UserAgent = ua }
}

// WithRegistry registry for discovery
func WithRegistry(reg registry.Registry) DialOption {
	return func(opts *DialOptions) { opts.Registry = reg }
}

// WithCliFlag client flag
func WithCliFlag(flag FlagOption) DialOption {
	return func(opts *DialOptions) { opts.Flag |= flag }
}
