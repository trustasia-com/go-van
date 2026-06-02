// Package handler provides ...
package handler

import (
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/trustasia-com/go-van/pkg/logx"
)

//
// the middleware from: https://github.com/rs/cors
//

var (
	headerVaryOrigin = []string{"Origin"}
	headerOriginAll  = []string{"*"}
	headerTrue       = []string{"true"}
)

// CORSOptions is a configuration container to setup the CORS middleware.
type CORSOptions struct {
	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// An origin may contain a wildcard (*) to replace 0 or more characters
	// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penalty.
	// Only one wildcard can be used per origin.
	// Default value is ["*"]
	AllowedOrigins []string
	// AllowOriginFunc is a custom function to validate the origin. It take the origin
	// as argument and returns true if allowed or false otherwise. If this option is
	// set, the content of AllowedOrigins is ignored.
	AllowOriginFunc func(origin string) bool
	// AllowOriginRequestFunc is a custom function to validate the origin. It takes the HTTP Request object and the origin as
	// argument and returns true if allowed or false otherwise. If this option is set, the content of `AllowedOrigins`
	// and `AllowOriginFunc` is ignored.
	AllowOriginRequestFunc func(r *http.Request, origin string) bool
	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods []string
	// AllowedHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders []string
	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposedHeaders []string
	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached. Default value is 0, which stands for no Access-Control-Max-Age
	// header to be sent back. Set MaxAge to a negative value to force a 0 max-age.
	MaxAge int
	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool
	// AllowPrivateNetwork indicates whether to accept cross-origin requests over a
	// private network.
	AllowPrivateNetwork bool
	// OptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	OptionsPassthrough bool
	// Provides a status code to use for successful OPTIONS requests.
	// Default value is http.StatusNoContent (204).
	OptionsSuccessStatus int
	// Debugging flag adds additional output to debug server side CORS issues
	Debug bool
}

// Logger generic interface for logger
type Logger interface {
	Printf(string, ...any)
}

// Cors http handler
type Cors struct {
	// Debug logger
	Debug bool
	// Normalized list of plain allowed origins
	allowedOrigins []string
	// List of allowed origins containing wildcards
	allowedWOrigins []wildcard
	// Optional origin validator function
	allowOriginFunc func(origin string) bool
	// Optional origin validator (with request) function
	allowOriginRequestFunc func(r *http.Request, origin string) bool
	// Normalized set of allowed headers (lowercase)
	allowedHeaders map[string]struct{}
	// Normalized list of allowed methods
	allowedMethods []string
	// Pre-computed normalized list of exposed headers
	exposedHeaders []string
	// Pre-computed maxAge header value
	maxAge []string
	// Set to true when allowed origins contains a "*"
	allowedOriginsAll bool
	// Set to true when allowed headers contains a "*"
	allowedHeadersAll bool
	// Status code to use for successful OPTIONS requests
	optionsSuccessStatus int
	allowCredentials     bool
	allowPrivateNetwork  bool
	optionPassthrough    bool
	preflightVary        []string
}

