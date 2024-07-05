// Package resolver provides ...
package resolver

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

type directBuilder struct{}

// ResolveNow could be called multiple times concurrently
func (d *directBuilder) ResolveNow(options resolver.ResolveNowOptions) {}

// Close close watcher
func (d *directBuilder) Close() {}

// Build implements resolver.Resolver
func (d *directBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {

	var addrs []resolver.Address
	endpoints := strings.FieldsFunc(target.Endpoint(), func(r rune) bool {
		return r == EndpointSepChar
	})

	for _, val := range endpoints {
		addrs = append(addrs, resolver.Address{Addr: val})
	}
	cc.UpdateState(resolver.State{Addresses: subset(addrs, subsetSize)})

	return d, nil
}

// Scheme return scheme
func (d *directBuilder) Scheme() string {
	return DirectScheme
}
