// Package codes provides ...
package codes

// WithTranslator specific translator
func WithTranslator(trans Translator) {
	globalI18n.translator = trans

	langs := trans.SupportedLang()
	if len(langs) == 0 {
		panic("codes: warning: not found supported lang")
	}
	globalI18n.supportedLang = langs
}

// i18nInstance international native
type i18nInstance struct {
	// cache from translator
	supportedLang []string
	// translator
	translator Translator
}
