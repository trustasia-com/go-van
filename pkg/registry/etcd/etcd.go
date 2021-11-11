// Package etcd provides ...
package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/trustasia-com/go-van/pkg/registry"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

// prefix store k/v prefix
const prefix = "/go-van/registry"

// NewRegistry return etcd regsitry
func NewRegistry(opts ...registry.Option) registry.Registry {
	options := registry.Options{
		TTL: time.Second * 15,
	}
	// apply option
	for _, o := range opts {
		o(&options)
	}
	if len(options.Addresses) == 0 {
		logx.Error("etcd: not found addresses, please specify")
	}
	// new etcd client
	config := clientv3.Config{
		TLS:         options.TLS,
		DialTimeout: time.Second * 5,
		Endpoints:   options.Addresses,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
		},
	}
	// auth cred
	if options.Context != nil {
		auth, ok := options.Context.Value(authKey{}).(*authCreds)
		if ok {
			config.Username = auth.username
			config.Password = auth.password
		}
	}
	// ignore error, will call handle error
	client, err := clientv3.New(config)
	if err != nil {
		logx.Errorf("etcd: new etcd client: %s", err)
	}
	return &etcdRegistry{options: options, client: client}
}

type etcdRegistry struct {
	options registry.Options

	client  *clientv3.Client
	leaseID clientv3.LeaseID
}

// Register register service to registry
func (r *etcdRegistry) Register(ctx context.Context, ins *registry.Instance) error {
	key := fmt.Sprintf("%s/%s/%s", prefix, ins.Name, ins.ID)

	data, err := json.Marshal(ins)
	if err != nil {
		return err
	}
	// lease id
	resp, err := r.client.Grant(ctx, int64(r.options.TTL.Seconds()))
	if err != nil {
		return err
	}
	_, err = r.client.Put(ctx, key, string(data), clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	r.leaseID = resp.ID
	return r.keepAliveAsync(ctx, ins)
}

// Deregister deregister service from registry
func (r *etcdRegistry) Deregister(ctx context.Context, ins *registry.Instance) error {
	key := fmt.Sprintf("%s/%s/%s", prefix, ins.Name, ins.ID)

	_, err := r.client.Delete(ctx, key)
	return err
}

// GetService get service from regsitry
func (r *etcdRegistry) GetService(ctx context.Context, name string) ([]*registry.Instance, error) {
	key := fmt.Sprintf("%s/%s", prefix, name)
	resp, err := r.client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}

	var items []*registry.Instance
	for _, kv := range resp.Kvs {
		srv := &registry.Instance{}
		err = json.Unmarshal(kv.Value, srv)
		if err != nil {
			return nil, err
		}
		items = append(items, srv)
	}
	return items, nil
}

// Watch service change
func (r *etcdRegistry) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	key := fmt.Sprintf("%s/%s", prefix, name)
	return newWatcher(ctx, key, r.client), nil
}

func (r *etcdRegistry) keepAliveAsync(ctx context.Context, ins *registry.Instance) error {
	ch, err := r.client.KeepAlive(context.TODO(), r.leaseID)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case _, ok := <-ch:
				if !ok {
					r.client.Revoke(ctx, r.leaseID)
					_ = r.Register(ctx, ins)
					return
				}
			}
		}
	}()
	return nil
}
