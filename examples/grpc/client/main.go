// Package main provides ...
package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/deepzz0/go-van/examples/grpc/helloworld"
	"github.com/deepzz0/go-van/pkg/codes"
	"github.com/deepzz0/go-van/pkg/codes/status"
	"github.com/deepzz0/go-van/pkg/registry"
	"github.com/deepzz0/go-van/pkg/registry/etcd"
	"github.com/deepzz0/go-van/pkg/server"
	"github.com/deepzz0/go-van/pkg/server/grpcx"
)

func main() {
	reg := etcd.NewRegistry(
		registry.WithTTL(time.Second*10),
		registry.WithAddress("localhost:2379"),
	)

	// grpc client
	conn, err := grpcx.DialContext(
		server.WithEndpoint("discov:///grpc"),
		server.WithRegistry(reg),
	)
	if err != nil {
		panic(err)
	}
	cli := pb.NewGreeterClient(conn)
	reply, err := cli.SayHello(context.Background(), &pb.HelloRequest{Name: "go-van"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("[gRPC] SayHello %+v\n", reply)

	// returns error
	_, err = cli.SayHello(context.Background(), &pb.HelloRequest{Name: "error"})
	if err != nil {
		code := status.Code(err)
		if code == codes.InvalidArgument {
			fmt.Println(err.Error())
		}
	}

}
