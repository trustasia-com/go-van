// Package main provides ...
package main

import (
	"fmt"

	"github.com/trustasia-com/go-van/pkg/codes"

	"github.com/unknwon/i18n"
)

type goI18nTranslator struct{}

// Tr translate lang, should not manual call
func (t goI18nTranslator) Tr(lang string, code codes.Code,
	args ...any) string {
	return i18n.Tr(lang, fmt.Sprint(int(code)), args...)
}

// SupportedLang supported language
func (t goI18nTranslator) SupportedLang() []string {
	return i18n.ListLangs()
}

var (
	invalidUsername codes.Code = 1002
	// and more code
	// ...
)

func init() {
	i18n.SetMessage(codes.LangZhCN, "zh-cn.ini")
	i18n.SetMessage(codes.LangEnUS, "en-us.ini")
	trans := goI18nTranslator{}

	codes.WithTranslator(trans)
}

func main() {
	fmt.Println(invalidUsername.Tr(codes.LangZhCN))
	fmt.Println(invalidUsername.Tr(codes.LangEnUS))
}
