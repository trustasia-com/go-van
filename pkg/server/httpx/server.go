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
	"github.com/trustasia-com/go-van/pkg/telemetry"

	"github.com/justinas/alice"
)

// NewServer new http server
func NewServer(opts ...server.ServerOption) *Server {
	options := server.ServerOptions{
		Network: "tcp",
		Address: ":0",
		Handler: http.DefaultServeMux,

		Flag: server.FlagRecover,
	}
	// apply option
	for _, o := range opts {
		o(&options)
	}

	svr := &Server{
		network: options.Network,
		address: options.Address,
		handler: options.Handler,
	}
	svr.Server = &http.Server{Handler: svr}

	chain := alice.New()

	// flag apply options
	if options.Flag&server.FlagRecover > 0 {
		chain = chain.Append(handler.RecoverHandler)
	}
	// telemetry
	if len(options.Telemetry) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var flag telemetry.FlagOption
		svr.shutdown, flag = telemetry.InitProvider(ctx, options.Telemetry...)

		// if flag&telemetry.FlagMeter > 0 {
		// 	chain = chain.Append(handler.MeterSrvHandler)
		// }
		if flag&telemetry.FlagTracer > 0 {
			chain = chain.Append(handler.TracerSrvHandler)
		}
	}
	// from context
	if options.Context != nil {
		h, ok := options.Context.Value(corsOptKey{}).(func(http.Handler) http.Handler)
		if ok {
			chain = chain.Append(h)
		}
	}
	svr.handler = chain.Then(options.Handler)

	return svr
}

// Server http server
type Server struct {
	network  string
	address  string
	handler  http.Handler
	shutdown func()

	*http.Server
}

// Start start http server
func (s *Server) Start() error {
	lis, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	logx.Infof("[HTTP] server listening on: %s", lis.Addr().String())
	return s.Serve(lis)
}

// Stop stop http server
func (s *Server) Stop() error {
	logx.Info("[HTTP] server stopping")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	// telemetry
	if s.shutdown != nil {
		s.shutdown()
	}
	return s.Shutdown(ctx)
}

// Endpoint return endpoint
func (s *Server) Endpoint() (string, error) {
	addr, err := internal.Extract(s.address)
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
		w.Write([]byte("pong"))
		return
	}
	wrapper := &handler.WrappedWriter{}
	wrapper.ResponseWriter = w
	s.handler.ServeHTTP(wrapper, r)
}
