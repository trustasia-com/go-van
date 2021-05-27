// Package resolver provides ...
package resolver

import (
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

func TestDirectResolver(t *testing.T) {
	_, err := grpc.Dial("direct:///localhost:50051,localhost:50052",
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
	)
	if err != nil {
		t.Fatal(err)
	}
}
