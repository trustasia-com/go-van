// Package http provides ...
package http

import "net/http"

type handlerChain []http.Handler
