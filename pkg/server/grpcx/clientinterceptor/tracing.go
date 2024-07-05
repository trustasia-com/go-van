// Package clientinterceptor provides ...
package clientinterceptor

import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

// OtelTraceHandler alias otelgrpc.ClientHandler
var OtelTraceHandler = otelgrpc.NewClientHandler
