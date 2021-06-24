// Package apollo provides ...
package apollo

import (
	"context"

	"github.com/deepzz0/go-van/pkg/confx"
)

type apolloLoader struct{}

// LoadFiles load config from backend
func (l *apolloLoader) LoadFiles(obj interface{}, glob string) error {

	return nil
}

// WatchFiles watch file change
func (l *apolloLoader) WatchFiles(ctx context.Context, do confx.WatchFunc, glob string) error {

	return nil
}
