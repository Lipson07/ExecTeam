package handler

import (
	"net/http"

	"backend/internal/model"

	"github.com/gin-gonic/gin"
)

func writeJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func writeError(c *gin.Context, status int, message string) {
	c.JSON(status, model.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}
