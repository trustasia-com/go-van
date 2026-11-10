// Package handler provides ...
package handler

import (
	"bufio"
	"net"
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

// Unwrap 保留底层响应控制能力，供 http.ResponseController 设置连接期限。
func (w *WrappedWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// Hijack 保留 HTTP/1 连接升级能力；不支持的传输返回标准错误。
func (w *WrappedWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return http.NewResponseController(w.ResponseWriter).Hijack()
}

// FlushError 透传刷新错误，避免流式响应把写入失败报告为成功。
func (w *WrappedWriter) FlushError() error {
	return http.NewResponseController(w.ResponseWriter).Flush()
}

// Flush 保留 http.Flusher 接口；需要诊断错误的调用方使用 FlushError。
func (w *WrappedWriter) Flush() {
	_ = w.FlushError()
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
