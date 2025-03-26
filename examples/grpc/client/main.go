// Package main provides ...
package main

import (
	"context"
	"fmt"

	pb "github.com/trustasia-com/go-van/examples/grpc/helloworld"
	"github.com/trustasia-com/go-van/pkg/codes"
	"github.com/trustasia-com/go-van/pkg/codes/status"
	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/grpcx"
)

func main() {
	// grpc client
	conn, err := grpcx.DialContext(
		server.WithEndpoint("localhost:8000"),
		server.WithCliFlag(server.ClientStdFlag, server.FlagInsecure),
	)
	if err != nil {
		panic(err)
	}
	cli := pb.NewGreeterClient(conn)

	for _, name := range []string{"go-van", "error"} {
		reply, err := cli.SayHello(context.Background(), &pb.HelloRequest{
			Name: name,
			Lang: codes.LangEnUS,
		})
		if err != nil {
			status, ok := status.FromError(err)
			if ok {
				fmt.Println(status.Message())
			} else {
				logx.Error("grpc client SayHello error", err)
			}
			continue
		}
		fmt.Printf("[gRPC] SayHello %+v\n", reply)
	}

}
