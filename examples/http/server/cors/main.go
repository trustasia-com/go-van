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
		w.Header().Add("X-WeKey-Token", "asfdjasklfjqwi12131")
		w.Header().Add("X-WeKey-Key", "asfdjasklfjqwi12131")
		w.Write([]byte("hello world"))
	})
	http.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("panic error")
	})

	corsOpt := httpx.CORSOptions{
		AllowedOrigins: []string{"https://*.example.com"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}
	corsOpt.AllowedHeaders = []string{"X-WeKey-Token"}
	srv := httpx.NewServer(
		server.WithAddress(":9000"),
		httpx.WithCORS(corsOpt),
	)
	service := van.NewService(
		van.WithName("net-http"),
		van.WithServer(srv),
	)
	if err := service.Run(); err != nil {
		logx.Fatal(err)
	}
}
