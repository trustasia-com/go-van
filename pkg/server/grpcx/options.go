// Package grpcx provides gRPC utilities for the server package
package grpcx

import (
	"context"

	"github.com/trustasia-com/go-van/pkg/server"

	"google.golang.org/grpc"
)

type grpcOptsKey struct{}

// WithDialOpt grpc client option
func WithDialOpt(options ...grpc.DialOption) server.DialOption {
	return func(opts *server.DialOptions) {
		if opts.Context == nil {
			opts.Context = context.Background()
		}

		// Safely extract existing options from context
		if existingOpts := opts.Context.Value(grpcOptsKey{}).([]grpc.DialOption); len(existingOpts) > 0 {
			options = append(existingOpts, options...)
		}

		opts.Context = context.WithValue(opts.Context, grpcOptsKey{}, options)
	}
}
