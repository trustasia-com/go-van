// Package resolver provides ...
package resolver

import (
	"context"
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
	reg := etcd.NewRegistry(registry.WithAddress("127.0.0.1:2379"))
	reg.Register(context.Background(), &registry.Instance{
		ID:        "1",
		Name:      "helloworld",
		Endpoints: []string{"grpc://localhost:50051"},
	})
	resolv = NewBuilder(reg)
}

func TestDiscovResolver(t *testing.T) {
	cc, err := grpc.Dial("discov:///helloworld",
		grpc.WithResolvers(resolv),
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
