// Package main provides ...
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/trustasia-com/go-van"
	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/httpx"
	"github.com/trustasia-com/go-van/pkg/telemetry"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var appName = "http-service-app"

func main() {
	r := gin.Default()
	r.GET("/user/:id", handleUserInfo)

	srv := httpx.NewServer(
		server.WithAddress(":9001"),
		server.WithHandler(r),
		server.WithSrvFlag(server.FlagTracing),
		server.WithTelemetry(
			telemetry.WithEndpoint("localhost:4317"),
			telemetry.WithName(appName),
			telemetry.WithOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
		),
	)
	service := van.NewService(
		van.WithName(appName),
		van.WithServer(srv),
	)
	if err := service.Run(); err != nil {
		logx.Fatal(err)
	}
}

func handleUserInfo(c *gin.Context) {
	meter := otel.Meter(appName)

	opt := metric.WithAttributes(
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	)
	// This is the equivalent of prometheus.NewCounterVec
	counter, err := meter.Float64Counter("foo", metric.WithDescription("a simple counter"))
	if err != nil {
		log.Fatal(err)
	}
	counter.Add(context.Background(), 1, opt)

	id := c.Param("id")
	if id != "1" { // err
		// logging with trace id
		logx.WithContext(c.Request.Context()).Error("not found")
		c.String(http.StatusBadRequest, "not found")
		return
	}
	u := map[string]any{
		"username": "bob",
		"age":      10,
	}
	c.JSON(http.StatusOK, u)
}
