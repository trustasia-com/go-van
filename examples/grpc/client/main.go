// Package main provides ...
package main

import (
	"context"
	"fmt"

	pb "github.com/trustasia-com/go-van/examples/grpc/helloworld"
	"github.com/trustasia-com/go-van/pkg/codes"
	"github.com/trustasia-com/go-van/pkg/codes/status"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/grpcx"
)

func main() {
	// grpc client
	conn, err := grpcx.DialContext(
		server.WithEndpoint("localhost:8000"),
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
