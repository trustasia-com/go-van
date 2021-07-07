// Package apollo provides ...
package apollo

import (
	"context"
	"fmt"
	"github.com/trustasia-com/go-van/pkg/codec/yaml"
	"testing"
	"time"
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
	loader *apolloLoader
)

func init() {
	load, err := NewApolloLoader(WithAppId("test2"),
		WithCluster("dev"),
		WithAddr("http://101.132.140.237:8080"),
		WithNamespaceNames([]string{"test.yml", "test2.yml"}),
		WithSecret("5a9940521184403b86150ccc5e8de75d"),
	)
	if err != nil {
		panic(err)
	}
	loader = load
}

func TestApolloLoader_LoadFiles(t *testing.T) {
	err := loader.LoadFiles(&conf, "test.yml", "test2.yml")
	if err != nil {
		panic(err)
	}
	t.Log(conf)
}

func TestApolloLoader_WatchFiles(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	err := loader.WatchFiles(ctx, watchFunc, "test.yml")
	if err != nil {
		t.Fatal(err)
	}
}

func watchFunc(name string, data []byte) error {
	fmt.Println(name)
	fmt.Println(string(data))
	err := yaml.NewCodec().Unmarshal(data, &conf)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf)
	return nil
}
