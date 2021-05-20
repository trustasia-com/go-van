// Package resolver provides ...
package resolver

import (
	"context"
	"fmt"
	"net"

	pb "github.com/deepzz0/go-van/examples/helloworld"

	"google.golang.org/grpc"
)

func init() {
	// grpc server
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			panic(err)
		}
		s := grpc.NewServer()
		pb.RegisterGreeterServer(s, &server{})
		if err = s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Welcome %+v!", in.Name)}, nil
}
