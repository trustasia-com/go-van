// Package registry provides ...
package registry

import "context"

// Registry for discovery
type Registry interface {
	// Register service to registry
	Register(ctx context.Context, srv *Service) error
	// Deregister service from registry
	Deregister(ctx context.Context, srv *Service) error
	// GetService return the service in memory according to the service name.
	GetService(ctx context.Context, srvName string) ([]*Service, error)
	// Watch creates a watcher according to the service name.
	Watch(ctx context.Context, srvName string) (Watcher, error)
}

// Service instance for registry
type Service struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []string          `json:"endpoints"`
}

// Watcher is service watcher.
type Watcher interface {
	// Next is blocking call
	Next() (*Result, error)
	// Stop close the watcher.
	Stop() error
}

// EventType registry event type
type EventType int

// enum actions
const (
	// create service
	Create EventType = iota
	// delete service
	Delete
	// update service
	Update
)

// Result watcher result
type Result struct {
	EventType EventType
	Service   *Service
}
