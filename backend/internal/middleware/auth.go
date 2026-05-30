package middleware

import (
	"net/http"
	"strings"

	"backend/internal/model"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService service.AuthService
}

func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) RequireAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{
			Error:   http.StatusText(http.StatusUnauthorized),
			Message: "токен не предоставлен",
		})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{
			Error:   http.StatusText(http.StatusUnauthorized),
			Message: "неверный формат токена",
		})
		return
	}

	userID, err := m.authService.ValidateToken(parts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{
			Error:   http.StatusText(http.StatusUnauthorized),
			Message: "невалидный токен",
		})
		return
	}

	c.Set("user_id", userID)
	c.Next()
}
