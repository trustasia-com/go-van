// Package apollo provides ...
package apollo

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/trustasia-com/go-van/pkg/codec/yaml"
	"github.com/trustasia-com/go-van/pkg/confx"
	"github.com/trustasia-com/go-van/pkg/logx"
	"github.com/zouyx/agollo/v4"
	"github.com/zouyx/agollo/v4/constant"
	"github.com/zouyx/agollo/v4/env/config"
	"github.com/zouyx/agollo/v4/extension"
	"github.com/zouyx/agollo/v4/storage"
	"strings"
	"sync"
)

type apolloLoader struct {
	conf   *config.AppConfig
	client *agollo.Client
}

// create a apolloLoader
func NewApolloLoader(appId, cluster, ip, namespaceNames, secret string) (*apolloLoader, error) {
	conf := &config.AppConfig{
		AppID:          appId,
		Cluster:        cluster,
		IP:             ip,
		NamespaceName:  namespaceNames,
		IsBackupConfig: false,
		Secret:         secret,
	}

	agollo.SetBackupFileHandler(&FileHandler{})

	str := strings.Split(namespaceNames, ",")
	for _, namespaceName := range str {

		t := strings.Split(namespaceName, ".")
		if len(t) < 2 {
			return nil, errors.New("namespaceName error")
		}
		suffix := "." + t[1]
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
		conf:   conf,
		client: client,
	}, nil
}

// LoadFiles load config from backend
func (l *apolloLoader) LoadFiles(obj interface{}, namespaceName string) error {
	if l.client == nil {
		return errors.New("apolloLoader  client is nil")
	}

	s := l.client.GetConfig(namespaceName)
	content := s.GetValue("content")
	if content == "" {
		return errors.New("content is empty")
	}
	logx.Infof("namespaceName %s content is : %s", namespaceName, content)

	t := strings.SplitAfter(namespaceName, ".")
	if len(t) < 2 {
		return errors.New("namespaceName error")
	}
	suffix := "." + t[1]

	var err error
	data := []byte(content)
	switch suffix {
	case ".yaml", ".yml":
		err = yaml.NewCodec().Unmarshal(data, obj)
	default:
		return errors.New("unsupported file suffix: " + suffix)
	}
	if err != nil {
		return errors.Wrap(err, namespaceName)
	}
	return nil
}

// WatchFiles watch file change
func (l *apolloLoader) WatchFiles(do confx.WatchFunc, namespaceName string) error {
	if l.client == nil {
		return errors.New("apolloLoader  client is nil")
	}

	if do == nil {
		return errors.New("do watchFunc is nil")
	}

	listener := &CustomChangeListener{
		doFunc: do,
	}

	listener.wg.Add(1)

	l.client.AddChangeListener(listener)

	listener.wg.Wait()
	return nil
}

type FileHandler struct {
}

// WriteConfigFile write config to file
func (fileHandler *FileHandler) WriteConfigFile(config *config.ApolloConfig, configPath string) error {
	//fmt.Println(config.Configurations)
	return nil
}

// GetConfigFile get real config file
func (fileHandler *FileHandler) GetConfigFile(configDir string, appID string, namespace string) string {
	return ""
}

//LoadConfigFile load config from file
func (fileHandler *FileHandler) LoadConfigFile(configDir string, appID string, namespace string) (*config.ApolloConfig, error) {
	return &config.ApolloConfig{}, nil
}

type CustomChangeListener struct {
	wg     sync.WaitGroup
	doFunc confx.WatchFunc
}

// listener config change
func (c *CustomChangeListener) OnChange(changeEvent *storage.ChangeEvent) {
	logx.Infof("change event: %+v", changeEvent.Changes)
	content := changeEvent.Changes["content"]
	logx.Infof("change value:%+v", content)
	data := []byte(fmt.Sprintf("%+v", content.NewValue))
	c.doFunc(changeEvent.Namespace, data)
}

func (c *CustomChangeListener) OnNewestChange(event *storage.FullChangeEvent) {
	//write your code here
}
