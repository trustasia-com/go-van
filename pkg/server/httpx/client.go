// Package httpx provides ...
package httpx

import (
	"context"
	"net/http"
	"time"

	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/httpx/handler"
)

// NewClient new http client, concurrent security
func NewClient(opts ...server.DialOption) server.HTTPClient {
	options := server.DialOptions{
		Timeout: time.Second * 5,
	}
	cli := &client{options: options}
	cli.Client = &http.Client{Transport: cli}
	// apply option
	for _, o := range opts {
		o(&options)
	}
	// transport apply
	if options.Context != nil {
		trans, ok := options.Context.Value(transportOptKey{}).(http.RoundTripper)
		if ok {
			cli.transport = trans
		}
	}

	// apply client flag
	if options.Flag&server.FlagTracing > 0 {
		cli.transport = handler.TraceCliHandler(cli.transport)
	}

	return cli
}

// client wrapper http client
type client struct {
	options   server.DialOptions
	transport http.RoundTripper

	*http.Client
}

// RoundTrip implements http.RoundTripper as http.Transport
func (c *client) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.UserAgent() == "" && c.options.UserAgent != "" {
		req.Header.Set("User-Agent", c.options.UserAgent)
	}
	// default transport
	if c.transport == nil {
		c.transport = http.DefaultTransport
	}
	return c.transport.RoundTrip(req)
}

// Do request to server
func (c *client) Do(ctx context.Context, req *http.Request, resp interface{}) error {
	if c.options.Registry != nil {
	}
	// TODO
	//
	// codec xml
	// codec json
	// codec text
	// codec bianry
	return nil
}
