// Package grpcx provides ...
package grpcx

import (
	"fmt"
	"net"

	"github.com/trustasia-com/go-van/pkg/internal"
	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/grpcx/serverinterceptor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// NewServer new grpc server
func NewServer(opts ...server.ServerOption) server.Server {
	options := server.ServerOptions{
		Network: "tcp",
		Address: ":0",

		Flag: server.ServerStdFlag,
	}
	// apply option
	for _, o := range opts {
		o(&options)
	}
	svr := &Server{options: options}
	// prepare grpc option
	grpcOpts := []grpc.ServerOption{}

	// flag apply options
	if options.Flag&server.FlagRecover > 0 {
		grpcOpts = append(grpcOpts,
			grpc.ChainUnaryInterceptor(serverinterceptor.UnaryServerInterceptor()),
			grpc.ChainStreamInterceptor(serverinterceptor.StreamServerInterceptor()),
		)
	}
	if options.Flag&server.FlagTracing > 0 {
		grpcOpts = append(grpcOpts,
			grpc.ChainUnaryInterceptor(serverinterceptor.UnaryTraceInterceptor()),
			grpc.ChainStreamInterceptor(serverinterceptor.StreamTraceInterceptor()),
		)
	}

	// other server option or middleware
	if len(options.Options) > 0 {
		grpcOpts = append(grpcOpts, options.Options...)
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

// Server grpc server
type Server struct {
	options  server.ServerOptions
	grpcOpts []grpc.ServerOption

	*grpc.Server
	healthSvr *health.Server
}

// Start server
func (s *Server) Start() error {
	lis, err := net.Listen(s.options.Network, s.options.Address)
	if err != nil {
		return err
	}

	s.healthSvr.Resume()
	logx.Infof("[gRPC] server listening on: %s", lis.Addr().String())
	return s.Serve(lis)
}

// Stop server
func (s *Server) Stop() error {
	s.GracefulStop()
	s.healthSvr.Shutdown()
	logx.Info("[gRPC] server stopping")
	return nil
}

// Endpoint return a real address to registry endpoint.
// examples:
//   grpc://127.0.0.1:9000?isSecure=false
func (s *Server) Endpoint() (string, error) {
	addr, err := internal.Extract(s.options.Address)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpc://%s", addr), nil
}
