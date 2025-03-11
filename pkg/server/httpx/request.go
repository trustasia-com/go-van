// Package httpx provides ...
package httpx

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
)

// NewRequest new request for client
func NewRequest(method, path, query string, body []byte) *Request {
	req := &Request{
		method: method,
		path:   path,
		query:  query,
		body:   body,

		header: make(map[string][]string),
	}
	return req
}

// Request for http request
type Request struct {
	method             string // http method
	path               string // url path
	query              string // query
	body               []byte // request body
	username, password string // basic auth

	header  http.Header     // header
	context context.Context // context
}

// HTTP generate http request
// host, eg. https://example.com:8080
// path, eg. /users
// query, eg. page=1&size=10
func (req *Request) HTTP(host string) (httpReq *http.Request, err error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	u.JoinPath(req.path)
	u.RawQuery = req.query
	if len(req.body) > 0 {
		httpReq, err = http.NewRequest(req.method, u.String(), bytes.NewReader(req.body))
	} else {
		httpReq, err = http.NewRequest(req.method, u.String(), nil)
	}
	if err != nil {
		return nil, err
	}
	// context
	if req.context != nil {
		httpReq = httpReq.WithContext(req.context)
	}
	if len(req.header) > 0 {
		httpReq.Header = req.header
	}
	if req.username != "" && req.password != "" {
		httpReq.SetBasicAuth(req.username, req.password)
	}
	return httpReq, nil
}

// GetMethod return method
func (req *Request) GetMethod() string {
	return req.method
}

// GetPath return path
func (req *Request) GetPath() string {
	return req.path
}

// GetQuery return query
func (req *Request) GetQuery() string {
	return req.query
}

// GetBody return body
func (req *Request) GetBody() []byte {
	return req.body
}

// GetHeader get header
func (req *Request) GetHeader() http.Header {
	return req.header
}

// GetContext return context
func (req *Request) GetContext() context.Context {
	return req.context
}

// SetBasicAuth basic auth
func (req *Request) SetBasicAuth(username, password string) {
	req.username = username
	req.password = password
}

// SetHeader set http header
func (req *Request) SetHeader(key, val string) {
	req.header.Set(key, val)
}

// AddHeader add http header
func (req *Request) AddHeader(key, val string) {
	req.header.Add(key, val)
}

// SetContext set context
func (req *Request) SetContext(ctx context.Context) {
	req.context = ctx
}
