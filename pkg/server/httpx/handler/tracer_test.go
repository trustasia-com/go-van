package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestTracerSrvHandlerUsesSafeLowCardinalitySpanName(t *testing.T) {
	testTracerSpanName(t, nil, "HTTP GET")
}

func TestTracerSrvHandlerAllowsReviewedSpanNameFormatter(t *testing.T) {
	testTracerSpanName(t, func(*http.Request) string {
		return "HTTP GET /whois/v1/:domain"
	}, "HTTP GET /whois/v1/:domain")
}

func testTracerSpanName(t *testing.T, formatter SpanNameFormatter, want string) {
	t.Helper()
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	previous := otel.GetTracerProvider()
	otel.SetTracerProvider(provider)
	t.Cleanup(func() { otel.SetTracerProvider(previous) })

	handler := TracerSrvHandler(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	}), formatter)
	request := httptest.NewRequest(http.MethodGet, "/whois/v1/deepzz.cn?token=secret", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("ended span count = %d, want 1", len(spans))
	}
	if got := spans[0].Name(); got != want {
		t.Fatalf("span name = %q, want %q", got, want)
	}
	for _, item := range spans[0].Attributes() {
		value := item.Value.Emit()
		if strings.Contains(value, "deepzz.cn") || strings.Contains(value, "secret") {
			t.Fatalf("span attribute %q leaks request URL value %q", item.Key, value)
		}
	}
}
