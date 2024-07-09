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

// TraceSrvHandler returns a middleware that trace the request.
func TraceSrvHandler(next http.Handler) http.Handler {
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

		spanName := r.RequestURI
		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(httpconv.ClientRequest(r)...),
			oteltrace.WithAttributes(httpconv.ServerRequest("", r)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", r.Method)
		}
		ctx, span := tracer.Start(ctx, spanName, opts...)
		defer span.End()

		// pass the span through the request context
		r = r.WithContext(ctx)

		// serve the request to the next middleware
		next.ServeHTTP(w, r)
	})
}

// TraceCliHandler returns a middleware that trace the request.
func TraceCliHandler(trans http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(trans)
}
