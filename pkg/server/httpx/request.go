// Package httpx provides ...
package httpx

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

// Request for http request
type Request struct {
	Method string      // http method
	Header http.Header // header
	URL    string      // request url
	Query  url.Values  // query uri
	Body   io.Reader   // request body

	Context context.Context // context
}

// httpRequest generate http request
func (req *Request) httpRequest() (*http.Request, error) {
	u, err := url.Parse(req.URL)
	if err != nil {
		return nil, err
	}
	if req.Query != nil {
		u.RawQuery = req.Query.Encode()
	}
	httpReq, err := http.NewRequest(req.Method, u.String(), req.Body)
	if err != nil {
		return nil, err
	}
	// context
	if req.Context != nil {
		httpReq = httpReq.WithContext(req.Context)
	}
	if req.Header != nil {
		httpReq.Header = req.Header
	}
	return httpReq, nil
}

// TODO more function migrate Request
