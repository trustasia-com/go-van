// Package handler provides ...
package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/trustasia-com/go-van/pkg/logx"
)

//
// the middleware from: https://github.com/rs/cors
//

// cors headers
const (
	HeaderOrigin              = "Origin"
	HeaderAccept              = "Accept"
	HeaderContentType         = "Content-Type"
	HeaderCustomRequestedWith = "X-Requested-With"

	CORSHeaderVary                          = "Vary"
	CORSHeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	CORSHeaderAccessControlRequestMethods   = "Access-Control-Request-Methods"
	CORSHeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	CORSHeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	CORSHeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	CORSHeaderAccessControlMaxAge           = "Access-Control-MaxAge"
	CORSHeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
)

// CORSOptions cors options
// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/CORS
type CORSOptions struct {
	// default: *
	AllowedOrigins []string
	AllowedMethods []string
	// Cache-Control、Content-Language、Content-Type、Expires、Last-Modified、Pragma
	AllowedHeaders   []string
	ExposedHeaders   []string
	MaxAge           int // s
	AllowCredentials bool
}

func (opts CORSOptions) handlePreflight(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()

	origin := r.Header.Get(HeaderOrigin)
	if r.Method != http.MethodOptions {
		logx.Warningf(" Preflight aborted: %s!=Options", r.Method)
		return
	}
	// Always set Vary headers
	// see https://github.com/rs/cors/issues/10,
	//     https://github.com/rs/cors/commit/dbdca4d95feaa7511a46e6f1efb3b3aa505bc43f#commitcomment-12352001
	headers.Add(CORSHeaderVary, HeaderOrigin)
	headers.Add(CORSHeaderVary, CORSHeaderAccessControlRequestMethod)
	headers.Add(CORSHeaderVary, CORSHeaderAccessControlRequestHeaders)

	if origin == "" {
		logx.Warningf(" Preflight aborted: empty origin")
		return
	}
	if !opts.isOriginAllowed(r, origin) {
		logx.Warningf(" Preflight aborted: origin '%s' not allowed", origin)
		return
	}
	reqMethod := r.Header.Get(CORSHeaderAccessControlRequestMethod)
	if !opts.isMethodAllowed(reqMethod) {
		logx.Warningf(" Preflight aborted: headers '%v' not allowed", reqMethod)
		return
	}
	reqHeaders := strings.Fields(r.Header.Get(CORSHeaderAccessControlRequestHeaders))
	if !opts.areHeadersAllowed(reqHeaders) {
		logx.Warningf(" Preflight aborted: headers '%v' not allowed", reqHeaders)
		return
	}

	if opts.AllowedOrigins[0] == "*" {
		headers.Set(CORSHeaderAccessControlAllowOrigin, "*")
	} else {
		headers.Set(CORSHeaderAccessControlAllowOrigin, origin)
	}
	// Spec says: Since the list of methods can be unbounded, simply returning the method indicated
	// by Access-Control-Request-Method (if supported) can be enough
	headers.Set(CORSHeaderAccessControlRequestMethods, strings.ToUpper(reqMethod))
	if len(reqHeaders) > 0 {
		// Spec says: Since the list of headers can be unbounded, simply returning supported headers
		// from Access-Control-Request-Headers can be enough
		headers.Set(CORSHeaderAccessControlRequestHeaders, strings.Join(reqHeaders, ", "))
	}
	if opts.AllowCredentials {
		headers.Set(CORSHeaderAccessControlAllowCredentials, "true")
	}
	if opts.MaxAge > 0 {
		headers.Set(CORSHeaderAccessControlMaxAge, strconv.Itoa(opts.MaxAge))
	}
	logx.Infof(" Preflight response headers: %v", headers)
}

