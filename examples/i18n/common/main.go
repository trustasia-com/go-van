// Package main provides ...
package main

import (
	"fmt"

	"github.com/trustasia-com/go-van/pkg/codes"
)

var (
	invalidUsername codes.Code = 1002
	// and more code
	// ...
)

func init() {
	trans := codes.DefaultTranslator{
		Code2Desc: map[string]map[codes.Code]string{
			codes.LangZhCN: {
				invalidUsername: "用户名错误",
			},
			codes.LangEnUS: {
				invalidUsername: "Invalid username",
			},
		},
	}
	codes.WithTranslator(trans)
}

func main() {
	fmt.Println(invalidUsername.Tr(codes.LangZhCN))
	fmt.Println(invalidUsername.Tr(codes.LangEnUS))

}
