// Package apollo provides ...
package apollo

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/trustasia-com/go-van/pkg/codec/yaml"
	"github.com/trustasia-com/go-van/pkg/confx"
	"github.com/trustasia-com/go-van/pkg/logx"

	"github.com/pkg/errors"
	"github.com/zouyx/agollo/v4"
	"github.com/zouyx/agollo/v4/constant"
	"github.com/zouyx/agollo/v4/env/config"
	"github.com/zouyx/agollo/v4/extension"
	"github.com/zouyx/agollo/v4/storage"
)

type apolloLoader struct {
	options Options
	client  *agollo.Client
}

// New a apolloLoader
func NewApolloLoader(opt ...Option) (*apolloLoader, error) {
	// apply opts
	opts := Options{}
	for _, o := range opt {
		o(&opts)
	}

	names := strings.Join(opts.NamespaceNames, ",")
	conf := &config.AppConfig{
		AppID:         opts.AppId,
		Cluster:       opts.Cluster,
		IP:            opts.Addr,
		NamespaceName: names,
		Secret:        opts.Secret,
	}

	for _, namespaceName := range opts.NamespaceNames {
		suffix := filepath.Ext(namespaceName)
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
		return conf, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "StartWithConfig fail")
	}

	return &apolloLoader{
		options: opts,
		client:  client,
	}, nil
}

// LoadFiles load configs from backend
func (l *apolloLoader) LoadFiles(obj interface{}, namespaceName ...string) error {
	if l.client == nil {
		return errors.New("apolloLoader client is nil")
	}
	if len(namespaceName) == 0 {
		return errors.New("The number of namespaceName cannot be 0")
	}

	buff := new(bytes.Buffer)
	for _, name := range namespaceName {
		suffix := filepath.Ext(name)
		if !(suffix == ".yaml" || suffix == ".yml") {
			return errors.New("unsupported file suffix: " + suffix)
		}

		if !in(name, l.options.GetNamespaceNames()) {
			return errors.New(fmt.Sprintf("namespaceName %s not loading", name))
		}

		s := l.client.GetConfigAndInit(name)
		content := s.GetValue("content")
		if content == "" {
			return errors.New(fmt.Sprintf("namespacename %s content is empty", name))
		}
		buff.WriteString(content)
	}

	data := buff.Bytes()
	logx.Infof("all namespaceName content is : %s", data)
	c := yaml.NewCodec()
	err := c.Unmarshal(data, obj)
	if err != nil {
		return errors.Wrap(err, "unmarshal fail")
	}
	return nil
}

// WatchFiles watch file change
func (l *apolloLoader) WatchFiles(ctx context.Context, do confx.WatchFunc, namespaceName ...string) error {
	if l.client == nil {
		return errors.New("apolloLoader client is nil")
	}

	if do == nil {
		return errors.New("do watchFunc is nil")
	}

	if len(namespaceName) == 0 {
		return errors.New("The number of namespaceName cannot be 0")
	}

	listener := &ConfigChangeListener{
		namespaceNames: namespaceName,
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
	return nil
}

func in(target string, strArray []string) bool {
	for _, element := range strArray {
		if target == element {
			return true
		}
	}
	return false
}

type event struct {
	Name string
	Data []byte
}

type ConfigChangeListener struct {
	// These namespace changes will being listened
	namespaceNames []string
	eventChan      chan event
}

// listener config change
func (c *ConfigChangeListener) OnChange(changeEvent *storage.ChangeEvent) {
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

func (c *ConfigChangeListener) OnNewestChange(event *storage.FullChangeEvent) {
}
