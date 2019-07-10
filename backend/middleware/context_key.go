package middleware

type contextKey string

var (
	userContextKey contextKey = "UserContextKey"
	ginContextKey  contextKey = "GinContextKey"
)
