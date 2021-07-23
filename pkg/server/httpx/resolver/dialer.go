// Package resolver provides ...
package resolver

import (
	"context"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/trustasia-com/go-van/pkg/registry"
)

var dialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
}

// NewBuilder create registry resolver
func NewBuilder(reg registry.Registry) Builder {
	return &resolveBuilder{
		registry: reg,
		cache:    &cache{},

		ch: make(chan struct{}),
	}
}

type resolveBuilder struct {
	registry registry.Registry
	cache    *cache

	w  registry.Watcher
	d  *net.Dialer
	ch chan struct{}
}

func (d *resolveBuilder) Build() (Dialer, error) {
	return d.DialContext, nil
}

func (d *resolveBuilder) DialContext(ctx context.Context, network,
	addr string) (net.Conn, error) {

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	// whether ip address
	if net.ParseIP(host) != nil || strings.Contains(addr, ".") {
		return dialer.DialContext(ctx, network, addr)
	}
	// discovery name
	ip := d.cache.pick()

	return dialer.DialContext(ctx, network, ip+":"+port)
}

// watch watch the registry change
func (d *resolveBuilder) watch() {
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

func (d *resolveBuilder) update(inss []*registry.Instance) error {
	var addrs []string
	for _, ins := range inss {
		for _, e := range ins.Endpoints {
			u, err := url.Parse(e)
			if err != nil {
				return err
			}
			// find http endpoint
			if u.Scheme != "http" {
				continue
			}
			vals := url.Values{}
			for k, v := range ins.Metadata {
				vals.Set(k, v)
			}
			addrs = append(addrs, u.Host)
		}
	}
	d.cache.updateAddresses(addrs)
	return nil
}

// cache
type cache struct {
	addrs []string
	l     sync.RWMutex
	next  int
}

// updateAddresses update addresses
func (c *cache) updateAddresses(addrs []string) {
	c.l.Lock()
	c.addrs = addrs
	c.l.Unlock()
}

// pick pick address
func (c *cache) pick() string {
	c.l.RLock()
	addr := c.addrs[c.next]
	c.next = (c.next + 1) % len(c.addrs)
	c.l.RUnlock()

	return addr
}
