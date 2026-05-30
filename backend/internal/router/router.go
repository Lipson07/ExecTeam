package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/internal/handler"
	"backend/internal/middleware"
)

func NewRouter(
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	oauthHandler *handler.OAuthHandler,
	authMW *middleware.AuthMiddleware,
) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	api := r.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.POST("/verify-email", authHandler.VerifyEmail)
		api.POST("/resend-code", authHandler.ResendCode)

		api.GET("/auth/:provider", oauthHandler.GetAuthURL)
		api.GET("/auth/:provider/callback", oauthHandler.Callback)

		auth := api.Group("/")
		auth.Use(authMW.RequireAuth)

		auth.GET("/users/me", userHandler.GetMe)
		auth.PUT("/users/me", userHandler.UpdateMe)
		auth.GET("/users", userHandler.GetAll)
		auth.GET("/users/search", userHandler.SearchByName)
		auth.GET("/users/:id", userHandler.GetByID)
		auth.DELETE("/users/:id", userHandler.Delete)
	}

	return r
}
