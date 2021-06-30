// Package main provides ...
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/trustasia-com/go-van/pkg/codes"
	"github.com/trustasia-com/go-van/pkg/codes/status"
	"github.com/trustasia-com/go-van/pkg/registry"
	"github.com/trustasia-com/go-van/pkg/registry/etcd"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/grpcx"
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

	for _, name := range []string{"go-van", "error", "panic"} {
		reply, err := cli.SayHello(context.Background(), &pb.HelloRequest{Name: name})
		if err != nil {
			code := status.Code(err)
			if code == codes.InvalidArgument {
				fmt.Println("codes.InvalidArgument: ", err.Error())
			} else {
				fmt.Println(err.Error())
			}
			continue
		}
		fmt.Printf("[gRPC] SayHello %+v\n", reply)
	}

}
