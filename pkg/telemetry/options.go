// Package telemetry provides ...
package telemetry

import (
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
)

// Signal 表示一个 OpenTelemetry Signal。
type Signal uint8

// Signal 列表。
const (
	SignalTracer Signal = 1 << iota
	SignalLogger
	SignalMeter
)

// Enabled 判断指定 Signal 是否全部启用。
func (signals Signal) Enabled(signal Signal) bool {
	return signal != 0 && signals&signal == signal
}

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

	// OpenTelemetry Signals
	signals Signal
	// OTLP gRPC 是否使用明文连接
	insecure bool
}

// WithEndpoint opentelemetry backend endpoint
func WithEndpoint(edp string) Option {
	return func(opts *options) { opts.endpoint = edp }
}

// WithName open with name
func WithName(name string) Option {
	return func(opts *options) { opts.name = name }
}

// WithSignals 选择需要初始化的 OpenTelemetry Signals。
func WithSignals(signals ...Signal) Option {
	return func(opts *options) {
		for _, signal := range signals {
			opts.signals |= signal
		}
	}
}

// WithInsecure 使用明文 OTLP gRPC Transport。
func WithInsecure() Option {
	return func(opts *options) { opts.insecure = true }
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
