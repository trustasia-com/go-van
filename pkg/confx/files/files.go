// Package files provides ...
package files

import (
	"bytes"
	"context"
	"os"
	"path/filepath"

	"github.com/trustasia-com/go-van/pkg/codec/yaml"
	"github.com/trustasia-com/go-van/pkg/confx"
	"github.com/trustasia-com/go-van/pkg/logx"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

type filesLoader struct {
	filepath string
}

func NewFilesLoader(filepath string) *filesLoader {
	return &filesLoader{
		filepath: filepath,
	}
}

// LoadFiles load config from backend
func (l *filesLoader) LoadFiles(obj interface{}, fileNames ...string) error {
	if l.filepath == "" {
		return errors.New("filepath not set")
	}

	buff := new(bytes.Buffer)
	for _, name := range fileNames {
		suffix := filepath.Ext(name)
		if !(suffix == ".yaml" || suffix == ".yml") {
			return errors.New("unsupported file suffix: " + suffix)
		}

		file := filepath.Join(l.filepath, name)

		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		buff.Write(data)
	}

	data := buff.Bytes()
	c := yaml.NewCodec()
	err := c.Unmarshal(data, obj)
	if err != nil {
		logx.Error("error:", err)
		return errors.Wrap(err, "unmarshal fail")
	}
	return nil
}

// WatchFiles watch file change
func (l *filesLoader) WatchFiles(ctx context.Context, do confx.WatchFunc, fileNames ...string) error {
	if l.filepath == "" {
		return errors.New("filepath not set")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// watch file
	for _, fileName := range fileNames {
		file := filepath.Join(l.filepath, fileName)
		err = watcher.Add(file)
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
	return nil
}
