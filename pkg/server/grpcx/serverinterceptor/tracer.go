// Package serverinterceptor provides ...
package serverinterceptor

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

// OTelTracerHandler alias otelgrpc.ServerHandler
var OTelTracerHandler = otelgrpc.NewServerHandler
