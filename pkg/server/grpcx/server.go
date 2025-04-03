// Package grpcx provides ...
package grpcx

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/trustasia-com/go-van/pkg/internal"
	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/grpcx/serverinterceptor"
	"github.com/trustasia-com/go-van/pkg/telemetry"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// NewServer new grpc server
func NewServer(opts ...server.ServerOption) *Server {
	options := server.ServerOptions{
		Network: "tcp",
		Address: ":0",

		Flag: server.FlagRecover,
	}
	// apply option
	for _, o := range opts {
		o(&options)
	}
	svr := &Server{
		network: options.Network,
		address: options.Address,
	}
	// prepare grpc option
	grpcOpts := []grpc.ServerOption{}

	// flag apply options
	if options.Flag&server.FlagRecover > 0 {
		grpcOpts = append(grpcOpts,
			grpc.ChainUnaryInterceptor(serverinterceptor.UnaryServerInterceptor()),
			grpc.ChainStreamInterceptor(serverinterceptor.StreamServerInterceptor()),
		)
	}
	// telemetry
	if len(options.Telemetry) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var flag telemetry.FlagOption
		svr.shutdown, flag = telemetry.InitProvider(ctx, options.Telemetry...)

		if flag&telemetry.FlagMeter > 0 {
			grpcOpts = append(grpcOpts,
				grpc.ChainUnaryInterceptor(serverinterceptor.UnaryMeterInterceptor()),
				grpc.ChainStreamInterceptor(serverinterceptor.StreamMeterInterceptor()),
			)
		}
		if flag&telemetry.FlagTracer > 0 {
			grpcOpts = append(grpcOpts, grpc.StatsHandler(serverinterceptor.OTelTracerHandler()))
		}
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
	network  string
	address  string
	shutdown func()

	*grpc.Server
	healthSvr *health.Server
}

// Start server
func (s *Server) Start() error {
	lis, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}

	s.healthSvr.Resume()
	logx.Infof("[gRPC] server listening on: %s", lis.Addr().String())
	return s.Serve(lis)
}

// Stop server
func (s *Server) Stop() error {
	logx.Info("[gRPC] server stopping")

	// telemetry
	if s.shutdown != nil {
		s.shutdown()
	}
	s.GracefulStop()
	s.healthSvr.Shutdown()
	return nil
}

// Endpoint return a real address to registry endpoint.
// examples:
//
//	grpc://127.0.0.1:9000?isSecure=false
func (s *Server) Endpoint() (string, error) {
	addr, err := internal.Extract(s.address)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpc://%s", addr), nil
}
