// Package files provides ...
package files

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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

func TestLoadFilesExpandsEnvironmentVariables(t *testing.T) {
	t.Setenv("CONFX_HOST", "database.internal")
	t.Setenv("CONFX_WORKERS", "8")

	dir := t.TempDir()
	data := []byte(`
dsn: "postgres://${CONFX_HOST}:${CONFX_PORT:-5432}/app"
workers: ${CONFX_WORKERS:-4}
enabled: ${CONFX_ENABLED:-true}
`)
	if err := os.WriteFile(filepath.Join(dir, "app.yml"), data, 0o600); err != nil {
		t.Fatal(err)
	}

	var got struct {
		DSN     string
		Workers int
		Enabled bool
	}
	if err := NewLoader(dir).LoadFiles(&got, "app.yml"); err != nil {
		t.Fatal(err)
	}
	if got.DSN != "postgres://database.internal:5432/app" {
		t.Fatalf("DSN = %q", got.DSN)
	}
	if got.Workers != 8 {
		t.Fatalf("Workers = %d", got.Workers)
	}
	if !got.Enabled {
		t.Fatal("Enabled = false")
	}
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
