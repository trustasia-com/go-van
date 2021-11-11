// Package httpx provides ...
package httpx

import (
	"fmt"
	"mime"
	"net/http"

	"github.com/trustasia-com/go-van/pkg/codec/json"
	"github.com/trustasia-com/go-van/pkg/codec/xml"
)

// Response wrapped http response
type Response struct {
	Response *http.Response

	Data []byte
}

// HTTP return http response, Body is nil (having already been consumed).
func (resp Response) HTTP() *http.Response {
	return resp.Response
}

// Scan scan data to struct
func (resp Response) Scan(p interface{}) error {
	if len(resp.Data) == 0 {
		return fmt.Errorf("httpx: no content")
	}

	ct := resp.Response.Header.Get("Content-Type")
	media, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return fmt.Errorf("httpx: invalid Content-Type: %s", ct)
	}
	switch media {
	case "application/xml", "application/xhtml+xml": // codec xml
		err = xml.NewCodec().Unmarshal(resp.Data, p)
	case "application/json": // codec json
		err = json.NewCodec().Unmarshal(resp.Data, p)
	default:
		*p.(*[]byte) = resp.Data
	}
	return nil
}
