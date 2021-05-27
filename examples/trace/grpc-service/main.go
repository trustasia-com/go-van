// Package main provides ...
package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/deepzz0/go-van"
	pb "github.com/deepzz0/go-van/examples/grpc/helloworld"
	"github.com/deepzz0/go-van/pkg/codes"
	"github.com/deepzz0/go-van/pkg/codes/status"
	"github.com/deepzz0/go-van/pkg/logx"
	"github.com/deepzz0/go-van/pkg/server"
	"github.com/deepzz0/go-van/pkg/server/grpcx"
	"github.com/deepzz0/go-van/pkg/telemetry"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var httpClient = http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

func main() {
	shutdown := telemetry.InitProvider(
		telemetry.WithEndpoint("0.0.0.0:4317"),
		telemetry.WithTracerName("grpc-service-app"),
	)
	defer shutdown()

	// grpc server
	srv := grpcx.NewServer(
		server.WithAddress(":8000"),
		server.WithSrvFlag(server.FlagTracing),
	)
	s := &serverGRPC{}
	pb.RegisterGreeterServer(srv, s)

	service := van.NewService(
		van.WithName("grpc-service"),
		van.WithServer(srv),
	)
	if err := service.Run(); err != nil {
		logx.Fatal(err)
	}
}

type userServer struct {
	pb.UnimplementedUserServer
}

func (s *userServer) GetUserInfo(ctx context.Context, in *pb.UserInfoReq) (*pb.UserInfoResp, error) {
	if in.Id != "1" {
		return nil, status.Err(codes.NotFound, "no user")
	}
	resp := &pb.UserInfoResp{Username: "bob"}
	return resp, nil
}

func (s *userServer) GetUserInfoProxy(ctx context.Context, in *pb.UserInfoReq) (*pb.UserInfoResp, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"http://localhost:9001/user/"+in.Id, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	resp2 := &pb.UserInfoResp{Username: m["username"].(string)}
	return resp2, nil
}
