// Package main provides ...
package main

import (
	"net/http"
	"time"

	"github.com/trustasia-com/go-van/pkg/registry"
	"github.com/trustasia-com/go-van/pkg/registry/etcd"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/httpx"
)

func main() {
	reg := etcd.NewRegistry(
		registry.WithTTL(time.Second*10),
		registry.WithAddress("192.168.252.177:2379"),
	)

	cli := httpx.NewClient(
		server.WithTimeout(time.Second*2),
		server.WithRegistry(reg),
		server.WithEndpoint("http://gin-http"),
	)

	req := &httpx.Request{
		Method: http.MethodGet,
		Path:   "/hello",
	}
	err := cli.Do(req, nil)
	if err != nil {
		panic(err)
	}
	// idle conn timeout
	time.Sleep(time.Second * 91)
	err = cli.Do(req, nil)
	if err != nil {
		panic(err)
	}
	err = cli.Do(req, nil)
	if err != nil {
		panic(err)
	}
}
