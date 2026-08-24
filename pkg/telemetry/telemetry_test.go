package telemetry

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
)

func TestStartRequiresExplicitSignal(t *testing.T) {
	runtime, err := Start(
		context.Background(),
		WithName("test-service"),
		WithEndpoint("127.0.0.1:4317"),
		WithInsecure(),
	)

	if err == nil || !strings.Contains(err.Error(), "signal") {
		t.Fatalf("Start() error = %v, want signal validation error", err)
	}
	if runtime != nil {
		t.Fatal("Start() returned a Runtime for invalid signal configuration")
	}
}

func TestStartUsesTLSByDefault(t *testing.T) {
	previousTracer := otel.GetTracerProvider()
	previousPropagator := otel.GetTextMapPropagator()
	t.Cleanup(func() {
		otel.SetTracerProvider(previousTracer)
		otel.SetTextMapPropagator(previousPropagator)
	})

	runtime, err := Start(
		context.Background(),
		WithName("test-service"),
		WithEndpoint("collector.internal:4317"),
		WithSignals(SignalTracer),
	)
	if err != nil {
		t.Fatalf("Start() error = %v, want secure OTLP client", err)
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := runtime.Shutdown(shutdownCtx); err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}
}

func TestRuntimeShutdownIsIdempotentAndAggregatesErrors(t *testing.T) {
	firstErr := errors.New("first shutdown")
	secondErr := errors.New("second shutdown")
	var calls []string
	runtime := newRuntime(SignalMeter, []shutdownFunc{
		func(context.Context) error {
			calls = append(calls, "tracer")
			return firstErr
		},
		func(context.Context) error {
			calls = append(calls, "meter")
			return secondErr
		},
	}, func() error {
		calls = append(calls, "connection")
		return nil
	})

	err := runtime.Shutdown(context.Background())
	if !errors.Is(err, firstErr) || !errors.Is(err, secondErr) {
		t.Fatalf("Shutdown() error = %v, want both shutdown errors", err)
	}
	wantCalls := []string{"meter", "tracer", "connection"}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf("Shutdown() calls = %v, want %v", calls, wantCalls)
	}

	err = runtime.Shutdown(context.Background())
	if !errors.Is(err, firstErr) || !errors.Is(err, secondErr) {
		t.Fatalf("second Shutdown() error = %v, want cached shutdown errors", err)
	}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf("second Shutdown() calls = %v, want no additional calls", calls)
	}
	if !runtime.Enabled(SignalMeter) || runtime.Enabled(SignalTracer) {
		t.Fatalf("Runtime Signal state does not match configured Meter-only Runtime")
	}
}

func TestOptionsSeparateSignalsFromTransport(t *testing.T) {
	configured := options{}
	WithSignals(SignalMeter, SignalTracer)(&configured)
	WithInsecure()(&configured)

	if !configured.signals.Enabled(SignalMeter) || !configured.signals.Enabled(SignalTracer) {
		t.Fatalf("configured Signals = %v, want Meter and Tracer", configured.signals)
	}
	if configured.signals.Enabled(SignalLogger) {
		t.Fatalf("configured Signals = %v, Logger must remain disabled", configured.signals)
	}
	if !configured.insecure {
		t.Fatal("WithInsecure() did not enable insecure transport")
	}
}

func TestInitializeSignalProvidersCleansUpBeforeReturningInitializationError(t *testing.T) {
	wantErr := errors.New("initialize meter")
	var calls []string
	initializers := []providerInitializer{
		{
			signal: SignalTracer,
			initialize: func(context.Context, *resource.Resource, *grpc.ClientConn) (signalProvider, error) {
				return signalProvider{shutdown: func(context.Context) error {
					calls = append(calls, "tracer-shutdown")
					return nil
				}}, nil
			},
		},
		{
			signal: SignalMeter,
			initialize: func(context.Context, *resource.Resource, *grpc.ClientConn) (signalProvider, error) {
				return signalProvider{}, wantErr
			},
		},
	}

	providers, err := initializeSignalProviders(
		context.Background(), SignalTracer|SignalMeter, nil, nil, initializers,
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf("initializeSignalProviders() error = %v, want %v", err, wantErr)
	}
	if providers != nil {
		t.Fatalf("initializeSignalProviders() providers = %v, want nil", providers)
	}
	if len(calls) != 1 || calls[0] != "tracer-shutdown" {
		t.Fatalf("cleanup calls = %v, want [tracer-shutdown]", calls)
	}
}
