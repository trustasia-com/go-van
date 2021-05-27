// Package httpx provides ...
package httpx

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/deepzz0/go-van/pkg/internal"
	"github.com/deepzz0/go-van/pkg/logx"
	"github.com/deepzz0/go-van/pkg/server"
	"github.com/deepzz0/go-van/pkg/server/httpx/handler"

	"github.com/justinas/alice"
)

// NewServer new http server
func NewServer(opts ...server.ServerOption) *httpServer {
	options := server.ServerOptions{
		Network: "tcp",
		Address: ":0",
		Handler: http.DefaultServeMux,
	}
	svr := &httpServer{
		options: options,
	}
	svr.Server = &http.Server{Handler: svr}
	// apply option
	for _, o := range opts {
		o(&svr.options)
	}
	// recover options
	chain := alice.New()
	if svr.options.Recover {
		chain = chain.Append(handler.RecoverHandler)
	}
	svr.Handler = chain.Then(svr.Handler)
	return svr
}

// httpServer http server
type httpServer struct {
	options server.ServerOptions

	*http.Server
}

// Start start http server
func (s *httpServer) Start() error {
	lis, err := net.Listen(s.options.Network, s.options.Address)
	if err != nil {
		return err
	}
	logx.Infof("[HTTP] server listening on: %s", lis.Addr().String())
	return s.Serve(lis)
}

// Stop stop http server
func (s *httpServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	s.Shutdown(ctx)
	logx.Info("[HTTP] server stopping")
	return nil
}

// Endpoint return endpoint
func (s *httpServer) Endpoint() (string, error) {
	addr, err := internal.Extract(s.options.Address)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s", addr), nil
}

// ServeHTTP wrapper http.Handler
func (s *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO
	// more
	// eg. health check
	s.options.Handler.ServeHTTP(w, r)
}
