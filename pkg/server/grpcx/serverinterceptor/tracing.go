// Package serverinterceptor provides ...
package serverinterceptor

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

// UnaryTraceInterceptor alias otelgrpc.UnaryServerInterceptor
var UnaryTraceInterceptor = otelgrpc.UnaryServerInterceptor

// StreamTraceInterceptor  alias otelgrpc.StreamServerInterceptor
var StreamTraceInterceptor = otelgrpc.StreamServerInterceptor
