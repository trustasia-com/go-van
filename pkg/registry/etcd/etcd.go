// Package etcd provides ...
package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deepzz0/go-van/pkg/registry"

	clientv3 "go.etcd.io/etcd/client/v3"
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
	reg := &etcdRegistry{options: options}
	// new etcd client
	config := clientv3.Config{
		TLS:         reg.options.TLS,
		DialTimeout: time.Second * 5,
		Endpoints:   reg.options.Addresses,
	}
	// auth cred
	if reg.options.Context != nil {
		auth, ok := reg.options.Context.Value(authKey{}).(*authCreds)
		if ok {
			config.Username = auth.username
			config.Password = auth.password
		}
	}
	// ignore error, will call handle error
	reg.client, _ = clientv3.New(config)
	return reg
}

type etcdRegistry struct {
	options registry.Options

	client *clientv3.Client
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
	// keep alive
	ch, err := r.client.KeepAlive(context.TODO(), resp.ID)
	if err != nil {
		return err
	}
	go func() {
		for range ch {
			// heartbeat
		}
	}()
	return nil
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
