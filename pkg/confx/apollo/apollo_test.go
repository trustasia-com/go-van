// Package apollo provides ...
package apollo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/trustasia-com/go-van/pkg/codec/yaml"
	"github.com/trustasia-com/go-van/pkg/confx"

	"github.com/apolloconfig/agollo/v4/env/config"
)

type Conf struct {
	Database struct {
		Driver string
		Source string
	}
	Ports []int
	Grpc  struct {
		User string
	}
}

var (
	conf   Conf
	loader confx.Confx
)

func init() {
	conf := config.AppConfig{
		AppID:         "SampleApp",
		Cluster:       "dev",
		IP:            "http://192.168.252.177:8080",
		NamespaceName: "app.yml,connect.yml",
		Secret:        "a5ce81b8767e4d4cbc0baf94fea57bfb",
	}
	var err error
	loader, err = NewLoader(
		WithConfig(conf),
	)
	if err != nil {
		panic(err)
	}
}

func TestApolloLoader_LoadFiles(t *testing.T) {
	err := loader.LoadFiles(&conf, "app.yml", "connect.yml")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(conf)
}

func TestApolloLoader_WatchFiles(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	err := loader.WatchFiles(ctx, watchFunc, "app.yml")
	if err != nil {
		t.Fatal(err)
	}
}

func watchFunc(name string, data []byte) error {
	fmt.Println(name)
	fmt.Println(string(data))
	err := yaml.NewCodec().Unmarshal(data, &conf)
	if err != nil {
		return err
	}
	fmt.Println(conf)
	return nil
}
