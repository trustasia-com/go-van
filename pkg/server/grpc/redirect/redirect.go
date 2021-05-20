// Package redirect provides ...
package redirect

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

type builder struct{}

// NewBuilder creates registry resolver direct
// example:
//   redirect://<authority>/127.0.0.1:9000,127.0.0.2:9000
func NewBuilder() resolver.Builder {
	return &builder{}
}

// ResolveNow could be called multiple times concurrently
func (d *builder) ResolveNow(opts resolver.ResolveNowOptions) {}

// Close close watcher
func (d *builder) Close() {}

// Scheme return scheme
func (d *builder) Scheme() string {
	return "redirect"
}

// Build implements resolver.Resolver
func (d *builder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {

	var addrs []resolver.Address
	for _, addr := range strings.Split(target.Endpoint, ",") {
		addrs = append(addrs, resolver.Address{Addr: addr})
	}
	cc.UpdateState(resolver.State{Addresses: addrs})
	return d, nil
}
