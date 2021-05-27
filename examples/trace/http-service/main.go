// Package main provides ...
package main

import (
	"net/http"

	"github.com/deepzz0/go-van"
	"github.com/deepzz0/go-van/pkg/logx"
	"github.com/deepzz0/go-van/pkg/server"
	"github.com/deepzz0/go-van/pkg/server/httpx"
	"github.com/deepzz0/go-van/pkg/telemetry"

	"github.com/gin-gonic/gin"
)

func main() {
	shutdown := telemetry.InitProvider(
		telemetry.WithEndpoint("0.0.0.0:4317"),
		telemetry.WithTracerName("http-service-app"),
	)
	defer shutdown()

	r := gin.Default()
	r.GET("/user/:id", handleUserInfo)

	srv := httpx.NewServer(
		server.WithAddress(":9001"),
		server.WithHandler(r),
		server.WithFlag(server.FlagTracing),
	)
	service := van.NewService(
		van.WithName("http-service"),
		van.WithServer(srv),
	)
	if err := service.Run(); err != nil {
		logx.Fatal(err)
	}
}

func handleUserInfo(c *gin.Context) {
	id := c.Param("id")
	if id != "1" { // err
		c.String(http.StatusBadRequest, "not found")
		return
	}
	u := map[string]interface{}{
		"username": "bob",
		"age":      10,
	}
	c.JSON(http.StatusOK, u)
}
