package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (midd *middleware) LocalizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accept := c.GetHeader("Accept-Language")
		localizer := i18n.NewLocalizer(midd.bundle, accept)
		c.Request = c.Request.WithContext(StoreLocalizerInContext(c.Request.Context(), localizer))
		c.Next()
	}
}

func StoreLocalizerInContext(ctx context.Context, localizer *i18n.Localizer) context.Context {
	return context.WithValue(ctx, localizerContextKey, localizer)
}

func LocalizerFromContext(ctx context.Context) (*i18n.Localizer, error) {
	localizer := ctx.Value(localizerContextKey)
	if localizer == nil {
		err := fmt.Errorf("Could not retrieve i18n.Localizer")
		return nil, err
	}

	lc, ok := localizer.(*i18n.Localizer)
	if !ok {
		err := fmt.Errorf("i18n.Localizer has wrong type")
		return nil, err
	}
	return lc, nil
}
