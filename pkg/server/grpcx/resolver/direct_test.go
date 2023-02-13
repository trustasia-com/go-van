// Package resolver provides ...
package resolver

import (
	"testing"

	"google.golang.org/grpc"
)

func TestDirectResolver(t *testing.T) {
	_, err := grpc.Dial("direct:///localhost:50051,localhost:50052",
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatal(err)
	}
}
