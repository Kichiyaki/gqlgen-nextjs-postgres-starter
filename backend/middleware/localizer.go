package middleware

import (
	_i18n "backend/i18n"
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var localizerContextKey contextKey = "localizer_context_key"

func LocalizerToContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			accept := req.Header.Get("Accept-Language")
			localizer := i18n.NewLocalizer(_i18n.Bundle, accept)
			ctx := StoreLocalizerInContext(req.Context(), localizer)
			c.SetRequest(req.WithContext(ctx))
			return next(c)
		}
	}
}

func StoreLocalizerInContext(ctx context.Context, localizer *i18n.Localizer) context.Context {
	return context.WithValue(ctx, localizerContextKey, localizer)
}

func LocalizerFromContext(ctx context.Context) (*i18n.Localizer, error) {
	localizer := ctx.Value(localizerContextKey)
	if localizer == nil {
		err := fmt.Errorf("Could not retrieve *i18n.Localizer")
		return nil, err
	}

	gc, ok := localizer.(*i18n.Localizer)
	if !ok {
		err := fmt.Errorf("*i18n.Localizer has wrong type")
		return nil, err
	}
	return gc, nil
}
