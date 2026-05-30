package handler

import (
	"net/http"

	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
	service service.OAuthService
}

func NewOAuthHandler(service service.OAuthService) *OAuthHandler {
	return &OAuthHandler{service: service}
}

func (h *OAuthHandler) GetAuthURL(c *gin.Context) {
	provider := c.Param("provider")

	url, err := h.service.GetOAuthURL(provider)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusFound, url)
}

func (h *OAuthHandler) Callback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")

	if code == "" {
		writeError(c, http.StatusBadRequest, "код не передан")
		return
	}

	resp, err := h.service.HandleCallback(provider, code)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	frontendURL := "http://localhost:4200/auth/callback?token=" + resp.Token
	c.Redirect(http.StatusFound, frontendURL)
}
