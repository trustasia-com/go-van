// Package van provides ...
package van

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type testServer struct {
	done chan struct{}
}

func (s *testServer) Start() error {
	go func() {
		for {
			fmt.Println(1)
			time.Sleep(time.Second)
		}
	}()
	<-s.done
	return nil
}

func (s *testServer) Stop() error {
	s.done <- struct{}{}
	return nil
}

func (s *testServer) Endpoint() (string, error) {
	return "", nil
}

func TestService(t *testing.T) {
	svr := &testServer{done: make(chan struct{})}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 5)
		cancel()
	}()
	srv := NewSrv(Server(svr), Context(ctx))

	if err := srv.Run(); err != nil {
		t.Fatal(err)
	}
}
