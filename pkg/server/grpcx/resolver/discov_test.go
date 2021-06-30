// Package resolver provides ...
package resolver

import (
	"context"
	"testing"

	"github.com/trustasia-com/go-van/pkg/registry"
	"github.com/trustasia-com/go-van/pkg/registry/etcd"

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
	_, err := grpc.Dial("discov:///helloworld",
		grpc.WithResolvers(resolv),
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
	)
	if err != nil {
		t.Fatal(err)
	}
}
