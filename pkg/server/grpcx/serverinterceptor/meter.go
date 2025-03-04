// Package serverinterceptor provides ...
package serverinterceptor

import (
	"context"
	"time"

	"github.com/trustasia-com/go-van/pkg/codes/status"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
)

const meterName = "go-van-meter"

// UnaryMeterInterceptor returns a new unary server interceptor for metrics.
func UnaryMeterInterceptor() grpc.UnaryServerInterceptor {
	meter := otel.Meter(meterName)

	// Tracks the number of gRPC requests
	rpcCounter, _ := meter.Int64Counter(
		"rpc.counter",
		metric.WithDescription("Number of gRPC API calls"),
		metric.WithUnit("{call}"),
	)
	// Duration of gRPC requests
	rpcDuration, _ := meter.Float64Histogram(
		"rpc.duration",
		metric.WithDescription("Duration of gRPC requests"),
		metric.WithUnit("s"),
	)
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

		now := time.Now()

		// exec gRPC handler
		resp, err = handler(ctx, req)

		// parse gRPC status
		st, _ := status.FromError(err)

		rpcCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("rpc.method", info.FullMethod),
				attribute.String("rpc.status_code", st.Code().String()),
			),
		)
		dur := time.Since(now).Seconds()
		rpcDuration.Record(ctx, dur,
			metric.WithAttributes(
				attribute.String("rpc.method", info.FullMethod),
				attribute.String("rpc.status_code", st.Code().String()),
			),
		)
		return
	}
}

// StreamMeterInterceptor returns a new streaming server interceptor for metrics.
func StreamMeterInterceptor() grpc.StreamServerInterceptor {
	meter := otel.Meter(meterName)

	// Tracks the number of gRPC requests
	rpcCounter, _ := meter.Int64Counter(
		"rpc.counter",
		metric.WithDescription("Number of gRPC API calls"),
		metric.WithUnit("{call}"),
	)
	// Duration of gRPC requests
	rpcDuration, _ := meter.Float64Histogram(
		"rpc.duration",
		metric.WithDescription("Duration of gRPC requests"),
		metric.WithUnit("s"),
	)
	return func(srv any, stream grpc.ServerStream,
		info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {

		now := time.Now()
		ctx := stream.Context()

		// exec gRPC handler
		err = handler(srv, stream)

		// parse gRPC status
		st, _ := status.FromError(err)

		rpcCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("rpc.method", info.FullMethod),
				attribute.String("rpc.status_code", st.Code().String()),
			),
		)
		dur := time.Since(now).Seconds()
		rpcDuration.Record(ctx, dur,
			metric.WithAttributes(
				attribute.String("rpc.method", info.FullMethod),
				attribute.String("rpc.status_code", st.Code().String()),
			),
		)
		return
	}
}
