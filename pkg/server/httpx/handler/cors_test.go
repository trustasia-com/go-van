package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewPanicsOnCredentialsWithWildcardOrigins(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for AllowCredentials with wildcard origins")
		}
	}()
	New(CORSOptions{
		AllowCredentials: true,
	})
}

func TestNewAllowsCredentialsWithExplicitOrigins(t *testing.T) {
	c := New(CORSOptions{
		AllowedOrigins:   []string{"https://example.com"},
		AllowCredentials: true,
	})
	if c == nil {
		t.Fatal("expected cors handler")
	}
}

func TestNewAllowsCredentialsWithOriginFunc(t *testing.T) {
	c := New(CORSOptions{
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://example.com"
		},
	})
	if c == nil {
		t.Fatal("expected cors handler")
	}
}

func TestPreflightAllowsRequestedHeaders(t *testing.T) {
	c := New(CORSOptions{
		AllowedOrigins: []string{"https://example.com"},
		AllowedMethods: []string{http.MethodGet},
		AllowedHeaders: []string{"X-Custom-Token"},
	})

	req := httptest.NewRequest(http.MethodOptions, "http://api.example.com/resource", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "x-custom-token")
	rec := httptest.NewRecorder()

	c.handlePreflight(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Headers"); got != "x-custom-token" {
		t.Fatalf("Access-Control-Allow-Headers = %q, want %q", got, "x-custom-token")
	}
}

func TestPreflightAllowsSplitRequestHeaders(t *testing.T) {
	c := New(CORSOptions{
		AllowedOrigins: []string{"https://example.com"},
		AllowedMethods: []string{http.MethodGet},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	req := httptest.NewRequest(http.MethodOptions, "http://api.example.com/resource", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header["Access-Control-Request-Headers"] = []string{"authorization", "content-type"}

	rec := httptest.NewRecorder()
	c.handlePreflight(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://example.com" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, "https://example.com")
	}
}

func TestPreflightRejectsDisallowedHeader(t *testing.T) {
	c := New(CORSOptions{
		AllowedOrigins: []string{"https://example.com"},
		AllowedMethods: []string{http.MethodGet},
		AllowedHeaders: []string{"Authorization"},
	})

	req := httptest.NewRequest(http.MethodOptions, "http://api.example.com/resource", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "x-denied")
	rec := httptest.NewRecorder()

	c.handlePreflight(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want empty", got)
	}
}

func TestActualRequestWithCredentialsEchoesOrigin(t *testing.T) {
	c := New(CORSOptions{
		AllowedOrigins:   []string{"https://example.com"},
		AllowedMethods:   []string{http.MethodGet},
		AllowCredentials: true,
	})

	req := httptest.NewRequest(http.MethodGet, "http://api.example.com/resource", nil)
	req.Header.Set("Origin", "https://example.com")
	rec := httptest.NewRecorder()

	c.handleActualRequest(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://example.com" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, "https://example.com")
	}
	if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("Access-Control-Allow-Credentials = %q, want true", got)
	}
}

func TestPreflightVaryMerged(t *testing.T) {
	c := Default()
	req := httptest.NewRequest(http.MethodOptions, "http://api.example.com/resource", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	rec := httptest.NewRecorder()
	rec.Header().Set("Vary", "Accept-Encoding")

	c.handlePreflight(rec, req)

	vary := rec.Header().Values("Vary")
	if len(vary) != 2 {
		t.Fatalf("Vary values = %v, want 2 entries", vary)
	}
	if !strings.Contains(vary[1], "Access-Control-Request-Method") {
		t.Fatalf("second Vary entry = %q, want preflight vary fields", vary[1])
	}
}

func TestOriginAllowedRequiresOriginHeader(t *testing.T) {
	c := Default()
	req := httptest.NewRequest(http.MethodGet, "http://api.example.com/resource", nil)
	if c.OriginAllowed(req) {
		t.Fatal("expected empty origin to be disallowed")
	}
}
