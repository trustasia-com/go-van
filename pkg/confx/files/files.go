// Package files provides ...
package files

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/trustasia-com/go-van/pkg/codec/yaml"
	"github.com/trustasia-com/go-van/pkg/confx"
	"github.com/trustasia-com/go-van/pkg/logx"

	"github.com/fsnotify/fsnotify"
)

type filesLoader struct {
	filesDir string
}

// NewLoader files loader instance
func NewLoader(dir string) confx.Confx {
	return &filesLoader{
		filesDir: dir,
	}
}

// LoadFiles load config from backend
func (l *filesLoader) LoadFiles(obj any, files ...string) error {
	buf := new(bytes.Buffer)
	for _, name := range files {
		suffix := filepath.Ext(name)
		if !(suffix == ".yaml" || suffix == ".yml") {
			return errors.New("unsupported file suffix: " + suffix)
		}

		path := filepath.Join(l.filesDir, name)

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		buf.Write(data)
		buf.WriteByte('\n')
	}

	data := buf.Bytes()
	c := yaml.NewCodec()
	err := c.Unmarshal(data, obj)
	if err != nil {
		return fmt.Errorf("unmarshal fail: %w", err)
	}
	return nil
}

// WatchFiles watch file change
func (l *filesLoader) WatchFiles(ctx context.Context, do confx.WatchFunc, fileNames ...string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// watch file
	for _, fileName := range fileNames {
		path := filepath.Join(l.filesDir, fileName)
		err = watcher.Add(path)
		if err != nil {
			return err
		}
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				data, err := os.ReadFile(event.Name)
				if err != nil {
					return err
				}
				do(event.Name, data)
			}
			// Remove Create Rename Chmod
		case err, ok := <-watcher.Errors:
			if !ok {
				logx.Error("error:", err)
				return err
			}

		case <-ctx.Done():
			logx.Error("error:", ctx.Err())
			return ctx.Err()
		}
	}
}
