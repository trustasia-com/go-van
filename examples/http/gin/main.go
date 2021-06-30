// Package main provides ...
package main

import (
	"github.com/trustasia-com/go-van"
	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/httpx"

	"github.com/gin-gonic/gin"
)

func main() {
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

	srv := httpx.NewServer(
		server.WithAddress(":9000"),
		server.WithHandler(e),
	)
	service := van.NewService(
		van.WithName("gin-http"),
		van.WithServer(srv),
	)
	if err := service.Run(); err != nil {
		logx.Fatal(err)
	}
}
