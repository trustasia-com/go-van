// Package main provides ...
package main

import (
	"context"
	"fmt"
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

	req := httpx.NewRequest(http.MethodGet, "/hello", "", nil)
	resp, err := cli.Do(context.Background(), req)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resp.Data))
	// idle conn timeout
	time.Sleep(time.Second * 91)
	_, err = cli.Do(context.Background(), req)
	if err != nil {
		panic(err)
	}
	_, err = cli.Do(context.Background(), req)
	if err != nil {
		panic(err)
	}
}
