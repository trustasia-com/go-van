// Package grpcx provides ...
package grpcx

import (
	"fmt"
	"net"

	"github.com/deepzz0/go-van/pkg/internal"
	"github.com/deepzz0/go-van/pkg/logx"
	"github.com/deepzz0/go-van/pkg/server"
	"github.com/deepzz0/go-van/pkg/server/grpcx/serverinterceptor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// NewServer new grpc server
func NewServer(opts ...server.ServerOption) *grpcServer {
	opt := server.ServerOptions{
		Network: "tcp",
		Address: ":0",

		Flag: server.ServerStdFlag,
	}
	svr := &grpcServer{options: opt}
	// apply option
	for _, o := range opts {
		o(&svr.options)
	}
	// prepare grpc option
	grpcOpts := []grpc.ServerOption{}

	// flag apply options
	if svr.options.Flag&server.FlagRecover > 0 {
		grpcOpts = append(grpcOpts,
			grpc.ChainUnaryInterceptor(serverinterceptor.UnaryServerInterceptor()),
			grpc.ChainStreamInterceptor(serverinterceptor.StreamServerInterceptor()),
		)
	}
	if svr.options.Flag&server.FlagTracing > 0 {
		grpcOpts = append(grpcOpts,
			grpc.ChainUnaryInterceptor(serverinterceptor.UnaryTraceInterceptor()),
			grpc.ChainStreamInterceptor(serverinterceptor.StreamTraceInterceptor()),
		)
	}

	// other server option or middleware
	if len(svr.options.Options) > 0 {
		grpcOpts = append(grpcOpts, svr.options.Options...)
	}
	// new grpc server
	svr.Server = grpc.NewServer(grpcOpts...)
	svr.healthSvr = health.NewServer()
	// Register reflection service on gRPC server.
	reflection.Register(svr.Server)
	// grpc health server
	healthpb.RegisterHealthServer(svr.Server, svr.healthSvr)
	return svr
}

// grpcServer grpc server
type grpcServer struct {
	options  server.ServerOptions
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
	logx.Infof("[gRPC] server listening on: %s", lis.Addr().String())
	return s.Serve(lis)
}

func (s *grpcServer) Stop() error {
	s.GracefulStop()
	s.healthSvr.Shutdown()
	logx.Info("[gRPC] server stopping")
	return nil
}

// Endpoint return a real address to registry endpoint.
// examples:
//   grpc://127.0.0.1:9000?isSecure=false
func (s *grpcServer) Endpoint() (string, error) {
	addr, err := internal.Extract(s.options.Address)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpc://%s", addr), nil
}
