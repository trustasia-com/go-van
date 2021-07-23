// Package main provides ...
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/trustasia-com/go-van"
	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/trustasia-com/go-van/pkg/registry"
	"github.com/trustasia-com/go-van/pkg/registry/etcd"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/httpx"

	"github.com/gin-gonic/gin"
)

func main() {
	reg := etcd.NewRegistry(
		registry.WithTTL(time.Second*10),
		registry.WithAddress("192.168.252.177:2379"),
	)

	// gin server
	e := gin.New()
	e.GET("/hello", func(c *gin.Context) {
		c.String(200, "hello world")
	})
	e.Use(func(c *gin.Context) {
		c.Writer.WriteString("gin middleware")
	})
	e.GET("/panic", func(c *gin.Context) {
		panic("panic error")
	})

	port := rand.Intn(999) + 9000
	srv := httpx.NewServer(
		server.WithAddress(fmt.Sprintf(":%d", port)),
		server.WithHandler(e),
	)
	service := van.NewService(
		van.WithRegistry(reg),
		van.WithName("gin-http"),
		van.WithServer(srv),
	)
	if err := service.Run(); err != nil {
		logx.Fatal(err)
	}
}
