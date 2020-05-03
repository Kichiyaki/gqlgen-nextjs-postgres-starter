package main

import (
	"backend/auth"
	_authUsecase "backend/auth/usecase"
	"backend/email"
	_graphqlHTTPDelivery "backend/graphql/delivery/http"
	"backend/graphql/resolvers"
	"backend/i18n"
	_middleware "backend/middleware"
	"backend/postgres"
	_userRepository "backend/user/repository"
	_userUsecase "backend/user/usecase"
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"golang.org/x/text/language"

	"github.com/gorilla/sessions"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"

	"github.com/go-pg/pg/v9"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

func init() {
	os.Setenv("TZ", "UTC")
	viper.SetConfigFile("config.json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal(err)
	}
	if viper.GetBool("application.debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}

	templatesDir, err := filepath.Abs(filepath.Join(dir, "email", "templates"))
	if err != nil {
		logrus.Fatal(err)
	}
	email.NewDialer(viper.GetString("email.host"),
		viper.GetInt("email.port"),
		viper.GetString("email.username"),
		viper.GetString("email.password"))
	if err := email.LoadTemplates(templatesDir); err != nil {
		logrus.Fatal(err)
	}

	//i18n
	lang, err := language.Parse(viper.GetString("application.defaultLanguage"))
	if err != nil {
		logrus.Fatal(err)
	}
	i18n.SetDefaultLanguage(lang)
	localesDir, err := filepath.Abs(filepath.Join(dir, "i18n", "locales"))
	if err != nil {
		logrus.Fatal(err)
	}
	if err := i18n.LoadMessageFiles(localesDir); err != nil {
		logrus.Fatal(err)
	}

	dbConnConfig := &pg.Options{
		Addr:            viper.GetString("db.addr"),
		User:            viper.GetString("db.user"),
		Password:        viper.GetString("db.password"),
		Database:        viper.GetString("db.name"),
		ApplicationName: viper.GetString("application.name"),
	}
	dbConn := pg.Connect(dbConnConfig)
	defer func() {
		err := dbConn.Close()
		if err != nil {
			logrus.Fatal(err)
		}
	}()

	userRepo, err := _userRepository.NewPostgreUserRepository(dbConn)
	if err != nil {
		logrus.Fatal(err)

	}
	if err := postgres.LoadFunctionsAndTriggers(dbConn); err != nil {
		logrus.Fatal(err)
	}

	authUcase := _authUsecase.NewAuthUsecase(_authUsecase.Config{
		UserRepo:                        userRepo,
		IntervalBetweenTokensGeneration: viper.GetInt("application.intervalBetweenTokensGeneration"),
		ResetPasswordTokenExpiresIn:     viper.GetInt("application.resetPasswordTokenExpiresIn"),
		RegistrationDisabled:            viper.GetBool("application.registrationDisabled"),
	})

	userUcase := _userUsecase.NewUserUsecase(_userUsecase.Config{
		UserRepo: userRepo,
	})

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     viper.GetStringSlice("application.cors.allowOrigins"),
		AllowHeaders:     middleware.DefaultCORSConfig.AllowHeaders,
		AllowCredentials: viper.GetBool("application.cors.allowCredentials"),
	}))

	//Gzip compression
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	//Session
	store := sessions.NewCookieStore([]byte(viper.GetString("session.secret")))
	store.Options.Secure = viper.GetBool("session.cookie.secure")
	store.Options.HttpOnly = viper.GetBool("session.cookie.httpOnly")
	if sessionName := viper.GetString("session.cookie.sessionName"); sessionName != "" {
		auth.SessionName = sessionName
	}
	if domain := viper.GetString("session.cookie.domain"); domain != "" {
		store.Options.Domain = domain
	}
	if sameSite := viper.GetString("session.cookie.sameSite"); sameSite != "" {
		store.Options.SameSite = convertToHTTPSameSite(sameSite)
	}
	if maxAge := viper.GetInt("session.cookie.maxAge"); maxAge > 0 {
		store.Options.MaxAge = maxAge
	}
	e.Use(session.Middleware(store))

	e.Use(_middleware.Logger())
	g := e.Group("")
	g.Use(middleware.Secure())
	g.Use(middleware.BodyLimit(viper.GetString("application.bodyLimit")))
	g.Use(_middleware.EchoContextToContext())
	g.Use(_middleware.LocalizerToContext())
	g.Use(_middleware.Authenticate(userRepo))
	_graphqlHTTPDelivery.NewGraphqlHandler(g, &resolvers.Resolver{
		FrontendURL: viper.GetString("application.frontend"),
		AuthUcase:   authUcase,
		UserUcase:   userUcase,
	})
	go func() {
		e.Start(viper.GetString("application.address"))
	}()
	logrus.Infof("Server is listening on port %s", viper.GetString("application.address"))

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	<-channel

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	e.Shutdown(ctx)
	logrus.Info("shutting down")
	os.Exit(0)
}

func convertToHTTPSameSite(sameSite string) http.SameSite {
	switch sameSite {
	case "lax":
		return http.SameSiteLaxMode
	case "strict":
		return http.SameSiteStrictMode
	default:
		return http.SameSiteDefaultMode
	}
}
