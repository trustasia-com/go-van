// Package httpx provides ...
package httpx

import (
	"context"
	"net/http"

	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/httpx/handler"
)

type transportOptKey struct{}

// WithTransport http transport for client
func WithTransport(trans *http.Transport) server.DialOption {
	return func(opts *server.DialOptions) {
		if opts.Context == nil {
			opts.Context = context.Background()
		}
		opts.Context = context.WithValue(opts.Context, transportOptKey{}, trans)
	}
}

type corsOptKey struct{}

// CORSOptions cors options
type CORSOptions = handler.CORSOptions

// WithCORS http cross origin resource share
func WithCORS(hOpts CORSOptions) server.ServerOption {
	cors := handler.New(hOpts)

	return func(opts *server.ServerOptions) {
		if opts.Context == nil {
			opts.Context = context.Background()
		}
		opts.Context = context.WithValue(opts.Context, corsOptKey{}, cors.Handler)
	}
}

type headerOptKey struct{}

// WithHeader http header
func WithHeader(h map[string]string) server.DialOption {
	return func(opts *server.DialOptions) {
		if opts.Context == nil {
			opts.Context = context.Background()
		}
		opts.Context = context.WithValue(opts.Context, headerOptKey{}, h)
	}
}
