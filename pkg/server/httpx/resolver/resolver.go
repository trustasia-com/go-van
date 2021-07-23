// Package resolver provides ...
package resolver

import (
	"context"
	"net"
)

// Builder creates a resolver
type Builder interface {
	Build() (Dialer, error)
}

// Dialer return dial function
type Dialer func(ctx context.Context, network, addr string) (net.Conn, error)
