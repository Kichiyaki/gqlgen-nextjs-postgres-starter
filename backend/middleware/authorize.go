package middleware

import (
	"backend/auth"
	"backend/models"
	"backend/user"
	"context"
	"fmt"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var userContextKey contextKey = "user_ctx_key"

func Authorize(repo user.Repository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := session.Get(auth.SessionName, c)
			login, ok1 := sess.Values["login"].(string)
			password, ok2 := sess.Values["password"].(string)
			req := c.Request()
			if ok1 && ok2 {
				user, err := repo.GetByCredentials(req.Context(), login, password)
				if err == nil {
					c.SetRequest(req.WithContext(StoreUserInContext(req.Context(), user)))
				}
			}
			return next(c)
		}
	}
}
func StoreUserInContext(ctx context.Context, u *models.User) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

func UserFromContext(ctx context.Context) (*models.User, error) {
	user := ctx.Value(userContextKey)
	if user == nil {
		err := fmt.Errorf("Could not retrieve *models.User")
		return nil, err
	}

	gc, ok := user.(*models.User)
	if !ok {
		err := fmt.Errorf("*models.User has wrong type")
		return nil, err
	}
	return gc, nil
}
