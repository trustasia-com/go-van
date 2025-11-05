// Package main provides ...
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/httpx"
)

func main() {
	cli := httpx.NewClient(
		server.WithEndpoint("https://api.thecatapi.com/v1"),
		httpx.WithHeader(map[string]string{
			"lang":         "zh-CN",
			"Content-Type": "application/json",
		}),
	)

	req := httpx.NewRequest(http.MethodGet, "/images/search", "", nil)
	resp, err := cli.Do(context.Background(), req)
	if err != nil {
		panic(err)
	}
	var result []map[string]any
	resp.Scan(&result)
	fmt.Printf("%#v\n", result)
}
