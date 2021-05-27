// Package main provides ...
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/deepzz0/go-van"
	pb "github.com/deepzz0/go-van/examples/grpc/helloworld"
	"github.com/deepzz0/go-van/pkg/codes"
	"github.com/deepzz0/go-van/pkg/codes/status"
	"github.com/deepzz0/go-van/pkg/logx"
	"github.com/deepzz0/go-van/pkg/registry"
	"github.com/deepzz0/go-van/pkg/registry/etcd"
	"github.com/deepzz0/go-van/pkg/server"
	"github.com/deepzz0/go-van/pkg/server/grpcx"
)

// serverGRPC is used to implement helloworld.GreeterServer.
type serverGRPC struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *serverGRPC) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if in.Name == "error" {
		return nil, status.Err(codes.InvalidArgument, "name invalid")
	}
	if in.Name == "panic" {
		panic("panic error")
	}
	return &pb.HelloReply{Message: fmt.Sprintf("Welcome %+v!", in.Name)}, nil
}

func main() {
	reg := etcd.NewRegistry(
		registry.WithTTL(time.Second*10),
		registry.WithAddress("localhost:2379"),
	)
	// grpc server
	srv := grpcx.NewServer(
		server.WithRecover(true),
		server.WithAddress(":8000"),
	)
	s := &serverGRPC{}
	pb.RegisterGreeterServer(srv, s)

	service := van.NewService(
		van.WithName("grpc"),
		van.WithServer(srv),
		van.WithRegistry(reg),
	)
	if err := service.Run(); err != nil {
		logx.Fatal(err)
	}
}
