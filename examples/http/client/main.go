// Package main provides ...
package main

import (
	"fmt"
	"net/http"

	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/httpx"
)

func main() {
	cli := httpx.NewClient(
		server.WithEndpoint("https://baidu.com"),
	)

	req := httpx.NewRequest(http.MethodGet, "", nil)
	resp, err := cli.Do(req)
	if err != nil {
		panic(err)
	}
	var result []byte
	resp.Scan(&result)
	fmt.Println(string(result))
}
