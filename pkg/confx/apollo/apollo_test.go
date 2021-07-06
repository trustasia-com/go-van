// Package apollo provides ...
package apollo

import (
	"fmt"
	"github.com/trustasia-com/go-van/pkg/codec/yaml"
	"testing"
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
	loader *apolloLoader
)

func init() {
	load, err := NewApolloLoader("test2", "dev", "http://101.132.140.237:8080", "test.yml,test2.yml", "5a9940521184403b86150ccc5e8de75d")
	if err != nil {
		panic(err)
	}
	loader = load
}

func TestApolloLoader_LoadFiles(t *testing.T) {
	err := loader.LoadFiles(&conf, "test2.yml")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(conf)
}

func TestApolloLoader_WatchFiles(t *testing.T) {
	err := loader.WatchFiles(watchFunc, "test.yml")
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
