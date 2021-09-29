// Package main provides ...
package main

import (
	"net/http"

	"github.com/trustasia-com/go-van"
	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/httpx"
)

func main() {
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
	)
	if err := service.Run(); err != nil {
		logx.Fatal(err)
	}
}
