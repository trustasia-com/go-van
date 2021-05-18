// Package etcd provides ...
package etcd

import (
	"context"

	"github.com/deepzz0/go-van/registry"
)

type authKey struct{}

type authCreds struct {
	username string
	password string
}

// Auth etcd auth creds
func Auth(username, password string) registry.Option {
	return func(opts *registry.Options) {
		if opts.Ctx == nil {
			opts.Ctx = context.Background()
		}
		creds := &authCreds{
			username: username,
			password: password,
		}
		opts.Ctx = context.WithValue(opts.Ctx, authKey{}, creds)
	}
}