func (opts CORSOptions) handleActualRequest(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get(HeaderOrigin)

	// Always set Vary, see https://github.com/rs/cors/issues/10
	headers.Add(CORSHeaderVary, HeaderOrigin)
	if origin == "" {
		logx.Warning("  Actual request no headers added: missing origin")
		return
	}
	if !opts.isOriginAllowed(r, origin) {
		logx.Warningf("  Actual request no headers added: origin '%s' not allowed", origin)
		return
	}
	// Note that spec does define a way to specifically disallow a simple method like GET or
	// POST. Access-Control-Allow-Methods is only used for pre-flight requests and the
	// spec doesn't instruct to check the allowed methods for simple cross-origin requests.
	// We think it's a nice feature to be able to have control on those methods though.
	if !opts.isMethodAllowed(r.Method) {
		logx.Warningf("  Actual request no headers added: method '%s' not allowed", r.Method)

		return
	}
	if opts.AllowedOrigins[0] == "*" {
		headers.Set(CORSHeaderAccessControlAllowOrigin, "*")
	} else {
		headers.Set(CORSHeaderAccessControlAllowOrigin, origin)
	}
	if len(opts.ExposedHeaders) > 0 {
		headers.Set(CORSHeaderAccessControlExposeHeaders, strings.Join(opts.ExposedHeaders, ", "))
	}
	if opts.AllowCredentials {
		headers.Set(CORSHeaderAccessControlAllowCredentials, "true")
	}
	logx.Infof("  Actual response added headers: %v", headers)
}

func (opts CORSOptions) isOriginAllowed(r *http.Request, origin string) bool {
	// allow all
	if opts.AllowedOrigins[0] == "*" {
		return true
	}

	origin = strings.ToLower(origin)
	for _, o := range opts.AllowedOrigins {
		if o == origin {
			return true
		}
		if idx := strings.IndexByte(origin, '*'); idx >= 0 {
			prefix := o[:idx]
			suffix := o[idx:]
			if len(o) >= len(prefix)+len(suffix) && strings.HasPrefix(origin, prefix) &&
				strings.HasSuffix(origin, suffix) {
				return true
			}
		}
	}
	return false
}

func (opts CORSOptions) isMethodAllowed(method string) bool {
	if len(opts.AllowedMethods) == 0 {
		// If no method allowed, always return false, even for preflight request
		return false
	}
	method = strings.ToUpper(method)
	if method == http.MethodOptions {
		// Always allow preflight requests
		return true
	}
	for _, m := range opts.AllowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

// areHeadersAllowed checks if a given list of headers are allowed to used within
// a cross-domain request.
func (opts CORSOptions) areHeadersAllowed(requestedHeaders []string) bool {
	if opts.AllowedHeaders[0] == "*" || len(requestedHeaders) == 0 {
		return true
	}
	for _, header := range requestedHeaders {
		header = http.CanonicalHeaderKey(header)
		found := false
		for _, h := range opts.AllowedHeaders {
			if h == header {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

type converter func(string) string

// convert converts a list of string using the passed converter function
func convert(s []string, c converter) []string {
	out := []string{}
	for _, i := range s {
		out = append(out, c(i))
	}
	return out
}

// CORSAllowAll create a new Cors handler with permissive configuration allowing all
// origins with all standard methods with any header and credentials.
func CORSAllowAll() CORSOptions {
	return CORSOptions{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}
}

// CORSHandler return a middleware that cors
func CORSHandler(opts CORSOptions) func(http.Handler) http.Handler {
	// exposed headers
	opts.ExposedHeaders = convert(opts.ExposedHeaders, http.CanonicalHeaderKey)
	// allowed headers
	if len(opts.AllowedHeaders) == 0 {
		opts.AllowedHeaders = []string{
			HeaderOrigin,
			HeaderAccept,
			HeaderContentType,
			HeaderCustomRequestedWith,
		}
	} else {
		// Origin is always appended as some browsers will always request for this header at preflight
		opts.AllowedHeaders = convert(append(opts.AllowedHeaders, HeaderOrigin), http.CanonicalHeaderKey)
	}
	// allowed origins
	if len(opts.AllowedOrigins) == 0 {
		opts.AllowedOrigins = []string{"*"}
	}
	// allowed methods
	if len(opts.AllowedMethods) == 0 {
		opts.AllowedMethods = []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodHead,
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions && r.Header.Get(CORSHeaderAccessControlRequestMethod) != "" {
				opts.handlePreflight(w, r)

				w.WriteHeader(http.StatusNoContent)
				return
			}
			opts.handleActualRequest(w, r)

			next.ServeHTTP(w, r)
		})
	}
}
