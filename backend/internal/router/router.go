package router

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"backend/internal/handler"
	"backend/internal/middleware"
)

func NewRouter(
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	oauthHandler *handler.OAuthHandler,
	authMW *middleware.AuthMiddleware,
	corsMW *cors.Cors,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(corsMW.Handler)

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/verify-email", authHandler.VerifyEmail)
		r.Post("/resend-code", authHandler.ResendCode)

		r.Get("/auth/{provider}", oauthHandler.GetAuthURL)
		r.Get("/auth/{provider}/callback", oauthHandler.Callback)

		r.Group(func(r chi.Router) {
			r.Use(authMW.RequireAuth)

			r.Get("/users/me", userHandler.GetMe)
			r.Put("/users/me", userHandler.UpdateMe)
			r.Get("/users", userHandler.GetAll)
			r.Get("/users/search", userHandler.SearchByName)
			r.Get("/users/{id}", userHandler.GetByID)
			r.Delete("/users/{id}", userHandler.Delete)
		})
	})

	return r
}
