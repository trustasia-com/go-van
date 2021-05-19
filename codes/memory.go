// Package codes provides ...
package codes

import (
	"fmt"
)

// MemoryTranslator memory translator
// it's example, you can implements i18n.go/translator
type MemoryTranslator struct {
	Code2Desc map[string]map[Code]string
}

// Tr translate lang, should not manual call
func (trans MemoryTranslator) Tr(lang string, code Code,
	args ...interface{}) string {

	str := code.String()

	codes := trans.Code2Desc[lang]
	desc, ok := codes[code]
	if !ok {
		return "Code" + str
	}
	return fmt.Sprintf(str+desc, args...)
}

// SupportedLang supported language
func (trans MemoryTranslator) SupportedLang() []string {
	return []string{LangZhCN, LangEnUS}
}
