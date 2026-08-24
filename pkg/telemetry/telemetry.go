// Package telemetry provides ...
package telemetry

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	logglobal "go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// docs:
//   https://opentelemetry.io/docs/languages/go/
// examples:
//   https://github.com/open-telemetry/opentelemetry-go-contrib
//   https://github.com/open-telemetry/opentelemetry-go
//   https://github.com/open-telemetry/opentelemetry-collector
//

const grpcServiceConfig = `{"loadBalancingPolicy":"round_robin"}`

const providerCleanupTimeout = 10 * time.Second

// shutdownFunc func
type shutdownFunc func(context.Context) error

type signalProvider struct {
	install  func()
	shutdown shutdownFunc
}

type providerInitializer struct {
	signal     Signal
	initialize func(context.Context, *resource.Resource, *grpc.ClientConn) (signalProvider, error)
}

// Runtime 持有一次 Telemetry 初始化创建的 Provider 与 OTLP 连接。
type Runtime struct {
	signals         Signal
	shutdowns       []shutdownFunc
	closeConnection func() error
	shutdownOnce    sync.Once
	shutdownErr     error
}

func newRuntime(signals Signal, shutdowns []shutdownFunc, closeConnection func() error) *Runtime {
	return &Runtime{
		signals: signals, shutdowns: append([]shutdownFunc(nil), shutdowns...),
		closeConnection: closeConnection,
	}
}

// Enabled 判断 Runtime 是否启用了指定 Signal。
func (r *Runtime) Enabled(signal Signal) bool {
	if r == nil {
		return false
	}
	return r.signals.Enabled(signal)
}

// Shutdown 按 Provider 创建顺序的逆序释放资源；重复调用返回首次释放结果。
func (r *Runtime) Shutdown(ctx context.Context) error {
	if r == nil {
		return nil
	}
	if ctx == nil {
		return errors.New("telemetry shutdown context is nil")
	}
	r.shutdownOnce.Do(func() {
		var result error
		for index := len(r.shutdowns) - 1; index >= 0; index-- {
			if r.shutdowns[index] != nil {
				result = errors.Join(result, r.shutdowns[index](ctx))
			}
		}
		if r.closeConnection != nil {
			result = errors.Join(result, r.closeConnection())
		}
		r.shutdownErr = result
	})
	return r.shutdownErr
}

// Start 创建并安装显式选择的全局 OpenTelemetry Provider。
func Start(ctx context.Context, opts ...Option) (*Runtime, error) {
	if ctx == nil {
		return nil, errors.New("telemetry start context is nil")
	}
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	if err := options.validate(); err != nil {
		return nil, err
	}

	resourceAttrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(options.name),
	}
	if len(options.attributes) > 0 {
		resourceAttrs = append(resourceAttrs, options.attributes...)
	}
	res, err := resource.New(ctx,
		resource.WithAttributes(resourceAttrs...),
	)
	if err != nil {
		return nil, fmt.Errorf("create telemetry resource: %w", err)
	}
	grpcOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(grpcServiceConfig),
	}
	if options.insecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})))
	}
	if len(options.options) > 0 {
		grpcOpts = append(grpcOpts, options.options...)
	}
	conn, err := grpc.NewClient(options.endpoint, grpcOpts...)
	if err != nil {
		return nil, fmt.Errorf("create telemetry OTLP connection: %w", err)
	}
	providers, err := initializeSignalProviders(ctx, options.signals, res, conn, []providerInitializer{
		{signal: SignalTracer, initialize: initTracer},
		{signal: SignalMeter, initialize: initMeter},
		{signal: SignalLogger, initialize: initLogger},
	})
	if err != nil {
		return nil, errors.Join(err, conn.Close())
	}
	shutdowns := make([]shutdownFunc, 0, len(providers))
	for _, provider := range providers {
		if provider.install != nil {
			provider.install()
		}
		shutdowns = append(shutdowns, provider.shutdown)
	}
	if options.signals.Enabled(SignalTracer) {
		otel.SetTextMapPropagator(newPropagator())
	}
	return newRuntime(options.signals, shutdowns, conn.Close), nil
}

func initializeSignalProviders(
	ctx context.Context,
	signals Signal,
	res *resource.Resource,
	conn *grpc.ClientConn,
	initializers []providerInitializer,
) ([]signalProvider, error) {
	providers := make([]signalProvider, 0, len(initializers))
	for _, initializer := range initializers {
		if !signals.Enabled(initializer.signal) {
			continue
		}
		provider, err := initializer.initialize(ctx, res, conn)
		if err == nil {
			providers = append(providers, provider)
			continue
		}
		cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), providerCleanupTimeout)
		cleanupErr := shutdownSignalProviders(cleanupCtx, providers)
		cancel()
		return nil, errors.Join(err, cleanupErr)
	}
	return providers, nil
}

func shutdownSignalProviders(ctx context.Context, providers []signalProvider) error {
	var result error
	for index := len(providers) - 1; index >= 0; index-- {
		if providers[index].shutdown != nil {
			result = errors.Join(result, providers[index].shutdown(ctx))
		}
	}
	return result
}

func (o options) validate() error {
	const knownSignals = SignalTracer | SignalLogger | SignalMeter
	if o.signals&^knownSignals != 0 {
		return errors.New("telemetry contains unknown signals")
	}
	if o.signals == 0 {
		return errors.New("telemetry signal is not configured")
	}
	if strings.TrimSpace(o.name) == "" {
		return errors.New("telemetry service name is not configured")
	}
	if strings.TrimSpace(o.endpoint) == "" {
		return errors.New("telemetry endpoint is not configured")
	}
	return nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

// initTracer trace provider
func initTracer(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (signalProvider, error) {
	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return signalProvider{}, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	return signalProvider{
		install: func() { otel.SetTracerProvider(tracerProvider) }, shutdown: tracerProvider.Shutdown,
	}, nil
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
func initMeter(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (signalProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return signalProvider{}, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	return signalProvider{
		install: func() { otel.SetMeterProvider(meterProvider) }, shutdown: meterProvider.Shutdown,
	}, nil
}

// initLogger logger provider
func initLogger(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (signalProvider, error) {
	loggerExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
	if err != nil {
		return signalProvider{}, fmt.Errorf("failed to create logger exporter: %w", err)
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(loggerExporter)),
		sdklog.WithResource(res),
	)
	return signalProvider{
		install: func() { logglobal.SetLoggerProvider(loggerProvider) }, shutdown: loggerProvider.Shutdown,
	}, nil
}
