// Package main provides ...
package main

import (
	"context"
	"fmt"

	"github.com/trustasia-com/go-van"
	pb "github.com/trustasia-com/go-van/examples/grpc/helloworld"
	"github.com/trustasia-com/go-van/pkg/codes"
	"github.com/trustasia-com/go-van/pkg/codes/status"
	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/grpcx"
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
	// grpc server
	srv := grpcx.NewServer(
		server.WithAddress(":8000"),
	)
	s := &serverGRPC{}
	pb.RegisterGreeterServer(srv, s)

	service := van.NewService(
		van.WithName("grpc"),
		van.WithServer(srv),
	)
	if err := service.Run(); err != nil {
		logx.Fatal(err)
	}
}
