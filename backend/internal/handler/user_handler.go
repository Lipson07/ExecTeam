package handler

import (
	"net/http"
	"strconv"

	"backend/internal/model"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userIDIface, exists := c.Get("user_id")
	if !exists {
		writeError(c, http.StatusInternalServerError, "не удалось получить пользователя")
		return
	}

	userID, ok := userIDIface.(int)
	if !ok {
		writeError(c, http.StatusInternalServerError, "не удалось получить пользователя")
		return
	}

	user, err := h.service.GetByID(userID)
	if err != nil {
		writeError(c, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, user)
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userIDIface, exists := c.Get("user_id")
	if !exists {
		writeError(c, http.StatusInternalServerError, "не удалось получить пользователя")
		return
	}

	userID, ok := userIDIface.(int)
	if !ok {
		writeError(c, http.StatusInternalServerError, "не удалось получить пользователя")
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "неверный формат данных")
		return
	}

	user, err := h.service.UpdateMe(userID, &req)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, user)
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.service.GetAll()
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, users)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(c, http.StatusBadRequest, "неверный ID")
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		writeError(c, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, user)
}

func (h *UserHandler) SearchByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		writeError(c, http.StatusBadRequest, "параметр name обязателен")
		return
	}

	users, err := h.service.SearchByName(name)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, users)
}

func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(c, http.StatusBadRequest, "неверный ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, map[string]string{"message": "пользователь удалён"})
}
