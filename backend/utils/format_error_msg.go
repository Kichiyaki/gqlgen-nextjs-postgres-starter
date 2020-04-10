package utils

import (
	"backend/errors"
	"backend/middleware"
	"backend/user/validation"
	"context"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func FormatErrorMsg(ctx context.Context, err2 error) error {
	localizer, err := middleware.LocalizerFromContext(ctx)
	if err != nil {
		return err
	}
	graphqlErr := errors.ToGqlError(err2)
	defaultMsg := &i18n.Message{
		ID:    graphqlErr.Message,
		One:   graphqlErr.Message,
		Other: graphqlErr.Message,
	}
	switch graphqlErr.Message {
	case errors.ErrLoginPolicy, errors.ErrDisplayNamePolicy:
		graphqlErr.Message = localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: graphqlErr.Message,
			TemplateData: map[string]interface{}{
				"MinLength": validation.MinimumLoginLength,
				"MaxLength": validation.MaximumLoginLength,
			},
			DefaultMessage: defaultMsg,
		})
	case errors.ErrPasswordPolicy:
		graphqlErr.Message = localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: graphqlErr.Message,
			TemplateData: map[string]interface{}{
				"MinLength": validation.MinimumPasswordLength,
				"MaxLength": validation.MaximumPasswordLength,
			},
			DefaultMessage: defaultMsg,
		})
	default:
		graphqlErr.Message = localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID:      graphqlErr.Message,
			DefaultMessage: defaultMsg,
		})
	}
	return graphqlErr
}
