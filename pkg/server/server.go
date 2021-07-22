// Package server provides ...
package server

import (
	"context"
	"net/http"
)

// Server micro server
type Server interface {
	Start() error
	Stop() error
	Endpoint() (string, error)
}

// HTTPClient http client
type HTTPClient interface {
	Do(ctx context.Context, req *http.Request, resp interface{}) error
}
