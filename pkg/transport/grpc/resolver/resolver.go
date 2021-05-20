// Package resolver provides ...
package resolver

import (
	"context"
	"net/url"
	"time"

	"github.com/deepzz0/go-van/registry"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

// NewBuilder create registry resolver
func NewBuilder(reg registry.Registry) resolver.Builder {
	return &builder{
		registry: reg,
		ch:       make(chan struct{}),
	}
}

type builder struct {
	registry registry.Registry
	w        registry.Watcher
	cc       resolver.ClientConn
	ch       chan struct{}
}

// ResolveNow could be called multiple times concurrently
func (d *builder) ResolveNow(opts resolver.ResolveNowOptions) {}

// Close close watcher
func (d *builder) Close() {
	close(d.ch)
	d.w.Stop()
}

// Scheme return scheme
func (d *builder) Scheme() string {
	return "discovery"
}

// Build implements resolver.Resolver
func (d *builder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {

	d.cc = cc
	ctx := context.Background()
	// init addrs
	srvs, err := d.registry.GetService(ctx, target.Endpoint)
	if err != nil {
		return nil, err
	}
	err = d.update(srvs)
	if err != nil {
		return nil, err
	}
	// watch addrs change
	w, err := d.registry.Watch(ctx, target.Endpoint)
	if err != nil {
		return nil, err
	}
	d.w = w
	go d.watch()
	return d, nil
}

// watch watch the registry change
func (d *builder) watch() {
	for {
		select {
		case <-d.ch:
			return
		default:
		}
		// apply action
		srvs, err := d.w.Next()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		_ = d.update(srvs)
	}
}

func (d *builder) update(srvs []*registry.Service) error {
	var addrs []resolver.Address
	for _, srv := range srvs {
		for _, e := range srv.Endpoints {
			u, err := url.Parse(e)
			if err != nil {
				return err
			}
			// find grpc endpoint
			if u.Scheme != "grpc" {
				continue
			}
			var pairs []interface{}
			for k, v := range srv.Metadata {
				pairs = append(pairs, k, v)
			}
			addr := resolver.Address{
				ServerName: srv.Name,
				Attributes: attributes.New(pairs...),
				Addr:       u.Host,
			}
			addrs = append(addrs, addr)
		}
	}
	d.cc.UpdateState(resolver.State{Addresses: addrs})
	return nil
}
