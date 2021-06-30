// Package files provides ...
package files

import (
	"context"
	"os"
	"path/filepath"

	"github.com/trustasia-com/go-van/pkg/codec/yaml"
	"github.com/trustasia-com/go-van/pkg/confx"
	"github.com/trustasia-com/go-van/pkg/logx"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

type filesLoader struct{}

// LoadFiles load config from backend
func (l *filesLoader) LoadFiles(obj interface{}, glob string) error {
	files, err := filepath.Glob(glob)
	if err != nil {
		return err
	}

	ext := filepath.Ext(glob)
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		switch ext {
		case ".yaml", ".yml":
			c := yaml.NewCodec()
			err = c.Unmarshal(data, obj)
		default:
			return errors.New("unsupported file ext: " + ext)
		}
		if err != nil {
			return errors.Wrap(err, file)
		}
	}
	return nil
}

// WatchFiles watch file change
func (l *filesLoader) WatchFiles(ctx context.Context, do confx.WatchFunc, glob string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	// watch file
	dir := filepath.Dir(glob)
	err = watcher.Add(dir)
	if err != nil {
		return err
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
				return err
			}
			logx.Error("error:", err)
		case <-ctx.Done():
			logx.Error("error:", ctx.Err())
			return ctx.Err()
		}
	}
	return nil
}
