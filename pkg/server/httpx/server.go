// Package httpx provides ...
package httpx

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/trustasia-com/go-van/pkg/internal"
	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/httpx/handler"

	"github.com/justinas/alice"
)

// NewServer new http server
func NewServer(opts ...server.ServerOption) *Server {
	options := server.ServerOptions{
		Network: "tcp",
		Address: ":0",
		Handler: http.DefaultServeMux,

		Flag: server.ServerStdFlag,
	}
	// apply option
	for _, o := range opts {
		o(&options)
	}

	svr := &Server{
		options: options,
	}
	svr.Server = &http.Server{Handler: svr}

	chain := alice.New()

	// flag apply options
	if options.Flag&server.FlagRecover > 0 {
		chain = chain.Append(handler.RecoverHandler)
	}
	if options.Flag&server.FlagTracing > 0 {
		chain = chain.Append(handler.TraceSrvHandler)
	}
	// from context
	if options.Context != nil {
		h, ok := options.Context.Value(corsOptKey{}).(func(http.Handler) http.Handler)
		if ok {
			chain = chain.Append(h)
		}
	}
	svr.Handler = chain.Then(svr.Handler)

	return svr
}

// Server http server
type Server struct {
	options server.ServerOptions

	*http.Server
}

// Start start http server
func (s *Server) Start() error {
	lis, err := net.Listen(s.options.Network, s.options.Address)
	if err != nil {
		return err
	}
	logx.Infof("[HTTP] server listening on: %s", lis.Addr().String())
	return s.Serve(lis)
}

// Stop stop http server
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	s.Shutdown(ctx)
	logx.Info("[HTTP] server stopping")
	return nil
}

// Endpoint return endpoint
func (s *Server) Endpoint() (string, error) {
	addr, err := internal.Extract(s.options.Address)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s", addr), nil
}

// ServeHTTP wrapper http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO

	// health check
	if r.URL.Path == "/ping" {
		w.Write([]byte("it's ok"))
		return
	}
	s.options.Handler.ServeHTTP(w, r)
}
