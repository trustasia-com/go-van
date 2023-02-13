// Package grpcx provides ...
package grpcx

import (
	"time"

	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/grpcx/clientinterceptor"
	"github.com/trustasia-com/go-van/pkg/server/grpcx/resolver"

	"google.golang.org/grpc"
)

const grpcServiceConfig = `{"loadBalancingPolicy":"round_robin"}`

// DialContext dial to grpc server
func DialContext(opts ...server.DialOption) (*grpc.ClientConn, error) {
	options := server.DialOptions{
		Endpoint: ":0",
		Timeout:  time.Second,

		Flag: server.ClientStdFlag,
	}
	for _, o := range opts {
		o(&options)
	}
	// default config
	// grpc dial options
	grpcOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(grpcServiceConfig),
	}
	// discovery
	if options.Registry != nil {
		builder := resolver.NewBuilder(options.Registry)
		grpcOpts = append(grpcOpts, grpc.WithResolvers(builder))
	}

	// flag apply option
	if options.Flag&server.FlagSecure == 0 {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}
	if options.Flag&server.FlagTracing > 0 {
		grpcOpts = append(grpcOpts,
			grpc.WithUnaryInterceptor(clientinterceptor.UnaryTraceInterceptor()),
			grpc.WithStreamInterceptor(clientinterceptor.StreamTraceInterceptor()),
		)
	}

	// context custom options
	if options.Context != nil {
		opts, ok := options.Context.Value(grpcOptsKey{}).([]grpc.DialOption)
		if ok {
			grpcOpts = append(grpcOpts, opts...)
		}
		return grpc.DialContext(options.Context, options.Endpoint, grpcOpts...)
	}
	return grpc.Dial(options.Endpoint, grpcOpts...)
}
