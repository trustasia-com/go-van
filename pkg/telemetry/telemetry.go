// Package telemetry provides ...
package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/trustasia-com/go-van/pkg/logx"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"google.golang.org/grpc"
)

// examples:
//   https://github.com/open-telemetry/opentelemetry-go-contrib
//   https://github.com/open-telemetry/opentelemetry-go
//

// shutdownFunc func
type shutdownFunc func(context.Context) error

// InitProvider init telemetry provider
func InitProvider(ctx context.Context, opts ...Option) (shutdown func(), flag FlagOption) {
	options := options{flag: DefaultStdFlag}
	// apply opts
	for _, o := range opts {
		o(&options)
	}

	var (
		err                            error
		tracerShutdown, metricShutdown shutdownFunc
	)
	// gRPC connection
	conn, err := grpc.NewClient(options.endpoint, options.options...)
	if err != nil {
		logx.Fatal(err)
	}
	// resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// The service name used to display traces in backends
			semconv.ServiceNameKey.String(options.name),
		),
	)
	if err != nil {
		logx.Fatal(err)
	}
	// meter
	if options.flag&FlagMeter > 0 {
		metricShutdown, err = initMeter(ctx, res, conn)
		if err != nil {
			logx.Fatal(err)
		}
	}
	// tracer
	if options.flag&FlagTracer > 0 {
		tracerShutdown, err = initTracer(ctx, res, conn)
		if err != nil {
			logx.Fatal(err)
		}
	}
	// logger
	//
	shutdown = func() {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
		defer cancel()

		if tracerShutdown != nil {
			if err = tracerShutdown(ctx); err != nil {
				logx.Fatalf("failed to shutdown tracer: %v", err)
			}
		}
		if metricShutdown != nil {
			if err = metricShutdown(ctx); err != nil {
				logx.Fatalf("failed to shutdown metric: %v", err)
			}
		}
	}
	return shutdown, options.flag
}

// initMetric metric provider
//
// eg.
//
//	appName := "example-api"
//	meter := otel.Meter(appName)
//	opt := api.WithAttributes(
//		attribute.Key("A").String("B"),
//		attribute.Key("C").String("D"),
//	)
//	counter, err := meter.Float64Counter("foo", api.WithDescription("a simple counter"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	counter.Add(ctx, 5, opt)
func initMeter(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (shutdownFunc, error) {
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown, nil
}

// initTracer trace provider
func initTracer(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (shutdownFunc, error) {
	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}
