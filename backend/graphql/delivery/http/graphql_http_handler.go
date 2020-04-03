package http

import (
	"backend/graphql/generated"
	"backend/graphql/resolvers"
	"fmt"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
)

func NewGraphqlHandler(g *echo.Group, r *resolvers.Resolver) error {
	if r == nil {
		return fmt.Errorf("Graphql resolver cannot be nil")
	}
	g.POST("/graphql", graphqlHandler(r))
	g.GET("/playground", playgroundHandler())
	return nil
}

// Defining the Graphql handler
func graphqlHandler(r *resolvers.Resolver) echo.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: r}))

	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

// Defining the Playground handler
func playgroundHandler() echo.HandlerFunc {
	h := playground.Handler("Playground", "/graphql")

	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
