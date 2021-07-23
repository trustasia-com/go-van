// Package httpx provides ...
package httpx

import (
	"net/http"
	"testing"
	"time"

	"github.com/trustasia-com/go-van/pkg/registry"
	"github.com/trustasia-com/go-van/pkg/registry/etcd"
	"github.com/trustasia-com/go-van/pkg/server"
)

const etcdEndpoit = "192.168.252.177:2379"

func TestClientDo(t *testing.T) {
	reg := etcd.NewRegistry(
		registry.WithAddress(etcdEndpoit),
	)

	cli := NewClient(
		server.WithTimeout(time.Second*2),
		server.WithRegistry(reg),
	)

	req, err := http.NewRequest(http.MethodGet, "https://baidu.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	var html string
	err = cli.Do(req, &html)
	if err != nil {
		panic(err)
	}
	t.Log(html)
}
