// Package main provides ...
package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/deepzz0/go-van"
	pb "github.com/deepzz0/go-van/examples/grpc/helloworld"
	"github.com/deepzz0/go-van/pkg/logx"
	"github.com/deepzz0/go-van/pkg/server"
	"github.com/deepzz0/go-van/pkg/server/grpcx"
	"github.com/deepzz0/go-van/pkg/server/httpx"
	"github.com/deepzz0/go-van/pkg/telemetry"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
	grpcClient = pb.NewGreeterClient(conn)
	// http client
	httpClient = httpx.NewClient(
		server.WithTransport(otelhttp.NewTransport(http.DefaultTransport)),
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
		server.WithFlag(server.FlagTracing),
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
	resp, err := httpClient.Do(req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
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
