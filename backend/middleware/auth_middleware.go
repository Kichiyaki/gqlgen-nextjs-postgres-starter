package middleware

import (
	"context"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kichiyaki/graphql-starter/backend/models"
)

func (midd *middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("user")
		id, ok := v.(int)
		if ok {
			user, err := midd.userRepo.GetByID(context.Background(), id)
			if err == nil && user.ID > 0 {
				c.Request = c.Request.WithContext(StoreUserInContext(c.Request.Context(), user))
			}
		}
		c.Next()
	}
}

// For test purposes
func StoreUserInContext(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// UserFromContext returns user from graphql context (if exists)
func UserFromContext(ctx context.Context) (*models.User, error) {
	user := ctx.Value(userContextKey)
	if user == nil {
		return nil, fmt.Errorf("Could not retrieve models.User")
	}

	u, ok := user.(*models.User)
	if !ok {
		return nil, fmt.Errorf("models.User has wrong type")
	}
	return u, nil
}
