// Package files provides ...
package files

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/trustasia-com/go-van/pkg/confx"
)

var (
	conf struct {
		Database struct {
			Driver string
			Source string
		}
		Ports []int
		Grpc  map[string]string
	}
	loader confx.Confx
)

func init() {
	loader = NewLoader("../testdata/")
}

func TestLoadFiles(t *testing.T) {
	err := loader.LoadFiles(&conf, "test.yml", "test2.yml")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(conf)
}

func TestWatchFiles(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := loader.WatchFiles(ctx, watchFunc, "test.yml", "test2.yml")
	if err != nil {
		t.Fatal(err)
	}
}

func watchFunc(name string, data []byte) error {
	fmt.Println(name)
	fmt.Println(string(data))
	loader.LoadFiles(&conf, name)
	fmt.Println(conf)
	return nil
}
