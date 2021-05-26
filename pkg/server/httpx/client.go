// Package httpx provides ...
package httpx

import (
	"context"
	"net/http"
	"time"

	"github.com/deepzz0/go-van/pkg/server"
)

// NewClient new http client
func NewClient(opts ...server.DialOption) *Client {
	options := server.DialOptions{
		Timeout: time.Second * 5,
	}
	cli := &Client{
		options:   options,
		transport: http.DefaultTransport,
	}
	cli.Client = &http.Client{Transport: cli}
	// apply option
	for _, o := range opts {
		o(&options)
	}
	// change transport?
	if options.Context != nil {
		trans, ok := options.Context.Value(transportOptKey{}).(*http.Transport)
		if ok {
			cli.transport = trans
		}
	}
	return cli
}

// Client wrapper http client
type Client struct {
	options   server.DialOptions
	transport http.RoundTripper

	*http.Client
}

// RoundTrip implements http.RoundTripper as http.Transport
func (c *Client) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.UserAgent() == "" && c.options.UserAgent != "" {
		req.Header.Set("User-Agent", c.options.UserAgent)
	}
	return c.transport.RoundTrip(req)
}

// Do request to server
func (c *Client) Do(ctx context.Context, req *http.Request, resp interface{}) error {
	// TODO
	return nil
}
