package utils

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func GetErrorMsg(localizer *i18n.Localizer, msgID string) error {
	return GetErrorMsgWithData(localizer, msgID, map[string]interface{}{})
}

func GetErrorMsgWithData(localizer *i18n.Localizer, msgID string, templateData map[string]interface{}) error {
	return fmt.Errorf(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: msgID,
		TemplateData: templateData}))
}

func GetErrorMsgWithDataAndPluralCount(localizer *i18n.Localizer, msgID string, templateData map[string]interface{}, pluralCount int) error {
	return fmt.Errorf(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: msgID,
		TemplateData: templateData,
		PluralCount:  pluralCount}))
}
