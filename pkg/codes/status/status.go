// Package status provides ...
package status

import (
	"context"

	"github.com/trustasia-com/go-van/pkg/codes"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	gcodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// New returns a Status representing c and msg.
func New(c codes.Code, args ...any) *status.Status {
	scode := gcodes.Code(c)
	return status.New(scode, c.Tr(codes.LangEnUS, args...))
}

// Err returns an error representing c and msg.  If c is OK, returns nil.
func Err(c codes.Code, args ...any) error {
	return New(c, args...).Err()
}

// ErrorProto returns an error representing s.  If s.Code is OK, returns nil.
func ErrorProto(s *spb.Status) error {
	return status.FromProto(s).Err()
}

// FromProto returns a Status representing s.
func FromProto(s *spb.Status) *status.Status {
	return status.FromProto(s)
}

// FromError returns a Status representing err if it was produced by this
// package or has a method `GRPCStatus() *Status`.
// If err is nil, a Status is returned with codes.OK and no message.
// Otherwise, ok is false and a Status is returned with codes.Unknown and
// the original error message.
func FromError(err error) (s *status.Status, ok bool) {
	if err == nil {
		return nil, true
	}
	if se, ok := err.(interface {
		GRPCStatus() *status.Status
	}); ok {
		return se.GRPCStatus(), true
	}
	return New(codes.Unknown, err.Error()), false
}

// Convert is a convenience function which removes the need to handle the
// boolean return value from FromError.
func Convert(err error) *status.Status {
	s, _ := FromError(err)
	return s
}

// Code returns the Code of the error if it is a Status error, codes.OK if err
// is nil, or codes.Unknown otherwise.
func Code(err error) codes.Code {
	// Don't use FromError to avoid allocation of OK status.
	if err == nil {
		return codes.OK
	}
	if se, ok := err.(interface {
		GRPCStatus() *status.Status
	}); ok {
		return codes.Code(se.GRPCStatus().Code())
	}
	return codes.Unknown
}

// FromContextError converts a context error into a Status.  It returns a
// Status with codes.OK if err is nil, or a Status with codes.Unknown if err is
// non-nil and not a context error.
func FromContextError(err error) *status.Status {
	switch err {
	case nil:
		return nil
	case context.DeadlineExceeded:
		return New(codes.DeadlineExceeded, err.Error())
	case context.Canceled:
		return New(codes.Canceled, err.Error())
	default:
		return New(codes.Unknown, err.Error())
	}
}
