// Package clientinterceptor provides ...
package clientinterceptor

import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

// OTelTracerHandler alias otelgrpc.ClientHandler
var OTelTracerHandler = otelgrpc.NewClientHandler
