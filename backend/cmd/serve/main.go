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

	"golang.org/x/text/language"

	"github.com/go-redis/redis"

	"github.com/robfig/cron"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
	_authUsecase "github.com/kichiyaki/graphql-starter/backend/auth/usecase"
	_email "github.com/kichiyaki/graphql-starter/backend/email"
	_graphqlHandler "github.com/kichiyaki/graphql-starter/backend/graphql/delivery/http"
	_middleware "github.com/kichiyaki/graphql-starter/backend/middleware"
	"github.com/kichiyaki/graphql-starter/backend/postgre"
	redisStore "github.com/kichiyaki/graphql-starter/backend/sessions/redis"
	_tokenCron "github.com/kichiyaki/graphql-starter/backend/token/cron"
	_tokenRepo "github.com/kichiyaki/graphql-starter/backend/token/repository"
	_userRepo "github.com/kichiyaki/graphql-starter/backend/user/repository"
	_userUsecase "github.com/kichiyaki/graphql-starter/backend/user/usecase"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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
		SetApplicationName(viper.GetString("application.name"))
	conn, err := postgre.NewDatabase(dbCfg)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	redisConn := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("session.store.address"),
		Password: viper.GetString("session.store.password"), // no password set
		DB:       viper.GetInt("session.store.db"),          // use default DB
	})
	defer redisConn.Close()
	sessionStore := redisStore.NewRedisStore(redisConn, []byte(viper.GetString("session.secretKey")))

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

	authUsecase := _authUsecase.NewAuthUsecase(&_authUsecase.Config{
		SessStore:       sessionStore,
		UserRepo:        userRepo,
		TokenRepo:       tokenRepo,
		Email:           email,
		FrontendURL:     viper.GetString("application.frontend"),
		ApplicationName: viper.GetString("application.name"),
	})
	userUsecase := _userUsecase.NewUserUsecase(userRepo, authUsecase)

	c := cron.New()
	defer c.Stop()
	_tokenCron.InitTokenCron(c, tokenRepo)

	go func() {
		c.Start()
	}()

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{viper.GetString("application.frontend")}),
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
	bundle := i18n.NewBundle(language.Polish)
	bundle.MustLoadMessageFile("../../i18n/locales/active.pl.json")
	middleware := _middleware.NewMiddleware(userRepo, bundle)
	router := gin.Default()
	// store := cookie.NewStore([]byte(viper.GetString("session.secretKey")))
	router.Use(sessions.Sessions(viper.GetString("session.name"), sessionStore))
	router.Use(middleware.GinContextToContextMiddleware())
	router.Use(middleware.AuthMiddleware())
	router.Use(middleware.LocalizeMiddleware())
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

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	<-channel

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
