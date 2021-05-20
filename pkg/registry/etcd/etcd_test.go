// Package etcd provides ...
package etcd

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/deepzz0/go-van/pkg/registry"
)

var reg registry.Registry

func init() {
	reg = NewRegistry(registry.WithEndpoint("127.0.0.1:2379"))
	w, err := reg.Watch(context.Background(), "server1")
	if err != nil {
		panic(err)
	}
	go func() {
		time.Sleep(time.Second * 2)
		w.Stop()
	}()
	go watch(w)
}

func watch(w registry.Watcher) {
	for {
		srvs, err := w.Next()
		if err != nil {
			log.Println("error", err)
			if err.Error() == "watcher was canceled" {
				break
			}
			continue
		}
		log.Println("server1 list", len(srvs))
	}
}

func TestRegistry(t *testing.T) {
	srv1 := &registry.Service{
		ID:   "1",
		Name: "server1",
	}
	srv2 := &registry.Service{
		ID:   "2",
		Name: "server2",
	}
	// register
	err := reg.Register(context.Background(), srv1)
	if err != nil {
		t.Fatal(err)
	}
	err = reg.Register(context.Background(), srv2)
	if err != nil {
		t.Fatal(err)
	}
	// get service
	srvs, err := reg.GetService(context.Background(), srv2.Name)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range srvs {
		t.Log(v.ID)
	}
	// deregister
	err = reg.Deregister(context.Background(), srv1)
	if err != nil {
		t.Fatal(err)
	}
	// sleep 3 seconds wait for watcher stop
	time.Sleep(time.Second * 3)
}
