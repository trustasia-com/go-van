// Package registry provides ...
package registry

import "context"

// Registry for discovery
type Registry interface {
	Register(ctx context.Context, srv *Service) error
	Deregister(ctx context.Context, srv *Service) error
}

// Service instance for registry
type Service struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []string          `json:"endpoints"`
	// Endpoints []*Endpoint       `json:"endpoints"`
	// Nodes     []*Node           `json:"nodes"`
}

// Discovery is service discovery.
type Discovery interface {
	// GetService return the service instances in memory according to the service name.
	GetService(ctx context.Context, serviceName string) ([]*Service, error)
	// Watch creates a watcher according to the service name.
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}

// Watcher is service watcher.
type Watcher interface {
	// Next returns services in the following two cases:
	// 1.the first time to watch and the service instance list is not empty.
	// 2.any service instance changes found.
	// if the above two conditions are not met, it will block until context deadline exceeded or canceled
	Next() ([]*Service, error)
	// Stop close the watcher.
	Stop() error
}