// New creates a new Cors handler with the provided options.
func New(options CORSOptions) *Cors {
	c := &Cors{
		Debug:                  options.Debug,
		exposedHeaders:         convert(options.ExposedHeaders, http.CanonicalHeaderKey),
		allowOriginFunc:        options.AllowOriginFunc,
		allowOriginRequestFunc: options.AllowOriginRequestFunc,
		allowCredentials:       options.AllowCredentials,
		allowPrivateNetwork:    options.AllowPrivateNetwork,
		optionPassthrough:      options.OptionsPassthrough,
	}

	// Normalize options
	// Note: for origins and methods matching, the spec requires a case-sensitive matching.
	// As it may error prone, we chose to ignore the spec here.

	// Allowed Origins
	if len(options.AllowedOrigins) == 0 {
		if options.AllowOriginFunc == nil && options.AllowOriginRequestFunc == nil {
			// Default is all origins
			c.allowedOriginsAll = true
		}
	} else {
		c.allowedOrigins = []string{}
		c.allowedWOrigins = []wildcard{}
		for _, origin := range options.AllowedOrigins {
			// Normalize
			origin = strings.ToLower(origin)
			if origin == "*" {
				// If "*" is present in the list, turn the whole list into a match all
				c.allowedOriginsAll = true
				c.allowedOrigins = nil
				c.allowedWOrigins = nil
				break
			} else if prefix, suffix, ok := strings.Cut(origin, "*"); ok {
				w := wildcard{prefix, suffix}
				c.allowedWOrigins = append(c.allowedWOrigins, w)
			} else {
				c.allowedOrigins = append(c.allowedOrigins, origin)
			}
		}
	}

	if c.allowCredentials && c.allowedOriginsAll &&
		c.allowOriginFunc == nil && c.allowOriginRequestFunc == nil {
		panic("cors: AllowCredentials cannot be used with AllowedOrigins=[\"*\"] or empty AllowedOrigins; specify explicit origins or use AllowOriginFunc")
	}

	// Allowed Headers
	// Note: the Fetch standard guarantees that CORS-unsafe request-header names are lowercase.
	if len(options.AllowedHeaders) == 0 {
		c.allowedHeaders = newAllowedHeaderSet("accept", "content-type", "origin", "x-requested-with")
	} else {
		normalized := convert(append(options.AllowedHeaders, "Origin"), strings.ToLower)
		if slices.Contains(options.AllowedHeaders, "*") {
			c.allowedHeadersAll = true
		} else {
			c.allowedHeaders = newAllowedHeaderSet(normalized...)
		}
	}

	// Allowed Methods
	if len(options.AllowedMethods) == 0 {
		// Default is spec's "simple" methods
		c.allowedMethods = []string{http.MethodGet, http.MethodPost, http.MethodHead}
	} else {
		c.allowedMethods = convert(options.AllowedMethods, strings.ToUpper)
	}

	// Options Success Status Code
	if options.OptionsSuccessStatus == 0 {
		c.optionsSuccessStatus = http.StatusNoContent
	} else {
		c.optionsSuccessStatus = options.OptionsSuccessStatus
	}

	if c.allowPrivateNetwork {
		c.preflightVary = []string{"Origin, Access-Control-Request-Method, Access-Control-Request-Headers, Access-Control-Request-Private-Network"}
	} else {
		c.preflightVary = []string{"Origin, Access-Control-Request-Method, Access-Control-Request-Headers"}
	}

	if options.MaxAge > 0 {
		c.maxAge = []string{strconv.Itoa(options.MaxAge)}
	} else if options.MaxAge < 0 {
		c.maxAge = []string{"0"}
	}

	return c
}

// Default creates a new Cors handler with default options.
func Default() *Cors {
	return New(CORSOptions{})
}

// AllowAll create a new Cors handler with permissive configuration allowing all
// origins with all standard methods with any header.
func AllowAll() *Cors {
	return New(CORSOptions{
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
	})
}

// Handler apply the CORS specification on the request, and add relevant CORS headers
// as necessary.
func (c *Cors) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			c.logf("Handler: Preflight request")
			c.handlePreflight(w, r)
			// Preflight requests are standalone and should stop the chain as some other
			// middleware may not handle OPTIONS requests correctly. One typical example
			// is authentication middleware ; OPTIONS requests won't carry authentication
			// headers (see #1)
			if c.optionPassthrough {
				h.ServeHTTP(w, r)
			} else {
				w.WriteHeader(c.optionsSuccessStatus)
			}
		} else {
			c.logf("Handler: Actual request")
			c.handleActualRequest(w, r)
			h.ServeHTTP(w, r)
		}
	})
}

// HandlerFunc provides Martini compatible handler
func (c *Cors) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
		c.logf("HandlerFunc: Preflight request")
		c.handlePreflight(w, r)
		if !c.optionPassthrough {
			w.WriteHeader(c.optionsSuccessStatus)
		}
	} else {
		c.logf("HandlerFunc: Actual request")
		c.handleActualRequest(w, r)
	}
}

