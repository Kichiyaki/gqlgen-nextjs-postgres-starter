package http

import (
	"net/http"
	"time"

	"github.com/kichiyaki/graphql-starter/backend/auth"
	"github.com/kichiyaki/graphql-starter/backend/user"

	_gqlgenHandler "github.com/99designs/gqlgen/handler"
	"github.com/NYTimes/gziphandler"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kichiyaki/graphql-starter/backend/graphql/generated"
	"github.com/kichiyaki/graphql-starter/backend/graphql/resolvers"
)

type handler struct {
	userUcase user.Usecase
	authUcase auth.Usecase
}

func NewGraphqlHandler(r *gin.RouterGroup, userUcase user.Usecase, authUcase auth.Usecase) {
	h := handler{userUcase, authUcase}
	r.POST("/graphql", h.handle())
	r.GET("/graphql", h.handle())
	if gin.Mode() == gin.DebugMode {
		r.GET("/playground", h.handlePlayground(r.BasePath()+"/graphql"))
	}
}

func (handler *handler) handle() gin.HandlerFunc {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	}
	h := gziphandler.GzipHandler(
		_gqlgenHandler.GraphQL(
			generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{handler.userUcase, handler.authUcase}}),
			_gqlgenHandler.WebsocketUpgrader(upgrader),
			_gqlgenHandler.WebsocketKeepAliveDuration(15*time.Second)),
	)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (handler *handler) handlePlayground(path string) gin.HandlerFunc {
	h := _gqlgenHandler.Playground("API playground", path)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
