// Package grpc provides ...
package grpc

import (
	"time"

	"github.com/deepzz0/go-van/pkg/server"
	"github.com/deepzz0/go-van/pkg/server/grpc/resolver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

// DialContext dial to grpc server
func DialContext(opts ...server.Option) (*grpc.ClientConn, error) {
	options := server.Options{
		Endpoint: ":0",
		Timeout:  time.Second,
	}
	for _, o := range opts {
		o(&options)
	}
	// grpc dial options
	grpcOpts := []grpc.DialOption{
		grpc.WithBalancerName(roundrobin.Name),
	}
	// discovery
	if options.Registry != nil {
		builder := resolver.NewBuilder(options.Registry)
		grpcOpts = append(grpcOpts, grpc.WithResolvers(builder))
	}
	// tls secure
	if !options.Secure {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
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
