// Package etcd provides ...
package etcd

import (
	"context"

	"github.com/deepzz0/go-van/pkg/registry"
)

type authKey struct{}

type authCreds struct {
	username string
	password string
}

// WithAuth etcd auth creds
func WithAuth(username, password string) registry.Option {
	return func(opts *registry.Options) {
		if opts.Context == nil {
			opts.Context = context.Background()
		}
		creds := &authCreds{
			username: username,
			password: password,
		}
		opts.Context = context.WithValue(opts.Context, authKey{}, creds)
	}
}
