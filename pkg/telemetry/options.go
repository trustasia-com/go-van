// Package telemetry provides ...
package telemetry

import (
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
)

// FlagOption to flag with 0/1
type FlagOption int

// flag list
const (
	// secure for grpc
	FlagInsecure FlagOption = 1 << iota
	// opentelemetry tracing
	FlagTracer
	// opentelemetry logger
	FlagLogger
	// opentelemetry metrics
	FlagMeter
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
	// custom resource attributes that will be applied to all telemetry data
	attributes []attribute.KeyValue

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
func WithFlag(flag FlagOption) Option {
	return func(opts *options) { opts.flag |= flag }
}

// WithOptions otlpgrpc options
func WithOptions(dialOpts ...grpc.DialOption) Option {
	return func(opts *options) {
		opts.options = append(opts.options, dialOpts...)
	}
}

// WithAttributes add custom resource attributes
// These attributes will be applied to all telemetry data (traces, metrics, logs)
func WithAttributes(attrs ...attribute.KeyValue) Option {
	return func(opts *options) {
		opts.attributes = append(opts.attributes, attrs...)
	}
}
