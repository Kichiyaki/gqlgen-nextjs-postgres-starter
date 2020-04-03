package resolvers

import (
	"backend/email"
	"backend/middleware"
	"context"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/gomail.v2"
)

func sendEmail(ctx context.Context, title, content, to string, data map[string]interface{}) error {
	localizer, err := middleware.LocalizerFromContext(ctx)
	if err != nil {
		return err
	}
	msg := gomail.NewMessage()
	msg.SetHeader("From", email.GetAddress())
	msg.SetHeader("To", to)
	t, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    title,
		TemplateData: data,
		DefaultMessage: &i18n.Message{
			ID:    title,
			One:   title,
			Other: title,
		},
	})
	if err != nil {
		return err
	}
	c, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    content,
		TemplateData: data,
		DefaultMessage: &i18n.Message{
			ID:    content,
			One:   content,
			Other: content,
		},
	})
	if err != nil {
		return err
	}
	msg.SetHeader("Subject", t)
	body, err := email.GetTemplate("default.gohtml", map[string]interface{}{
		"Title":   t,
		"Content": c,
	})
	if err != nil {
		return err
	}
	msg.SetBody("text/html", body)
	return email.Send(msg)
}
