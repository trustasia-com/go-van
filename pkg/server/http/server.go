// Package http provides ...
package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/deepzz0/go-van/pkg/server"
)

// NewServer new http server
func NewServer(opts ...server.Option) (server.Server, error) {
	opt := server.Options{
		Network:  "tcp",
		Endpoint: ":0",
		Timeout:  time.Second,
		Context:  context.Background(),
	}
	svr := &httpServer{options: opt}
	// apply option
	for _, o := range opts {
		o(&svr.options)
	}
	// handler opts from context
	h, ok := svr.options.Context.Value(handlerOptKey{}).(http.Handler)
	if ok {
		svr.Server = &http.Server{Handler: h}
	}
	return svr, nil
}

// httpServer http server
type httpServer struct {
	options server.Options

	*http.Server
}

// Start start http server
func (s *httpServer) Start() error {
	lis, err := net.Listen(s.options.Network, s.options.Endpoint)
	if err != nil {
		return err
	}
	return s.Serve(lis)
}

// Stop stop http server
func (s *httpServer) Stop() error {
	s.Shutdown(s.options.Context)
	return nil
}

// Endpoint return endpoint
func (s *httpServer) Endpoint() (string, error) {
	return fmt.Sprintf("http://%s", s.options.Endpoint), nil
}
