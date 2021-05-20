// Package registry provides ...
package registry

import "context"

// Registry for discovery
type Registry interface {
	// Register service to registry
	Register(ctx context.Context, ins *Instance) error
	// Deregister service from registry
	Deregister(ctx context.Context, ins *Instance) error
	// GetService return the service in memory according to the service name.
	GetService(ctx context.Context, name string) ([]*Instance, error)
	// Watch creates a watcher according to the service name.
	Watch(ctx context.Context, name string) (Watcher, error)
}

// Instance instance for registry
type Instance struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []string          `json:"endpoints"`
}

// Watcher is service watcher.
type Watcher interface {
	// Next is blocking call
	Next() ([]*Instance, error)
	// Stop close the watcher.
	Stop() error
}
