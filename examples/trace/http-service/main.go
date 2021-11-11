// Package main provides ...
package main

import (
	"context"
	"net/http"

	"github.com/trustasia-com/go-van"
	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/httpx"
	"github.com/trustasia-com/go-van/pkg/telemetry"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	shutdown := telemetry.InitProvider(
		context.Background(),
		telemetry.WithEndpoint("192.168.252.177:4317"),
		telemetry.WithTracerName("http-service-app"),
		telemetry.WithOptions(grpc.WithInsecure()),
	)
	defer shutdown()

	r := gin.Default()
	r.GET("/user/:id", handleUserInfo)

	srv := httpx.NewServer(
		server.WithAddress(":9001"),
		server.WithHandler(r),
		server.WithSrvFlag(server.FlagTracing),
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
