package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kichiyaki/graphql-starter/backend/user"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Middleware interface {
	AuthMiddleware() gin.HandlerFunc
	GinContextToContextMiddleware() gin.HandlerFunc
	LocalizeMiddleware() gin.HandlerFunc
}

type middleware struct {
	userRepo user.Repository
	bundle   *i18n.Bundle
}

func NewMiddleware(userRepo user.Repository, bundle *i18n.Bundle) Middleware {
	return &middleware{userRepo, bundle}
}
