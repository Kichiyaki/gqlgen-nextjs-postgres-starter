package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
	_authUsecase "github.com/kichiyaki/graphql-starter/backend/auth/usecase"
	_email "github.com/kichiyaki/graphql-starter/backend/email"
	_graphqlHandler "github.com/kichiyaki/graphql-starter/backend/graphql/delivery/http"
	_middleware "github.com/kichiyaki/graphql-starter/backend/middleware"
	"github.com/kichiyaki/graphql-starter/backend/postgre"
	_tokenRepo "github.com/kichiyaki/graphql-starter/backend/token/repository"
	_userRepo "github.com/kichiyaki/graphql-starter/backend/user/repository"
	_userUsecase "github.com/kichiyaki/graphql-starter/backend/user/usecase"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	dbCfg := postgre.
		NewConfig().
		SetDBName(viper.GetString("database.name")).
		SetPassword(viper.GetString("database.password")).
		SetPort(viper.GetString("database.port")).
		SetURI(viper.GetString("database.uri")).
		SetUser(viper.GetString("database.username")).
		SetApplicationName(viper.GetString("database.applicationName"))
	conn, err := postgre.NewDatabase(dbCfg)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	userRepo, err := _userRepo.NewPostgreUserRepository(conn)
	if err != nil {
		panic(err)
	}
	tokenRepo, err := _tokenRepo.NewPostgreTokenRepository(conn)
	if err != nil {
		panic(err)
	}
	email, err := _email.NewEmail(_email.
		NewConfig().
		SetAddress(viper.GetString("email.address")).
		SetUsername(viper.GetString("email.username")).
		SetPort(viper.GetInt("email.port")).
		SetPassword(viper.GetString("email.password")).
		SetURI(viper.GetString("email.uri")))

	userUsecase := _userUsecase.NewUserUsecase(userRepo)
	authUsecase := _authUsecase.NewAuthUsecase(userRepo, tokenRepo, email)

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000", "http://localhost:3001"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{
			"Accept",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization"}),
		handlers.AllowCredentials(),
	)
	middleware := _middleware.NewMiddleware(userRepo)
	router := gin.Default()
	store := cookie.NewStore([]byte(viper.GetString("session.secretKey")))
	router.Use(sessions.Sessions(viper.GetString("session.name"), store))
	router.Use(middleware.GinContextToContextMiddleware())
	router.Use(middleware.AuthMiddleware())
	_graphqlHandler.NewGraphqlHandler(router.Group("/api"), userUsecase, authUsecase)

	srv := &http.Server{
		Addr:        fmt.Sprintf(":%d", viper.GetInt("application.port")),
		Handler:     cors(router),
		ReadTimeout: 2 * time.Minute,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
