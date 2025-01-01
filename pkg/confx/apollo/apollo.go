// Package apollo provides ...
package apollo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/trustasia-com/go-van/pkg/codec/yaml"
	"github.com/trustasia-com/go-van/pkg/confx"
	"github.com/trustasia-com/go-van/pkg/logx"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/constant"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/extension"
	"github.com/apolloconfig/agollo/v4/storage"
)

type apolloLoader struct {
	client agollo.Client
}

// NewLoader a apolloLoader
func NewLoader(opts ...Option) (confx.Confx, error) {
	// apply opts
	options := Options{}
	for _, o := range opts {
		o(&options)
	}

	names := strings.Split(options.NamespaceName, ",")
	for _, ns := range names {
		suffix := filepath.Ext(ns)
		switch constant.ConfigFileFormat(suffix) {
		case constant.YML:
			extension.AddFormatParser(constant.YML, nil)
		case constant.YAML:
			extension.AddFormatParser(constant.YAML, nil)
		default:
			return nil, errors.New("unsupported file suffix: " + suffix)
		}
	}

	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return &options, nil
	})
	if err != nil {
		return nil, fmt.Errorf("StartWithConfig fail %w", err)
	}

	return &apolloLoader{client: client}, nil
}

// LoadFiles load configs from backend
func (l *apolloLoader) LoadFiles(obj interface{}, namespaces ...string) error {
	if len(namespaces) == 0 {
		return errors.New("please specific need load namespace")
	}

	buf := new(bytes.Buffer)
	for _, name := range namespaces {
		suffix := filepath.Ext(name)
		if !(suffix == ".yaml" || suffix == ".yml") {
			return errors.New("unsupported file suffix: " + suffix)
		}

		s := l.client.GetConfigAndInit(name)
		if s == nil || !s.GetIsInit() {
			return errors.New("namespace not init with NewLoader")
		}
		content := s.GetValue("content")
		if content != "" {
			buf.WriteString(content)
			buf.WriteByte('\n')
		}
	}

	data := buf.Bytes()
	c := yaml.NewCodec()
	err := c.Unmarshal(data, obj)
	if err != nil {
		return fmt.Errorf("unmarshal fail %w", err)
	}
	return nil
}

// WatchFiles watch file change
func (l *apolloLoader) WatchFiles(ctx context.Context, do confx.WatchFunc, namespaces ...string) error {
	if l.client == nil {
		return errors.New("apolloLoader client is nil")
	}

	if do == nil {
		return errors.New("do watchFunc is nil")
	}

	if len(namespaces) == 0 {
		return errors.New("The number of namespaceName cannot be 0")
	}

	listener := &configChangeListener{
		namespaceNames: namespaces,
		eventChan:      make(chan event),
	}

	l.client.AddChangeListener(listener)
	for {
		select {
		case ev, ok := <-listener.eventChan:
			if !ok {
				return nil
			}
			do(ev.Name, ev.Data)
		case <-ctx.Done():
			logx.Error("error:", ctx.Err())
			return ctx.Err()
		}
	}
}

type event struct {
	Name string
	Data []byte
}

type configChangeListener struct {
	// These namespace changes will being listened
	namespaceNames []string
	eventChan      chan event
}

// listener config change
func (c *configChangeListener) OnChange(changeEvent *storage.ChangeEvent) {
	logx.Infof("change namespace: %s", changeEvent.Namespace)
	changeNamespace := changeEvent.Namespace
	for _, space := range c.namespaceNames {
		if space == changeNamespace {
			content := changeEvent.Changes["content"]
			logx.Infof("change value:%+v", content)
			dataByte := []byte(fmt.Sprintf("%+v", content.NewValue))
			c.eventChan <- event{Name: changeNamespace, Data: dataByte}
		}
	}
}

// OnNewestChange implementations listener
func (c *configChangeListener) OnNewestChange(event *storage.FullChangeEvent) {}
