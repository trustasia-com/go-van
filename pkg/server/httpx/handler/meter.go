// Package handler provides ...
package handler

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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

		// increment the counter
		status := -1
		if writer, ok := w.(*WrappedWriter); ok {
			status = writer.statusCode
		}
		apiCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("http.method", r.Method),
			attribute.Int("http.status_code", status),
			attribute.String("http.target", r.URL.Path),
		))
		dur := time.Since(now).Seconds()
		apiDuration.Record(ctx, dur, metric.WithAttributes(
			attribute.String("http.method", r.Method),
			attribute.Int("http.status_code", status),
			attribute.String("http.target", r.URL.Path),
		))
	})
}
