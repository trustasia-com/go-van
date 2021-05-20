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
	opt := registry.Options{TTL: time.Second * 15}
	reg := &etcdRegistry{options: opt}
	// apply option
	for _, o := range opts {
		o(&reg.options)
	}

	// new etcd client
	config := clientv3.Config{
		TLS:         reg.options.TLSConfig,
		DialTimeout: time.Second * 5,
		Endpoints:   reg.options.Enpoints,
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
	client, _ := clientv3.New(config)
	reg.client = client
	return reg
}

type etcdRegistry struct {
	options registry.Options

	client *clientv3.Client
	lease  clientv3.Lease
}

// Register register service to registry
func (r *etcdRegistry) Register(ctx context.Context, srv *registry.Service) error {
	key := fmt.Sprintf("%s/%s/%s", prefix, srv.Name, srv.ID)
	data, err := json.Marshal(srv)
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
		}
	}()
	return nil
}

// Deregister deregister service from registry
func (r *etcdRegistry) Deregister(ctx context.Context, srv *registry.Service) error {
	key := fmt.Sprintf("%s/%s/%s", prefix, srv.Name, srv.ID)
	_, err := r.client.Delete(ctx, key)
	return err
}

// GetService get service from regsitry
func (r *etcdRegistry) GetService(ctx context.Context, srvName string) ([]*registry.Service, error) {
	key := fmt.Sprintf("%s/%s", prefix, srvName)
	resp, err := r.client.Get(ctx, key, clientv3.WithPrefix(),
		clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}
	var items []*registry.Service
	for _, kv := range resp.Kvs {
		srv := &registry.Service{}
		err = json.Unmarshal(kv.Value, srv)
		if err != nil {
			return nil, err
		}
		items = append(items, srv)
	}
	return items, nil
}

// Watch service change
func (r *etcdRegistry) Watch(ctx context.Context, srvName string) (registry.Watcher, error) {
	key := fmt.Sprintf("%s/%s", prefix, srvName)
	return newWatcher(ctx, key, r.client), nil
}
