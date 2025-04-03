// Package main provides ...
package main

import (
	"net/http"
	"time"

	"github.com/trustasia-com/go-van"
	pb "github.com/trustasia-com/go-van/examples/trace/proto"
	"github.com/trustasia-com/go-van/pkg/codes/status"
	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/grpcx"
	"github.com/trustasia-com/go-van/pkg/server/httpx"
	"github.com/trustasia-com/go-van/pkg/telemetry"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

var (
	httpClient httpx.Client
	grpcClient pb.UserClient
)

func init() {
	// grpc client
	conn, err := grpcx.DialContext(
		server.WithEndpoint("localhost:8000"),
		server.WithCliFlag(server.FlagTracing|server.FlagInsecure),
	)
	if err != nil {
		panic(err)
	}
	grpcClient = pb.NewUserClient(conn)
}

func main() {
	// http client
	httpClient = httpx.NewClient(
		server.WithEndpoint("http://localhost:9001"),
		server.WithCliFlag(server.FlagTracing),
	)

	r := gin.Default()
	r.GET("/http-to-http/:id", handleHTTP2HTTP)
	r.GET("/http-to-grpc/:id", handleHTTP2GRPC)
	r.GET("/http-to-grpc-to-http/:id", handleHTTP2GRPC2HTTP)
	srv := httpx.NewServer(
		server.WithAddress(":9000"),
		server.WithHandler(r),
		server.WithTelemetry(
			telemetry.WithEndpoint("localhost:4317"),
			telemetry.WithName("http-interface-app"),
			telemetry.WithFlag(telemetry.FlagInsecure|telemetry.FlagTracer|telemetry.FlagMeter),
		),
	)
	service := van.NewService(
		van.WithName("http-interface"),
		van.WithServer(srv),
	)
	if err := service.Run(); err != nil {
		logx.Fatal(err)
	}
}

func handleHTTP2HTTP(c *gin.Context) {
	id := c.Param("id")

	// span
	req := httpx.NewRequest(http.MethodGet, "/user/"+id, "", nil)
	ctx := c.Request.Context()

	resp, err := httpClient.Do(ctx, req)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			c.String(http.StatusBadRequest, s.Message())
		} else {
			c.String(http.StatusBadRequest, err.Error())
		}
		return
	}
	c.Data(http.StatusOK, "text/html", resp.Data)
}

func handleHTTP2GRPC(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"id", id,
	))
	resp, err := grpcClient.GetUserInfo(ctx, &pb.UserInfoReq{Id: id})
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			c.String(http.StatusBadRequest, s.Message())
		} else {
			c.String(http.StatusBadRequest, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, resp)
}

func handleHTTP2GRPC2HTTP(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"id", id,
	))
	resp, err := grpcClient.GetUserInfoProxy(ctx, &pb.UserInfoReq{Id: id})
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			c.String(http.StatusBadRequest, s.Message())
		} else {
			c.String(http.StatusBadRequest, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, resp)
}
