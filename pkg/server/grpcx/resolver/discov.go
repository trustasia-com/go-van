// Package resolver provides ...
package resolver

import (
	"context"
	"net/url"
	"time"

	"github.com/trustasia-com/go-van/pkg/registry"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

// NewBuilder create registry resolver
func NewBuilder(reg registry.Registry) resolver.Builder {
	return &discovBuilder{
		registry: reg,
		ch:       make(chan struct{}),
	}
}

type discovBuilder struct {
	registry registry.Registry
	w        registry.Watcher
	cc       resolver.ClientConn
	ch       chan struct{}
}

// ResolveNow could be called multiple times concurrently
func (d *discovBuilder) ResolveNow(opts resolver.ResolveNowOptions) {}

// Close close watcher
func (d *discovBuilder) Close() {
	d.w.Stop()
	close(d.ch)
}

// Scheme return scheme
func (d *discovBuilder) Scheme() string {
	return DiscovScheme
}

// Build implements resolver.Resolver
func (d *discovBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {

	d.cc = cc
	// watch addrs change
	w, err := d.registry.Watch(context.TODO(), target.Endpoint)
	if err != nil {
		return nil, err
	}
	d.w = w
	go d.watch()
	return d, nil
}

// watch watch the registry change
func (d *discovBuilder) watch() {
	for {
		select {
		case <-d.ch:
			return
		default:
		}
		// apply action
		inss, err := d.w.Next()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		_ = d.update(inss)
	}
}

func (d *discovBuilder) update(inss []*registry.Instance) error {
	var addrs []resolver.Address
	for _, ins := range inss {
		for _, e := range ins.Endpoints {
			u, err := url.Parse(e)
			if err != nil {
				return err
			}
			// find grpc endpoint
			if u.Scheme != "grpc" {
				continue
			}
			var pairs []interface{}
			for k, v := range ins.Metadata {
				pairs = append(pairs, k, v)
			}
			addr := resolver.Address{
				ServerName: ins.Name,
				Attributes: attributes.New(pairs...),
				Addr:       u.Host,
			}
			addrs = append(addrs, addr)
		}
	}
	addrs = subset(addrs, subsetSize)
	d.cc.UpdateState(resolver.State{Addresses: addrs})
	return nil
}
