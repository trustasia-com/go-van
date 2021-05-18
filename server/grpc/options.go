// Package grpc provides ...
package grpc

import (
	"context"

	"github.com/deepzz0/go-van/server"

	"google.golang.org/grpc"
)

type grpcOptsKey struct{}

// ServerOption grpc option
func ServerOption(opts ...grpc.ServerOption) server.Option {
	return func(opts *server.Options) {
		if opts.Ctx == nil {
			opts.Ctx = context.Background()
		}
		opts.Ctx = context.WithValue(opts.Ctx, grpcOptsKey{}, opts)
	}
}
