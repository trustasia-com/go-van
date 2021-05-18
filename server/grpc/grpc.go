// Package grpc provides ...
package grpc

import (
	"fmt"
	"net"
	"time"

	"github.com/deepzz0/go-van/server"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// NewServer new grpc server
func NewServer(opts ...server.Option) server.Server {
	opt := server.Options{
		Network: "tcp",
		Address: ":0",
		Timeout: time.Second,
	}
	svr := &grpcServer{options: opt}
	// apply option
	for _, o := range opts {
		o(&svr.options)
	}
	// prepare grpc option
	var grpcOpts []grpc.ServerOption
	if svr.options.Ctx != nil {
		opts, ok := svr.options.Ctx.Value(grpcOptsKey{}).([]grpc.ServerOption)
		if ok {
			grpcOpts = opts
		}
	}
	// opentelemetry tracer
	if svr.options.Trace {
		grpcOpts = append(grpcOpts,
			grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
			grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		)
	}
	// new grpc server
	svr.Server = grpc.NewServer(grpcOpts...)
	// Register reflection service on gRPC server.
	reflection.Register(svr.Server)
	// grpc health server
	healthpb.RegisterHealthServer(svr.Server, svr.healthSvr)
	return svr
}

// grpcServer grpc server
type grpcServer struct {
	options  server.Options
	grpcOpts []grpc.ServerOption

	*grpc.Server
	healthSvr *health.Server
}

func (s *grpcServer) Start() error {
	lis, err := net.Listen(s.options.Network, s.options.Address)
	if err != nil {
		return err
	}

	s.healthSvr.Resume()
	return s.Serve(lis)
}

func (s *grpcServer) Stop() error {
	s.GracefulStop()
	s.healthSvr.Shutdown()
	return nil
}

// Endpoint return a real address to registry endpoint.
// examples:
//   grpc://127.0.0.1:9000?isSecure=false
func (s *grpcServer) Endpoint() (string, error) {
	return fmt.Sprintf("grpc://%s" + s.options.Address), nil
}
