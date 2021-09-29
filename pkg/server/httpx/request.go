// Package httpx provides ...
package httpx

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

// NewRequest new request for client
func NewRequest(method, path string, body io.Reader) *Request {
	req := &Request{
		method: method,
		path:   path,
		body:   body,

		header: make(map[string][]string),
	}
	return req
}

// Request for http request
type Request struct {
	method string    // http method
	path   string    // url path & query
	body   io.Reader // request body

	header  http.Header     // header
	context context.Context // context
}

// ToHTTP generate http request
func (req *Request) ToHTTP(host string) (*http.Request, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	u.Path = req.path
	httpReq, err := http.NewRequest(req.method, u.String(), req.body)
	if err != nil {
		return nil, err
	}
	// context
	if req.context != nil {
		httpReq = httpReq.WithContext(req.context)
	}
	if req.header != nil {
		httpReq.Header = req.header
	}
	return httpReq, nil
}

// TODO more function migrate Request

// SetHeader set http header
func (req *Request) SetHeader(key, val string) {
	req.header.Set(key, val)
}

// AddHeader add http header
func (req *Request) AddHeader(key, val string) {
	req.header.Add(key, val)
}

// Context set context
func (req *Request) Context(ctx context.Context) {
	req.context = ctx
}
