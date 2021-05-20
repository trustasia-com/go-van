// Package http provides ...
package http

import (
	"context"
	"net/http"

	"github.com/deepzz0/go-van/pkg/server"
)

type handlerOptKey struct{}

// WithHandler http handler
func WithHandler(h http.Handler) server.Option {
	return func(opts *server.Options) {
		if opts.Context == nil {
			opts.Context = context.Background()
		}
		opts.Context = context.WithValue(opts.Context, handlerOptKey{}, opts)
	}
}

type transportOptKey struct{}

// WithTransport http transport
func WithTransport(trans http.Transport) server.Option {
	return func(opts *server.Options) {
		if opts.Context == nil {
			opts.Context = context.Background()
		}
		opts.Context = context.WithValue(opts.Context, transportOptKey{}, opts)
	}
}
