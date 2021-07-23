// Package httpx provides ...
package httpx

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"time"

	"github.com/trustasia-com/go-van/pkg/codec/json"
	"github.com/trustasia-com/go-van/pkg/codec/xml"
	"github.com/trustasia-com/go-van/pkg/server"
	"github.com/trustasia-com/go-van/pkg/server/httpx/handler"
	"github.com/trustasia-com/go-van/pkg/server/httpx/resolver"
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
	transport := http.DefaultTransport
	if options.Context != nil {
		trans, ok := options.Context.Value(transportOptKey{}).(*http.Transport)
		if ok {
			transport = trans
		}
	}
	// discovery
	if options.Registry != nil {
		builder := resolver.NewBuilder(options.Registry)
		transport.(*http.Transport).DialContext, _ = builder.Build()
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

// RoundTrip implements http.RoundTripper as http.Transport
func (c *client) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.UserAgent() == "" && c.options.UserAgent != "" {
		req.Header.Set("User-Agent", c.options.UserAgent)
	}
	return c.transport.RoundTrip(req)
}

// Do request to server
func (c *client) Do(req *http.Request, resp interface{}) error {
	httpResp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	// check length
	if httpResp.StatusCode == 201 || resp == nil {
		return nil
	}

	ct := httpResp.Header.Get("Content-Type")
	media, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil
	}
	fmt.Println(media)
	switch media {
	case "application/xml": // codec xml
		err = xml.NewCodec().Unmarshal(data, resp)
	case "application/json": // codec json
		err = json.NewCodec().Unmarshal(data, resp)
	case "text/html": // codec text
		str := string(data)
		*resp.(*string) = str
	default:
		*resp.(*[]byte) = data
	}
	return err
}
