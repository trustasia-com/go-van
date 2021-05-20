// Package resolver provides ...
package resolver

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	pb "github.com/deepzz0/go-van/examples/helloworld"
	"github.com/deepzz0/go-van/pkg/registry"
	"github.com/deepzz0/go-van/pkg/registry/etcd"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
)

var resolv resolver.Builder

func init() {
	reg := etcd.NewRegistry(registry.Addr("127.0.0.1:2379"))
	reg.Register(context.Background(), &registry.Service{
		ID:        "1",
		Name:      "helloworld",
		Endpoints: []string{"grpc://localhost:50051"},
	})
	resolv = NewBuilder(reg)

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

func TestResolver(t *testing.T) {
	cc, err := grpc.Dial("discovery:///helloworld",
		grpc.WithResolvers(resolv),
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
	)
	t.Log(cc.GetState().String())
	if err != nil {
		t.Fatal(err)
	}
	client := pb.NewGreeterClient(cc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := client.SayHello(ctx, &pb.HelloRequest{
		Name: "go-van",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", resp)
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Welcome %+v!", in.Name)}, nil
}
