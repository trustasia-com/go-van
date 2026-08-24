// Package handler provides ...
package handler

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

const meterName = "go-van-meter"

// Status returns the HTTP response status code of the current request.
type httpStatusCode interface {
	Status() int
}

// WrappedWriter wrapped http response writer
type WrappedWriter struct {
	statusCode int

	http.ResponseWriter
}

// WriteHeader cover writer WriteHeader
func (w *WrappedWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode

	w.ResponseWriter.WriteHeader(statusCode)
}

// Flush cover writer Flush
func (w *WrappedWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Status 返回已写入的 HTTP Status；未显式写 Header 时为 200。
func (w *WrappedWriter) Status() int {
	if w.statusCode == 0 {
		return http.StatusOK
	}
	return w.statusCode
}

// MeterSrvHandler returns a middleware that metrics the request.
func MeterSrvHandler(next http.Handler) http.Handler {
	meter := otel.Meter(meterName)

	// Tracks the number of HTTP requests
	apiCounter, _ := meter.Int64Counter(
		"api.counter",
		metric.WithDescription("Number of API calls"),
		metric.WithUnit("{call}"),
	)
	// Duration of HTTP requests
	apiDuration, _ := meter.Float64Histogram(
		"api.duration",
		metric.WithDescription("Duration of HTTP requests"),
		metric.WithUnit("s"),
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		now := time.Now()

		// serve the request to the next middleware
		next.ServeHTTP(w, r)

		if r.URL.Path == "/ping" {
			return
		}

		// increment the counter
		status := http.StatusOK
		if writer, ok := w.(httpStatusCode); ok {
			status = writer.Status()
		}
		attributes := metric.WithAttributes(
			metricHTTPMethod(r.Method),
			semconv.HTTPResponseStatusCode(status),
		)
		apiCounter.Add(ctx, 1, attributes)
		dur := time.Since(now).Seconds()
		apiDuration.Record(ctx, dur, attributes)
	})
}

func metricHTTPMethod(method string) attribute.KeyValue {
	switch method {
	case http.MethodConnect, http.MethodDelete, http.MethodGet, http.MethodHead,
		http.MethodOptions, http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace:
		return semconv.HTTPRequestMethodKey.String(method)
	default:
		return semconv.HTTPRequestMethodOther
	}
}
