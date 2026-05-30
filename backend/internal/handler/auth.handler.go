package handler

import (
	"backend/internal/model"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, 400, "неверный формат")
		return
	}

	resp, err := h.service.Register(&req)
	if err != nil {
		writeError(c, 409, err.Error())
		return
	}

	writeJSON(c, 201, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, 400, "неверный формат")
		return
	}

	resp, err := h.service.Login(&req)
	if err != nil {
		writeError(c, 401, err.Error())
		return
	}

	writeJSON(c, 200, resp)
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req model.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, 400, "неверный формат")
		return
	}

	resp, err := h.service.VerifyEmail(req.Email, req.Code)
	if err != nil {
		writeError(c, 400, err.Error())
		return
	}

	writeJSON(c, 200, resp)
}

func (h *AuthHandler) ResendCode(c *gin.Context) {
	var req model.ResendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, 400, "неверный формат")
		return
	}

	if err := h.service.ResendCode(req.Email); err != nil {
		writeError(c, 400, err.Error())
		return
	}

	writeJSON(c, 200, map[string]string{"message": "код отправлен"})
}
