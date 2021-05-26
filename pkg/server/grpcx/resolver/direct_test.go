// Package resolver provides ...
package resolver

import (
	"context"
	"testing"
	"time"

	pb "github.com/deepzz0/go-van/examples/helloworld"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

func TestDirectResolver(t *testing.T) {
	cc, err := grpc.Dial("direct:///localhost:50051,localhost:50052",
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
	)
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
