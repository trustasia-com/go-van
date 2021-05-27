// Package handler provides ...
package handler

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// TraceSrvHandler returns a middleware that trace the request.
func TraceSrvHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO

		next.ServeHTTP(w, r)
	})
}

// TraceCliHandler returns a middleware that trace the request.
func TraceCliHandler(trans http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(trans)
}
