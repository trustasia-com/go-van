// Package confx provides ...
package confx

// WatchFunc file change exec
type WatchFunc = func(name string, data []byte) error

// Confx config loader
type Confx interface {
	// LoadFiles load config from backend
	LoadFiles(obj interface{}, files ...string) error
	// WatchFiles watch file change
	WatchFiles(do WatchFunc, files ...string) error
}
