// Package httpx provides ...
package httpx

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/httpx/handler"
	"github.com/trustasia-com/go-van/pkg/server/httpx/resolver"
)

// Client http client
type Client interface {
	Do(ctx context.Context, req *Request) (Response, error)
}

// NewClient new http client, concurrent security
func NewClient(opts ...server.DialOption) Client {
	options := server.DialOptions{
		Timeout: time.Second * 5,
		Flag:    0, // default flag
	}
	// apply option
	for _, o := range opts {
		o(&options)
	}

	cli := &client{
		endpoint:  options.Endpoint,
		userAgent: options.UserAgent,
	}
	cli.Client = &http.Client{Transport: cli}

	// transport apply
	cli.transport = http.DefaultTransport
	if options.Context != nil {
		trans, ok := options.Context.Value(transportOptKey{}).(*http.Transport)
		if ok {
			cli.transport = trans
		}
		header, ok := options.Context.Value(headerOptKey{}).(map[string]string)
		if ok {
			cli.httpHeader = header
		}
	}
	// NOTE discovery, experimental nature, not recommended
	if options.Registry != nil {
		builder := resolver.NewBuilder(options.Registry)
		cli.transport.(*http.Transport).DialContext, _ = builder.Build(options.Endpoint)
	}
	if options.Flag&server.FlagInsecure > 0 {
		trans, ok := cli.transport.(*http.Transport)
		if ok {
			trans.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		} else {
			logx.Warningf("httpx: insecure transport not supported")
		}
	}
	// apply client flag
	if options.Flag&server.FlagTracing > 0 {
		cli.transport = handler.TracerCliHandler(cli.transport)
	}
	return cli
}

// client wrapper http client
type client struct {
	endpoint   string
	userAgent  string
	httpHeader map[string]string
	transport  http.RoundTripper

	*http.Client
	addresses []string
}

// Do request to server
func (c *client) Do(ctx context.Context, req *Request) (resp Response, err error) {
	httpReq, err := req.HTTP(c.endpoint)
	if err != nil {
		return
	}
	// apply client header
	for k, v := range c.httpHeader {
		httpReq.Header.Add(k, v)
	}
	httpReq = httpReq.WithContext(ctx)
	httpResp, err := c.Client.Do(httpReq)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()
	resp.Response = httpResp

	// read data
	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return
	}
	resp.Data = data

	// check http status code
	if httpResp.StatusCode/100 != 2 {
		err = fmt.Errorf("httpx: http status: %s, body: %s", httpResp.Status, data)
		return
	}
	// no content
	if httpResp.StatusCode == 201 {
		return
	}
	return
}

// RoundTrip implements http.RoundTripper as http.Transport
func (c *client) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.UserAgent() == "" && c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	return c.transport.RoundTrip(req)
}