// Negroni compatible interface
func (c *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
		c.logf("ServeHTTP: Preflight request")
		c.handlePreflight(w, r)
		// Preflight requests are standalone and should stop the chain as some other
		// middleware may not handle OPTIONS requests correctly. One typical example
		// is authentication middleware ; OPTIONS requests won't carry authentication
		// headers (see #1)
		if c.optionPassthrough {
			next(w, r)
		} else {
			w.WriteHeader(c.optionsSuccessStatus)
		}
	} else {
		c.logf("ServeHTTP: Actual request")
		c.handleActualRequest(w, r)
		next(w, r)
	}
}

// handlePreflight handles pre-flight CORS requests
func (c *Cors) handlePreflight(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	if r.Method != http.MethodOptions {
		c.logf("  Preflight aborted: %s!=OPTIONS", r.Method)
		return
	}
	// Always set Vary headers
	// see https://github.com/rs/cors/issues/10,
	//     https://github.com/rs/cors/commit/dbdca4d95feaa7511a46e6f1efb3b3aa505bc43f#commitcomment-12352001
	if vary, found := headers["Vary"]; found {
		headers["Vary"] = append(vary, c.preflightVary[0])
	} else {
		headers["Vary"] = c.preflightVary
	}

	if origin == "" {
		c.logf("  Preflight aborted: empty origin")
		return
	}
	if !c.isOriginAllowed(r, origin) {
		c.logf("  Preflight aborted: origin '%s' not allowed", origin)
		return
	}

	reqMethod := r.Header.Get("Access-Control-Request-Method")
	if !c.isMethodAllowed(reqMethod) {
		c.logf("  Preflight aborted: method '%s' not allowed", reqMethod)
		return
	}
	// Note: the Fetch standard guarantees at most one Access-Control-Request-Headers
	// header in preflight requests, but some gateways split it into multiple headers.
	reqHeaders, found := r.Header["Access-Control-Request-Headers"]
	if found && !c.allowedHeadersAll && !c.areRequestHeadersAllowed(reqHeaders) {
		c.logf("  Preflight aborted: headers '%v' not allowed", reqHeaders)
		return
	}
	if c.allowedOriginsAll {
		headers["Access-Control-Allow-Origin"] = headerOriginAll
	} else {
		headers["Access-Control-Allow-Origin"] = r.Header["Origin"]
	}
	// Spec says: Since the list of methods can be unbounded, simply returning the method indicated
	// by Access-Control-Request-Method (if supported) can be enough
	headers["Access-Control-Allow-Methods"] = r.Header["Access-Control-Request-Method"]
	if found && len(reqHeaders[0]) > 0 {
		// Spec says: Since the list of headers can be unbounded, simply returning supported headers
		// from Access-Control-Request-Headers can be enough
		headers["Access-Control-Allow-Headers"] = reqHeaders
	}
	if c.allowCredentials {
		headers["Access-Control-Allow-Credentials"] = headerTrue
	}
	if c.allowPrivateNetwork && r.Header.Get("Access-Control-Request-Private-Network") == "true" {
		headers["Access-Control-Allow-Private-Network"] = headerTrue
	}
	if len(c.maxAge) > 0 {
		headers["Access-Control-Max-Age"] = c.maxAge
	}
	c.logf("  Preflight response headers: %v", headers)
}

