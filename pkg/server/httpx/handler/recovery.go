// Package handler provides ...
package handler

import (
	"net/http"
	"runtime"

	"github.com/trustasia-com/go-van/pkg/logx"
)

// RecoverHandler returns a middleware that recover the request.
func RecoverHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				logx.Errorf("[Recovery]%v: %+v\n%s\n", e, r, buf)

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
