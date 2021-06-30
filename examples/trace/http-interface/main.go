// Package main provides ...
package main

import (
	"net/http"
	"time"

	"github.com/trustasia-com/go-van"
	pb "github.com/trustasia-com/go-van/examples/trace/proto"
	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/grpcx"
	"github.com/trustasia-com/go-van/pkg/server/httpx"
	"github.com/trustasia-com/go-van/pkg/telemetry"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

var (
	httpClient *httpx.Client
	grpcClient pb.UserClient
)

func main() {
	// grpc client
	conn, err := grpcx.DialContext(
		server.WithEndpoint("localhost:8000"),
		server.WithCliFlag(server.FlagTracing),
	)
	if err != nil {
		panic(err)
	}
	grpcClient = pb.NewUserClient(conn)
	// http client
	httpClient = httpx.NewClient(
		server.WithCliFlag(server.FlagTracing),
	)

	shutdown := telemetry.InitProvider(
		telemetry.WithEndpoint("0.0.0.0:4317"),
		telemetry.WithTracerName("http-interface-app"),
	)
	defer shutdown()

	r := gin.Default()
	r.GET("/http-to-http/:id", handleHTTP2HTTP)
	r.GET("/http-to-grpc/:id", handleHTTP2GRPC)
	r.GET("/http-to-grpc-to-http/:id", handleHTTP2GRPC2HTTP)
	srv := httpx.NewServer(
		server.WithAddress(":9000"),
		server.WithHandler(r),
		server.WithSrvFlag(server.FlagTracing),
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
	ctx := c.Request.Context()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"http://localhost:9001/user/"+id, nil)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var data []byte
	err = httpClient.Do(ctx, req, &data)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Data(http.StatusOK, "text/html", data)
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
		c.String(http.StatusBadRequest, err.Error())
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
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
