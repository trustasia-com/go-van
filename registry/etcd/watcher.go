// Package etcd provides ...
package etcd

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/deepzz0/go-van/registry"

	"go.etcd.io/etcd/clientv3"
)

type watcher struct {
	client *clientv3.Client
	watch  clientv3.WatchChan
	stop   chan struct{}
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
		stop:   stop,
		watch:  client.Watch(ctx, key, clientv3.WithPrefix()),
	}
	return w
}

func (w *watcher) Next() (*registry.Result, error) {
	for resp := range w.watch {
		if err := resp.Err(); err != nil {
			return nil, err
		}
		if resp.Canceled {
			break
		}
		for _, ev := range resp.Events {
			result := &registry.Result{}

			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					result.EventType = registry.Create
				} else {
					result.EventType = registry.Update
				}
				srv := &registry.Service{}
				json.Unmarshal(ev.Kv.Value, srv)
				result.Service = srv
			case clientv3.EventTypeDelete:
				srv := &registry.Service{}
				json.Unmarshal(ev.PrevKv.Value, srv)
				result.EventType = registry.Delete
				result.Service = srv
			}
			return result, nil
		}
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
