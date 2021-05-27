// Package telemetry provides ...
package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
)

// Option telemetry option
type Option func(opts *Options)

// options registry Options
type options struct {
	// connect to backend store, maybe is a cluster
	endpoint string
	// tracer name
	tracerName string
	// export metrics
	metrics bool
	// otel options
	options []otlpgrpc.Option
	// other options for implementations of the interface
	// can be stored in a context
	context context.Context
}

// WithEndpoint opentelemetry backend endpoint
func WithEndpoint(edp string) Option {
	return func(opts *options) { opts.endpoint = edp }
}

// WithTracerName open tracer with name
func WithTracerName(name string) Option {
	return func(opts *options) { opts.tracerName = name }
}

// WithMetrics open metrics
func WithMetrics() Option {
	return func(opts *options) { opts.metrics = true }
}

// WithOptions otlpgrpc options
func WithOptions(oopts ...otlpgrpc.Option) Option {
	return func(opts *options) {
		opts.options = append(opts.options, oopts...)
	}
}

// WithContext context
func WithContext(ctx context.Context) Option {
	return func(opts *options) { opts.context = ctx }
}
