// Package httpx provides ...
package httpx

import (
	"context"
	"net/http"

	"github.com/deepzz0/go-van/pkg/server"
)

type transportOptKey struct{}

// WithTransport http transport
func WithTransport(trans http.Transport) server.DialOption {
	return func(opts *server.DialOptions) {
		if opts.Context == nil {
			opts.Context = context.Background()
		}
		opts.Context = context.WithValue(opts.Context, transportOptKey{}, opts)
	}
}
