// Package server provides ...
package server

import (
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
	Do(req *http.Request, resp interface{}) error
}
