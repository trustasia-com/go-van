// Package serverinterceptor provides ...
package serverinterceptor

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

// OtelTraceHandler alias otelgrpc.ServerHandler
var OtelTraceHandler = otelgrpc.NewServerHandler
