// Package resolver provides ...
package resolver

import (
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestDirectResolver(t *testing.T) {
	_, err := grpc.Dial("direct:///localhost:50051,localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
}
