// Package http provides ...
package http

import (
	"net/http"

	"github.com/deepzz0/go-van/pkg/server"
)

// NewClient new http client
func NewClient(opts ...server.Option) (*http.Client, error) {
	options := server.Options{}

	for _, o := range opts {
		o(&options)
	}
	trans, ok := options.Context.Value(transportOptKey{}).(*http.Transport)
	if ok {
		return &http.Client{Transport: trans}, nil
	}
	return &http.Client{}, nil
}
