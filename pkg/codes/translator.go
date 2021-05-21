// Package codes provides ...
package codes

import "fmt"

// Translator translate code to desc
type Translator interface {
	Tr(lang string, code Code, args ...interface{}) string
	SupportedLang() []string
}

// DefaultTranslator memory translator, implements i18n.go/translator
type DefaultTranslator struct {
	Code2Desc map[string]map[Code]string
}

// Tr translate lang, should not manual call
func (t DefaultTranslator) Tr(lang string, code Code,
	args ...interface{}) string {

	str := code.String()

	codes := t.Code2Desc[lang]
	desc, ok := codes[code]
	if !ok {
		return "Code" + str
	}
	return fmt.Sprintf(str+desc, args...)
}

// SupportedLang supported language
func (t DefaultTranslator) SupportedLang() []string {
	return []string{LangZhCN, LangEnUS}
}
