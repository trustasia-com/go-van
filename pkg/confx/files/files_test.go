// Package files provides ...
package files

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"
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
	loader filesLoader
)

func TestLoadFiles(t *testing.T) {
	err := loader.LoadFiles(&conf, "../testdata/test.yml")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(conf)
}

func TestWatchFiles(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	err := loader.WatchFiles(ctx, watchFunc, "../testdata/test.yml")
	if err != nil {
		t.Fatal(err)
	}
}

func watchFunc(name string, data []byte) error {
	fmt.Println(name)
	fmt.Println(string(data))
	_, file := filepath.Split(name)
	switch file {
	case "test.yml":
		loader.LoadFiles(&conf, name)
	default:
	}
	fmt.Println(conf)
	return nil
}
