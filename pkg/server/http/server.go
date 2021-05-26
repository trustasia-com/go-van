// Package http provides ...
package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/deepzz0/go-van/pkg/logx"
	"github.com/deepzz0/go-van/pkg/server"
	"github.com/deepzz0/go-van/pkg/tools"
)

// NewServer new http server
func NewServer(opts ...server.Option) server.Server {
	opt := server.Options{
		Network:  "tcp",
		Endpoint: ":0",
		Timeout:  time.Second,
		Context:  context.Background(),
	}
	svr := &httpServer{
		handler: http.DefaultServeMux,
		options: opt,
	}
	svr.Server = &http.Server{Handler: svr}
	// apply option
	for _, o := range opts {
		o(&svr.options)
	}
	// handler opts from context
	h, ok := svr.options.Context.Value(handlerOptKey{}).(http.Handler)
	if ok {
		svr.handler = h
	}
	return svr
}

// httpServer http server
type httpServer struct {
	handler http.Handler
	options server.Options

	*http.Server
}

// Start start http server
func (s *httpServer) Start() error {
	lis, err := net.Listen(s.options.Network, s.options.Endpoint)
	if err != nil {
		return err
	}
	logx.Infof("[HTTP] server listening on: %s", lis.Addr().String())
	return s.Serve(lis)
}

// Stop stop http server
func (s *httpServer) Stop() error {
	s.Shutdown(s.options.Context)
	return nil
}

// Endpoint return endpoint
func (s *httpServer) Endpoint() (string, error) {
	addr, err := tools.Extract(s.options.Endpoint)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s", addr), nil
}

// ServeHTTP wrapper http.Handler
func (s *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// more
	s.handler.ServeHTTP(w, r)
}
