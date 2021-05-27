// Package telemetry provides ...
package telemetry

import (
	"context"
	"time"

	"github.com/deepzz0/go-van/pkg/logx"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"google.golang.org/grpc"
)

// examples:
//   https://github.com/open-telemetry/opentelemetry-go-contrib
//

// shutdownFunc func
type shutdownFunc func(context.Context) error

// InitProvider init telemetry provider
func InitProvider(opts ...Option) (shutdown func()) {
	options := options{
		context: context.Background(),
	}
	// apply opts
	for _, o := range opts {
		o(&options)
	}
	// init exporter
	options.options = append(options.options, otlpgrpc.WithDialOption(grpc.WithBlock()))
	exp, err := otlp.NewExporter(options.context, otlpgrpc.NewDriver(
		options.options...,
	))
	if err != nil {
		logx.Fatal(err)
	}

	var tracerShutdown, metricShutdown shutdownFunc
	// tracer
	if options.tracerName != "" {
		tracerShutdown, err = initTracer(exp, options)
		if err != nil {
			logx.Fatal(err)
		}
	}
	// metrics
	if options.metrics {
		metricShutdown, err = initMetric(exp, options)
		if err != nil {
			logx.Fatal(err)
		}
	}
	// logger
	//
	shutdown = func() {
		if tracerShutdown != nil {
			tracerShutdown(options.context)
		}
		if metricShutdown != nil {
			metricShutdown(options.context)
		}
		exp.Shutdown(options.context)
	}
	return shutdown
}

// initTracer trace provider
func initTracer(exp *otlp.Exporter, opts options) (shutdownFunc, error) {
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			// service name
			semconv.ServiceNameKey.String(opts.tracerName),
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
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tracerProvider)
	return tracerProvider.Shutdown, nil
}

// initMetric metric provider
func initMetric(exp *otlp.Exporter, opts options) (shutdownFunc, error) {
	cont := controller.New(
		processor.New(
			simple.NewWithExactDistribution(),
			exp,
		),
		controller.WithCollectPeriod(7*time.Second),
		controller.WithExporter(exp),
	)
	// set global propagator to tracecontext (the default is no-op).
	global.SetMeterProvider(cont.MeterProvider())
	err := cont.Start(context.Background())
	if err != nil {
		return nil, err
	}
	return cont.Stop, nil
}
