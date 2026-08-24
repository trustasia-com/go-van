package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestMeterSrvHandlerRecordsGETWithoutSensitiveURLAttributes(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	previous := otel.GetMeterProvider()
	otel.SetMeterProvider(provider)
	t.Cleanup(func() {
		otel.SetMeterProvider(previous)
		_ = provider.Shutdown(context.Background())
	})

	handler := MeterSrvHandler(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	}))
	response := httptest.NewRecorder()
	handler.ServeHTTP(
		&WrappedWriter{ResponseWriter: response},
		httptest.NewRequest(http.MethodGet, "/whois/v1/deepzz.cn?token=secret", nil),
	)

	var resourceMetrics metricdata.ResourceMetrics
	if err := reader.Collect(context.Background(), &resourceMetrics); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	var names []string
	for _, scope := range resourceMetrics.ScopeMetrics {
		for _, candidate := range scope.Metrics {
			names = append(names, candidate.Name)
			assertMetricDataHasNoSensitiveURL(t, candidate.Data)
		}
	}
	sort.Strings(names)
	want := []string{"api.counter", "api.duration"}
	if len(names) != len(want) || names[0] != want[0] || names[1] != want[1] {
		t.Fatalf("metric names = %v, want %v", names, want)
	}
}

func TestMetricHTTPMethodNormalizesUnknownValues(t *testing.T) {
	if got := metricHTTPMethod("CUSTOM-deepzz.cn").Value.Emit(); got != "_OTHER" {
		t.Fatalf("metricHTTPMethod() = %q, want _OTHER", got)
	}
}

func assertMetricDataHasNoSensitiveURL(t *testing.T, data metricdata.Aggregation) {
	t.Helper()
	var attributes []attribute.KeyValue
	switch value := data.(type) {
	case metricdata.Sum[int64]:
		for _, point := range value.DataPoints {
			attributes = append(attributes, point.Attributes.ToSlice()...)
		}
	case metricdata.Histogram[float64]:
		for _, point := range value.DataPoints {
			attributes = append(attributes, point.Attributes.ToSlice()...)
		}
	default:
		t.Fatalf("unexpected metric aggregation %T", data)
	}
	values := make(map[string]string, len(attributes))
	for _, item := range attributes {
		value := item.Value.Emit()
		values[string(item.Key)] = value
		if strings.Contains(value, "deepzz.cn") || strings.Contains(value, "secret") {
			t.Fatalf("metric attribute leaks request URL value %q", value)
		}
	}
	if values["http.request.method"] != http.MethodGet {
		t.Fatalf("http.request.method = %q, want GET", values["http.request.method"])
	}
	if values["http.response.status_code"] != "204" {
		t.Fatalf("http.response.status_code = %q, want 204", values["http.response.status_code"])
	}
	if _, exists := values["http.target"]; exists {
		t.Fatal("http.target must not be exported")
	}
}
