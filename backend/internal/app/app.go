package app

import (
	"log"
	"net/http"

	"backend/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	repository "backend/internal/repo"
	"backend/internal/router"
	"backend/internal/service"
)

type App struct {
	Server *http.Server
}

func New() *App {
	cfg := config.Load()

	db, err := repository.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	userRepo := repository.NewUserPostgres(db)
	emailService := service.NewEmailService(cfg)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, emailService)

	oauthCfgs := map[string]service.OAuthConfig{
		"github": {
			ClientID:     cfg.GitHubClientID,
			ClientSecret: cfg.GitHubClientSecret,
			RedirectURL:  cfg.GitHubRedirectURL,
			AuthURL:      "https://github.com/login/oauth/authorize",
			TokenURL:     "https://github.com/login/oauth/access_token",
			UserInfoURL:  "https://api.github.com/user",
		},
		"google": {
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			AuthURL:      "https://accounts.google.com/o/oauth2/auth",
			TokenURL:     "https://oauth2.googleapis.com/token",
			UserInfoURL:  "https://openidconnect.googleapis.com/v1/userinfo",
		},
	}

	oauthService := service.NewOAuthService(userRepo, cfg.JWTSecret, oauthCfgs)

	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	oauthHandler := handler.NewOAuthHandler(oauthService)

	authMW := middleware.NewAuthMiddleware(authService)

	r := router.NewRouter(userHandler, authHandler, oauthHandler, authMW)

	router.PrintRoutes(r)

	addr := ":" + cfg.Port

	return &App{
		Server: &http.Server{
			Addr:    addr,
			Handler: r,
		},
	}
}

func (a *App) Run() error {
	log.Printf("Сервер запущен на http://localhost%s", a.Server.Addr)
	return a.Server.ListenAndServe()
}
