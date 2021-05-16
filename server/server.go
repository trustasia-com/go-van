// Package server provides ...
package server

// Server micro server
type Server interface {
	Start() error
	Stop() error
	Endpoint() (string, error)
}
