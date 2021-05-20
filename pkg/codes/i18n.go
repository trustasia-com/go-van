// Package codes provides ...
package codes

type translator interface {
	Tr(lang string, code Code, args ...interface{}) string
	SupportedLang() []string
}

// i18nInstance international native
type i18nInstance struct {
	// cache from translator
	supportedLang []string
	// translator
	translator translator
}