// handleActualRequest handles simple cross-origin requests, actual request or redirects
func (c *Cors) handleActualRequest(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	// Always set Vary, see https://github.com/rs/cors/issues/10
	if vary := headers["Vary"]; vary == nil {
		headers["Vary"] = headerVaryOrigin
	} else {
		headers["Vary"] = append(vary, headerVaryOrigin[0])
	}
	if origin == "" {
		c.logf("  Actual request no headers added: missing origin")
		return
	}
	if !c.isOriginAllowed(r, origin) {
		c.logf("  Actual request no headers added: origin '%s' not allowed", origin)
		return
	}

	// Note that spec does define a way to specifically disallow a simple method like GET or
	// POST. Access-Control-Allow-Methods is only used for pre-flight requests and the
	// spec doesn't instruct to check the allowed methods for simple cross-origin requests.
	// We think it's a nice feature to be able to have control on those methods though.
	if !c.isMethodAllowed(r.Method) {
		c.logf("  Actual request no headers added: method '%s' not allowed", r.Method)

		return
	}
	if c.allowedOriginsAll {
		headers["Access-Control-Allow-Origin"] = headerOriginAll
	} else {
		headers["Access-Control-Allow-Origin"] = r.Header["Origin"]
	}
	if len(c.exposedHeaders) > 0 {
		headers.Set("Access-Control-Expose-Headers", strings.Join(c.exposedHeaders, ", "))
	}
	if c.allowCredentials {
		headers["Access-Control-Allow-Credentials"] = headerTrue
	}
	c.logf("  Actual response added headers: %v", headers)
}

// convenience method. checks if a logger is set.
func (c *Cors) logf(format string, a ...any) {
	if c.Debug {
		logx.Infof(format, a...)
	}
}

// OriginAllowed reports whether the request origin is allowed.
func (c *Cors) OriginAllowed(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return false
	}
	return c.isOriginAllowed(r, origin)
}

// isOriginAllowed checks if a given origin is allowed to perform cross-domain requests
// on the endpoint
func (c *Cors) isOriginAllowed(r *http.Request, origin string) bool {
	if c.allowOriginRequestFunc != nil {
		return c.allowOriginRequestFunc(r, origin)
	}
	if c.allowOriginFunc != nil {
		return c.allowOriginFunc(origin)
	}
	if c.allowedOriginsAll {
		return true
	}
	origin = strings.ToLower(origin)
	if slices.Contains(c.allowedOrigins, origin) {
		return true
	}
	return slices.ContainsFunc(c.allowedWOrigins, func(w wildcard) bool {
		return w.match(origin)
	})
}

// isMethodAllowed checks if a given method can be used as part of a cross-domain request
// on the endpoint
func (c *Cors) isMethodAllowed(method string) bool {
	if len(c.allowedMethods) == 0 {
		// If no method allowed, always return false, even for preflight request
		return false
	}
	method = strings.ToUpper(method)
	if method == http.MethodOptions {
		// Always allow preflight requests
		return true
	}
	return slices.Contains(c.allowedMethods, method)
}

func (c *Cors) areRequestHeadersAllowed(values []string) bool {
	for _, value := range values {
		for _, name := range parseRequestHeaderNames(value) {
			if name == "" {
				continue
			}
			if _, ok := c.allowedHeaders[name]; !ok {
				return false
			}
		}
	}
	return true
}

type converter func(string) string

type wildcard struct {
	prefix string
	suffix string
}

func (w wildcard) match(s string) bool {
	return len(s) >= len(w.prefix)+len(w.suffix) && strings.HasPrefix(s, w.prefix) && strings.HasSuffix(s, w.suffix)
}

func newAllowedHeaderSet(headers ...string) map[string]struct{} {
	set := make(map[string]struct{}, len(headers))
	for _, header := range headers {
		set[header] = struct{}{}
	}
	return set
}

// convert converts a list of string using the passed converter function
func convert(s []string, c converter) []string {
	out := make([]string, 0, len(s))
	for _, i := range s {
		out = append(out, c(i))
	}
	return out
}

func parseRequestHeaderNames(headerList string) []string {
	if headerList == "" {
		return nil
	}
	parts := strings.Split(headerList, ",")
	names := make([]string, 0, len(parts))
	for _, part := range parts {
		name := strings.ToLower(strings.TrimSpace(part))
		if name != "" {
			names = append(names, name)
		}
	}
	return names
}
