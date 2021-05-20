// Package grpc provides ...
package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/deepzz0/go-van/pkg/middleware/recovery"
	"github.com/deepzz0/go-van/pkg/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// NewServer new grpc server
func NewServer(opts ...server.Option) (server.Server, error) {
	opt := server.Options{
		Network:  "tcp",
		Endpoint: ":0",
		Timeout:  time.Second,
		Context:  context.Background(),
	}
	svr := &grpcServer{options: opt}
	// apply option
	for _, o := range opts {
		o(&svr.options)
	}
	// prepare grpc option
	grpcOpts := []grpc.ServerOption{}
	// recover options
	if svr.options.Recover {
		grpcOpts = append(grpcOpts,
			grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor()),
			grpc.ChainStreamInterceptor(recovery.StreamServerInterceptor()),
		)
	}
	// other server option or middleware
	if svr.options.Context != nil {
		opts, ok := svr.options.Context.Value(grpcOptsKey{}).([]grpc.ServerOption)
		if ok {
			grpcOpts = append(grpcOpts, opts...)
		}
	}
	// new grpc server
	svr.Server = grpc.NewServer(grpcOpts...)
	// Register reflection service on gRPC server.
	reflection.Register(svr.Server)
	// grpc health server
	healthpb.RegisterHealthServer(svr.Server, svr.healthSvr)
	return svr, nil
}

// grpcServer grpc server
type grpcServer struct {
	options  server.Options
	grpcOpts []grpc.ServerOption

	*grpc.Server
	healthSvr *health.Server
}

func (s *grpcServer) Start() error {
	lis, err := net.Listen(s.options.Network, s.options.Endpoint)
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
	return fmt.Sprintf("grpc://%s" + s.options.Endpoint), nil
}
