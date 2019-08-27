package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (midd *middleware) GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := StoreGinContextInContext(c.Request.Context(), c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func StoreGinContextInContext(ctx context.Context, c *gin.Context) context.Context {
	return context.WithValue(ctx, ginContextKey, c)
}

// GinContextFromContext returns gin context from http context (if exists)
func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(ginContextKey)
	if ginContext == nil {
		err := fmt.Errorf("Could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}
