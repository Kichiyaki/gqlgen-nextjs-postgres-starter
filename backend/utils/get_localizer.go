package utils

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func GetLocalizer(language language.Tag, path string) *i18n.Localizer {
	bundle := i18n.NewBundle(language)
	bundle.MustLoadMessageFile(path)

	return i18n.NewLocalizer(bundle, language.String())
}
