// Package telemetry provides ...
package telemetry

import (
	"context"
	"time"

	"github.com/trustasia-com/go-van/pkg/logx"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
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
	// init exporter
	exp, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(opts.endpoint),
		otlptracegrpc.WithDialOption(opts.options...),
	)

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			// service name
			semconv.ServiceNameKey.String(opts.name),
		),
	)
	if err != nil {
		return nil, err
	}
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

var globalMeter metric.Meter

// initMetric metric provider
func initMetric(ctx context.Context, opts options) (shutdownFunc, error) {
	// init exporter
	exp, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(opts.endpoint),
		otlpmetricgrpc.WithDialOption(opts.options...),
	)

	pusher := controller.New(
		processor.NewFactory(
			simple.NewWithHistogramDistribution(),
			exp,
		),
		controller.WithCollectPeriod(7*time.Second),
		controller.WithExporter(exp),
	)
	// set global propagator to tracecontext (the default is no-op).
	global.SetMeterProvider(pusher)
	err = pusher.Start(context.Background())
	if err != nil {
		return nil, err
	}
	globalMeter = pusher.Meter(opts.name)
	return pusher.Stop, nil
}
