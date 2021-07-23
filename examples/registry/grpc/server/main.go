// Package main provides ...
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/trustasia-com/go-van"
	pb "github.com/trustasia-com/go-van/examples/grpc/helloworld"
	"github.com/trustasia-com/go-van/pkg/codes"
	"github.com/trustasia-com/go-van/pkg/codes/status"
	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/trustasia-com/go-van/pkg/registry"
	"github.com/trustasia-com/go-van/pkg/registry/etcd"
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
	reg := etcd.NewRegistry(
		registry.WithTTL(time.Second*10),
		registry.WithAddress("192.168.252.177:2379"),
	)

	// grpc server
	port := rand.Intn(999) + 8000
	srv := grpcx.NewServer(
		server.WithAddress(fmt.Sprintf(":%d", port)),
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
