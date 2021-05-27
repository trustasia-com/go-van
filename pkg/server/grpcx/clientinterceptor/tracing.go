// Package clientinterceptor provides ...
package clientinterceptor

import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

// UnaryTraceInterceptor alias otelgrpc.UnaryClientInterceptor
var UnaryTraceInterceptor = otelgrpc.UnaryClientInterceptor

// StreamTraceInterceptor alias otelgrpc.StreamClientInterceptor
var StreamTraceInterceptor = otelgrpc.StreamClientInterceptor
