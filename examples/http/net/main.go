// Package main provides ...
package main

import (
	"net/http"
	"time"

	"github.com/deepzz0/go-van"
	"github.com/deepzz0/go-van/pkg/logx"
	"github.com/deepzz0/go-van/pkg/registry"
	"github.com/deepzz0/go-van/pkg/registry/etcd"
	"github.com/deepzz0/go-van/pkg/server"
	"github.com/deepzz0/go-van/pkg/server/httpx"
)

func main() {
	reg := etcd.NewRegistry(
		registry.WithTTL(time.Second*10),
		registry.WithAddress("localhost:2379"),
	)
	// net/http handler
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	http.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("panic error")
	})

	srv := httpx.NewServer(
		server.WithAddress(":9000"),
	)
	service := van.NewService(
		van.WithName("net-http"),
		van.WithServer(srv),
		van.WithRegistry(reg),
	)
	if err := service.Run(); err != nil {
		logx.Fatal(err)
	}
}
