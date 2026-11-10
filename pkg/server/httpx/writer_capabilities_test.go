package httpx

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/trustasia-com/go-van/pkg/server"
)

func TestServerPreservesStreamingControl(t *testing.T) {
	release := make(chan struct{})
	controlErrors := make(chan error, 3)
	h := NewServer(server.WithHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		controller := http.NewResponseController(w)
		for _, setDeadline := range []func(time.Time) error{controller.SetReadDeadline, controller.SetWriteDeadline} {
			err := setDeadline(time.Now().Add(3 * time.Second))
			controlErrors <- err
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "data: first\n\n")
		controlErrors <- controller.Flush()
		<-release
		_, _ = io.WriteString(w, "data: last\n\n")
	})))
	s := httptest.NewServer(h)
	defer s.Close()
	defer close(release)
	client := &http.Client{Timeout: 2 * time.Second}
	response, err := client.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("stream status = %d, deadline error = %v", response.StatusCode, <-controlErrors)
	}
	line, err := bufio.NewReader(response.Body).ReadString('\n')
	if err != nil || line != "data: first\n" {
		t.Fatalf("first event = %q, error = %v", line, err)
	}
	for range 3 {
		if err := <-controlErrors; err != nil {
			t.Fatalf("stream control error = %v", err)
		}
	}
}

func TestServerPreservesHijacking(t *testing.T) {
	h := NewServer(server.WithHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hijacker, ok := w.(http.Hijacker)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		conn, rw, err := hijacker.Hijack()
		if err != nil {
			return
		}
		defer func() { _ = conn.Close() }()
		_, _ = rw.WriteString("HTTP/1.1 101 Switching Protocols\r\nConnection: Upgrade\r\nUpgrade: test\r\n\r\n")
		_ = rw.Flush()
		line, err := rw.ReadString('\n')
		if err == nil {
			_, _ = rw.WriteString(line)
			_ = rw.Flush()
		}
	})))
	s := httptest.NewServer(h)
	defer s.Close()
	conn, err := net.DialTimeout("tcp", strings.TrimPrefix(s.URL, "http://"), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	if err := conn.SetDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatal(err)
	}
	_, _ = io.WriteString(conn, "GET / HTTP/1.1\r\nHost: localhost\r\nConnection: Upgrade\r\nUpgrade: test\r\n\r\n")
	reader := bufio.NewReader(conn)
	response, err := http.ReadResponse(reader, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("upgrade status = %d, want 101", response.StatusCode)
	}
	_, _ = io.WriteString(conn, "echo\n")
	line, err := reader.ReadString('\n')
	if err != nil || line != "echo\n" {
		t.Fatalf("hijacked echo = %q, error = %v", line, err)
	}
}
