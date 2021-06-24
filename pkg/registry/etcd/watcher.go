// Package etcd provides ...
package etcd

import (
	"context"
	"encoding/json"

	"github.com/deepzz0/go-van/pkg/registry"

	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type watcher struct {
	client *clientv3.Client
	watch  clientv3.WatchChan

	stop chan struct{}
	key  string
	ctx  context.Context
}

func newWatcher(ctx context.Context, key string, client *clientv3.Client) *watcher {
	ctx, cancel := context.WithCancel(ctx)
	stop := make(chan struct{}, 1)
	go func() {
		<-stop
		cancel()
	}()

	w := &watcher{
		client: client,
		watch:  client.Watch(ctx, key, clientv3.WithPrefix(), clientv3.WithRev(0)),
		stop:   stop,
		key:    key,
		ctx:    ctx,
	}
	client.RequestProgress(ctx)
	return w
}

func (w *watcher) Next() ([]*registry.Instance, error) {
	for resp := range w.watch {
		if err := resp.Err(); err != nil {
			return nil, err
		}
		if resp.Canceled {
			break
		}

		resp, err := w.client.Get(w.ctx, w.key, clientv3.WithPrefix())
		if err != nil {
			return nil, err
		}
		var items []*registry.Instance
		for _, kv := range resp.Kvs {
			srv := &registry.Instance{}
			err = json.Unmarshal(kv.Value, &srv)
			if err != nil {
				return nil, err
			}
			items = append(items, srv)
		}
		return items, nil
	}
	return nil, errors.New("watcher was canceled")
}

func (w *watcher) Stop() error {
	select {
	case <-w.stop:
	default:
		close(w.stop)
	}
	return nil
}
