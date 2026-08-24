// Package handler provides ...
package handler

import (
	"fmt"
	"net/http"

	"github.com/trustasia-com/go-van/pkg"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const tracerName = "go-van-tracer"

// SpanNameFormatter 将受控路由信息转换为低基数 Span Name。
type SpanNameFormatter func(*http.Request) string

// DefaultSpanNameFormatter 不读取 URL，避免 Path 参数与 Query 进入 Telemetry。
func DefaultSpanNameFormatter(request *http.Request) string {
	if request == nil || request.Method == "" {
		return "HTTP UNKNOWN"
	}
	return fmt.Sprintf("HTTP %s", request.Method)
}

// TracerSrvHandler returns a middleware that trace the request.
func TracerSrvHandler(next http.Handler, formatter SpanNameFormatter) http.Handler {
	propagators := otel.GetTextMapPropagator()
	tracer := otel.Tracer(
		tracerName,
		oteltrace.WithInstrumentationVersion(pkg.Version),
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		savedCtx := r.Context()
		defer func() {
			r = r.WithContext(savedCtx)
		}()
		ctx := propagators.Extract(savedCtx, propagation.HeaderCarrier(r.Header))

		spanName := DefaultSpanNameFormatter(r)
		if formatter != nil {
			spanName = formatter(r)
			if spanName == "" {
				spanName = DefaultSpanNameFormatter(r)
			}
		}
		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(httpconv.ServerRequest("", r)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}
		ctx, span := tracer.Start(ctx, spanName, opts...)
		defer span.End()

		// pass the span through the request context
		r = r.WithContext(ctx)

		// serve the request to the next middleware
		w.Header().Set("X-Trace-Id", span.SpanContext().TraceID().String())
		next.ServeHTTP(w, r)
	})
}

// TracerCliHandler returns a middleware that trace the request.
func TracerCliHandler(trans http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(trans)
}
