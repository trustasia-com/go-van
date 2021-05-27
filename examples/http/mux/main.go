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

	"github.com/gorilla/mux"
)

func main() {
	reg := etcd.NewRegistry(
		registry.WithTTL(time.Second*10),
		registry.WithAddress("localhost:2379"),
	)
	// httprouter server
	r := mux.NewRouter()
	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("gin middleware"))

			next.ServeHTTP(w, r)
		})
	})
	r.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("panic error")
	})

	srv := httpx.NewServer(
		server.WithAddress(":9000"),
		server.WithHandler(r),
	)
	service := van.NewService(
		van.WithName("mux-http"),
		van.WithServer(srv),
		van.WithRegistry(reg),
	)
	if err := service.Run(); err != nil {
		logx.Fatal(err)
	}

}
