package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type failingFlushWriter struct {
	http.ResponseWriter
	err error
}

func (w failingFlushWriter) FlushError() error { return w.err }

func TestWrappedWriterPreservesFlushErrors(t *testing.T) {
	want := errors.New("socket flush failed")
	w := &WrappedWriter{ResponseWriter: failingFlushWriter{httptest.NewRecorder(), want}}
	if err := http.NewResponseController(w).Flush(); !errors.Is(err, want) {
		t.Fatalf("Flush() = %v, want %v", err, want)
	}
}

func TestWrappedWriterUnsupportedHijackReturnsError(t *testing.T) {
	w := &WrappedWriter{ResponseWriter: httptest.NewRecorder()}
	_, _, err := http.NewResponseController(w).Hijack()
	if !errors.Is(err, http.ErrNotSupported) {
		t.Fatalf("Hijack() = %v, want ErrNotSupported", err)
	}
}
