// Package telemetry provides ...
package telemetry

import (
	"context"
	"time"

	"github.com/trustasia-com/go-van/pkg/logx"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"
)

// examples:
//   https://github.com/open-telemetry/opentelemetry-go-contrib
//

// shutdownFunc func
type shutdownFunc func(context.Context) error

// InitProvider init telemetry provider
func InitProvider(ctx context.Context, opts ...Option) (shutdown func()) {
	options := options{}
	// apply opts
	for _, o := range opts {
		o(&options)
	}

	var (
		err                            error
		tracerShutdown, metricShutdown shutdownFunc
	)
	// tracer
	if options.name != "" {
		tracerShutdown, err = initTracer(ctx, options)
		if err != nil {
			logx.Fatal(err)
		}
	}
	// metrics
	if options.metrics {
		metricShutdown, err = initMetric(ctx, options)
		if err != nil {
			logx.Fatal(err)
		}
	}
	// logger
	//
	shutdown = func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*20)
		defer cancel()

		if tracerShutdown != nil {
			tracerShutdown(ctx)
		}
		if metricShutdown != nil {
			metricShutdown(ctx)
		}
	}
	return shutdown
}

// initTracer trace provider
func initTracer(ctx context.Context, opts options) (shutdownFunc, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(opts.name),
		),
	)
	if err != nil {
		return nil, err
	}
	// grpc conn
	conn, err := grpc.Dial(opts.endpoint, opts.options...)
	if err != nil {
		return nil, err
	}
	// init exporter
	exp, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithGRPCConn(conn),
	)

	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	// set global propagator to tracecontext (the default is no-op).
	ctmp := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{},
		propagation.Baggage{})
	otel.SetTextMapPropagator(ctmp)
	otel.SetTracerProvider(tracerProvider)
	return tracerProvider.Shutdown, nil
}

// initMetric metric provider
func initMetric(ctx context.Context, opts options) (shutdownFunc, error) {
	// grpc conn
	conn, err := grpc.Dial(opts.endpoint, opts.options...)
	if err != nil {
		return nil, err
	}
	// init exporter
	exp, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithGRPCConn(conn),
	)
	provider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exp)),
	)

	// set global propagator to tracecontext (the default is no-op).
	global.SetMeterProvider(provider)
	return provider.Shutdown, nil
}
