// Package httpx provides ...
package httpx

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/httpx/handler"
	"github.com/trustasia-com/go-van/pkg/server/httpx/resolver"
)

// HTTPClient http client
type HTTPClient interface {
	Do(req *Request) (Response, error)
}

// NewClient new http client, concurrent security
func NewClient(opts ...server.DialOption) HTTPClient {
	options := server.DialOptions{
		Timeout: time.Second * 5,
	}
	// apply option
	for _, o := range opts {
		o(&options)
	}

	cli := &client{options: options}
	cli.Client = &http.Client{Transport: cli}

	// transport apply
	transport := http.DefaultTransport
	if options.Context != nil {
		trans, ok := options.Context.Value(transportOptKey{}).(*http.Transport)
		if ok {
			transport = trans
		}
	}
	// NOTE discovery, experimental nature, not recommended
	if options.Registry != nil {
		builder := resolver.NewBuilder(options.Registry)
		transport.(*http.Transport).DialContext, _ = builder.Build(options.Endpoint)
	}
	cli.transport = transport

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
	addresses []string
}

// Do request to server
func (c *client) Do(req *Request) (resp Response, err error) {
	httpReq, err := req.ToHTTP(c.options.Endpoint)
	if err != nil {
		return
	}
	httpResp, err := c.Client.Do(httpReq)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()
	resp.response = httpResp

	// check http status code
	if httpResp.StatusCode/100 != 2 {
		err = fmt.Errorf("httpx: http status: %s", httpResp.Status)
		return
	}
	// no content
	if httpResp.StatusCode == 201 {
		return
	}
	// check content length
	if httpResp.ContentLength > 1<<10 { // 1m
		err = fmt.Errorf("httpx: too large: %d", httpResp.ContentLength)
		return
	}
	// read data
	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return
	}
	resp.Data = data
	return
}

// RoundTrip implements http.RoundTripper as http.Transport
func (c *client) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.UserAgent() == "" && c.options.UserAgent != "" {
		req.Header.Set("User-Agent", c.options.UserAgent)
	}
	return c.transport.RoundTrip(req)
}
