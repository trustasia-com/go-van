package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/telemetry"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type testTelemetryRuntime struct {
	signals telemetry.Signal
}

func (runtime testTelemetryRuntime) Enabled(signal telemetry.Signal) bool {
	return runtime.signals.Enabled(signal)
}

func TestServerAutomaticallyInstallsTracingFromRuntime(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	previous := otel.GetTracerProvider()
	otel.SetTracerProvider(provider)
	t.Cleanup(func() { otel.SetTracerProvider(previous) })

	httpServer := NewServer(
		server.WithHandler(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(http.StatusNoContent)
		})),
		server.WithTelemetry(testTelemetryRuntime{signals: telemetry.SignalTracer}),
	)
	httpServer.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest(http.MethodGet, "/whois/v1/deepzz.cn?token=secret", nil),
	)

	if got := len(recorder.Ended()); got != 1 {
		t.Fatalf("ended span count = %d, want 1", got)
	}
}

func TestServerDoesNotInstallTelemetryWithoutRuntime(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	previous := otel.GetTracerProvider()
	otel.SetTracerProvider(provider)
	t.Cleanup(func() { otel.SetTracerProvider(previous) })

	var runtime *telemetry.Runtime
	httpServer := NewServer(
		server.WithHandler(http.HandlerFunc(
			func(writer http.ResponseWriter, _ *http.Request) { writer.WriteHeader(http.StatusNoContent) },
		)),
		server.WithTelemetry(runtime),
	)
	response := httptest.NewRecorder()
	httpServer.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/health", nil))

	if got := len(recorder.Ended()); got != 0 {
		t.Fatalf("ended span count = %d, want 0", got)
	}
	if got := response.Header().Get("X-Trace-Id"); got != "" {
		t.Fatalf("X-Trace-Id = %q, want empty", got)
	}
}
