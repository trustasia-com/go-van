// Package httpx provides ...
package httpx

import (
	"context"
	"net/http"

	"github.com/trustasia-com/go-van/pkg/server"
)

type transportOptKey struct{}

// WithTransport http transport for client
func WithTransport(trans http.RoundTripper) server.DialOption {
	return func(opts *server.DialOptions) {
		if opts.Context == nil {
			opts.Context = context.Background()
		}
		opts.Context = context.WithValue(opts.Context, transportOptKey{}, trans)
	}
}
