package middleware

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
)

var echoContextKey contextKey = "echo_context_key"

func EchoContextToContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := StoreEchoContextInContext(req.Context(), c)
			c.SetRequest(req.WithContext(ctx))
			return next(c)
		}
	}
}

func StoreEchoContextInContext(ctx context.Context, c echo.Context) context.Context {
	return context.WithValue(ctx, echoContextKey, c)
}

// EchoContextFromContext returns gin context from http context (if exists)
func EchoContextFromContext(ctx context.Context) (echo.Context, error) {
	echoContext := ctx.Value(echoContextKey)
	if echoContext == nil {
		err := fmt.Errorf("Could not retrieve echo.Context")
		return nil, err
	}

	gc, ok := echoContext.(echo.Context)
	if !ok {
		err := fmt.Errorf("echo.Context has wrong type")
		return nil, err
	}
	return gc, nil
}
