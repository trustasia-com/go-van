// Package interceptor provides ...
package interceptor

import (
	"context"
	"fmt"
	"runtime"

	"github.com/deepzz0/go-van/pkg/codes"
	"github.com/deepzz0/go-van/pkg/codes/status"
	"github.com/deepzz0/go-van/pkg/logx"

	"google.golang.org/grpc"
)

// HandlerFunc recover handler func
type HandlerFunc func(ctx context.Context, p interface{}) error

// Option recovery option
type Option func(*options)

// WithHandler panic handler
func WithHandler(h HandlerFunc) Option {
	return func(opts *options) { opts.handler = h }
}

type options struct {
	handler HandlerFunc
}

// UnaryServerInterceptor returns a new unary server interceptor for panic recovery.
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	options := options{handler: defaultHandler}
	for _, o := range opts {
		o(&options)
	}
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {

		panicked := true

		defer func() {
			if e := recover(); e != nil || panicked {
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				logx.Errorf("[Recovery]%v: %+v\n%s\n", e, req, buf)
				err = options.handler(ctx, e)
			}
		}()

		resp, err := handler(ctx, req)
		panicked = false
		return resp, err
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for panic recovery.
func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	options := options{handler: defaultHandler}
	for _, o := range opts {
		o(&options)
	}
	return func(srv interface{}, stream grpc.ServerStream,
		info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {

		panicked := true

		defer func() {
			if e := recover(); e != nil || panicked {
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				logx.Errorf("[Recovery]%v: %+v\n%s\n", e, info, buf)
				err = options.handler(stream.Context(), e)
			}
		}()

		err = handler(srv, stream)
		panicked = false
		return err
	}
}

func defaultHandler(ctx context.Context, p interface{}) error {
	return status.Err(codes.Internal, fmt.Sprintf("%v", p))
}
