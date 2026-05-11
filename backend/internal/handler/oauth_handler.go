package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/service"
)

type OAuthHandler struct {
	service service.OAuthService
}

func NewOAuthHandler(service service.OAuthService) *OAuthHandler {
	return &OAuthHandler{service: service}
}

func (h *OAuthHandler) GetAuthURL(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	url, err := h.service.GetOAuthURL(provider)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Редирект на OAuth провайдера
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *OAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	code := r.URL.Query().Get("code")

	if code == "" {
		writeError(w, http.StatusBadRequest, "код не передан")
		return
	}

	resp, err := h.service.HandleCallback(provider, code)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Редирект на фронт с JWT токеном
	frontendURL := "http://localhost:4200/auth/callback?token=" + resp.Token
	http.Redirect(w, r, frontendURL, http.StatusFound)
}
