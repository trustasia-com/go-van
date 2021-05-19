// Package grpc provides ...
package grpc

import (
	"context"

	"github.com/deepzz0/go-van/server"

	"google.golang.org/grpc"
)

type grpcOptsKey struct{}

// WithServerOpt grpc server option
func WithServerOpt(opts ...grpc.ServerOption) server.Option {
	return func(opts *server.Options) {
		if opts.Context == nil {
			opts.Context = context.Background()
		}
		opts.Context = context.WithValue(opts.Context, grpcOptsKey{}, opts)
	}
}

// WithDialOpt grpc client option
func WithDialOpt(opts ...grpc.DialOption) server.Option {
	return func(opts *server.Options) {
		if opts.Context == nil {
			opts.Context = context.Background()
		}
		opts.Context = context.WithValue(opts.Context, grpcOptsKey{}, opts)
	}
}
