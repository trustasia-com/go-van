// Package telemetry provides ...
package telemetry

import (
	"google.golang.org/grpc"
)

// Option telemetry option
type Option func(opts *options)

// options registry Options
type options struct {
	// connect to backend store, maybe is a cluster
	endpoint string
	// name
	name string
	// export metrics
	metrics bool
	// otel tracer options
	options []grpc.DialOption
}

// WithEndpoint opentelemetry backend endpoint
func WithEndpoint(edp string) Option {
	return func(opts *options) { opts.endpoint = edp }
}

// WithName open with name
func WithName(name string) Option {
	return func(opts *options) { opts.name = name }
}

// WithMetrics open metrics
func WithMetrics(metrics bool) Option {
	return func(opts *options) { opts.metrics = metrics }
}

// WithOptions otlpgrpc options
func WithOptions(dialOpts ...grpc.DialOption) Option {
	return func(opts *options) {
		opts.options = append(opts.options, dialOpts...)
	}
}
