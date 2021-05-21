// Package status provides ...
package status

import (
	"errors"
	"fmt"
	"testing"

	"github.com/deepzz0/go-van/pkg/codes"
)

var testCode codes.Code = 2000

func init() {
	trans := codes.DefaultTranslator{
		Code2Desc: map[string]map[codes.Code]string{
			codes.LangEnUS: {
				testCode: "test: %s",
			},
		},
	}
	codes.WithTranslator(trans)
}

func TestNew(t *testing.T) {
	var cc = struct {
		CC string
		BB string
		DD map[string]string
	}{
		CC: "111",
		BB: "222",
		DD: map[string]string{"key": "value"},
	}
	status := New(codes.FailedPrecondition, "hello", "world", cc)
	t.Log(status.Err())
}

func TestErr(t *testing.T) {
	err := Err(codes.FailedPrecondition, "hello", "world")
	t.Log(err)
}

func TestCode(t *testing.T) {
	err := Err(codes.FailedPrecondition, "hello", "world")
	code := Code(err)
	t.Log(code == codes.FailedPrecondition)
}

func TestCustom(t *testing.T) {
	err := Err(testCode, "hahaha")
	t.Log(err)
}

func TestWrapErr(t *testing.T) {
	err := Err(codes.Unknown)
	err = fmt.Errorf("wrap %w", err)
	err = fmt.Errorf("wrap2 %w", err)
	t.Log(err)

	code := Code(err)
	t.Log(code == codes.Unknown)
}

func TestUnwrapErr(t *testing.T) {
	err := Err(codes.Unknown)
	err = fmt.Errorf("wrap %w", err)

	err = errors.Unwrap(err)
	t.Log(err)
	code := Code(err)
	t.Log(code == codes.Unknown)
}
