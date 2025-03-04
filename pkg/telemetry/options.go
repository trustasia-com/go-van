// Package telemetry provides ...
package telemetry

import (
	"google.golang.org/grpc"
)

// FlagOption to flag with 0/1
type FlagOption int

// flag list
const (
	// opentelemetry metrics
	FlagMeter = 1 << iota
	// opentelemetry tracing
	FlagTracer

	DefaultStdFlag = FlagMeter | FlagTracer
)

// Option telemetry option
type Option func(opts *options)

// options telemetry Options
type options struct {
	// connect to backend store, maybe is a cluster
	endpoint string
	// app name
	name string
	// otel collector options
	options []grpc.DialOption

	// opentelemetry switch
	flag FlagOption
}

// WithEndpoint opentelemetry backend endpoint
func WithEndpoint(edp string) Option {
	return func(opts *options) { opts.endpoint = edp }
}

// WithName open with name
func WithName(name string) Option {
	return func(opts *options) { opts.name = name }
}

// WithFlag opentelemetry switch
func WithFlag(flags ...FlagOption) Option {
	return func(opts *options) {
		for _, f := range flags {
			opts.flag |= f
		}
	}
}

// WithOptions otlpgrpc options
func WithOptions(dialOpts ...grpc.DialOption) Option {
	return func(opts *options) {
		opts.options = append(opts.options, dialOpts...)
	}
}
