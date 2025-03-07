// Package confx provides ...
package confx

import "context"

// WatchFunc file change exec
type WatchFunc = func(name string, data []byte) error

// Confx config loader
type Confx interface {
	// LoadFiles load config from backend
	LoadFiles(obj any, files ...string) error
	// WatchFiles watch file change
	WatchFiles(ctx context.Context, do WatchFunc, files ...string) error
}
