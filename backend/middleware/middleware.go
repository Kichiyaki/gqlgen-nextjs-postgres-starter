package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kichiyaki/graphql-starter/backend/user"
)

type Middleware interface {
	AuthMiddleware() gin.HandlerFunc
	GinContextToContextMiddleware() gin.HandlerFunc
}

type middleware struct {
	userRepo user.Repository
}

func NewMiddleware(userRepo user.Repository) Middleware {
	return &middleware{userRepo}
}
