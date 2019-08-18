package usecase

import (
	"golang.org/x/text/language"

	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/kichiyaki/graphql-starter/backend/utils"
)

var localizer *i18n.Localizer

func init() {
	localizer = utils.GetLocalizer(language.Polish, "../../i18n/locales/active.pl.json")
}
