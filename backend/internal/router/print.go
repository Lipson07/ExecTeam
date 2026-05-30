package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func PrintRoutes(r *gin.Engine) {
	fmt.Println("Зарегистрированные маршруты:")
	fmt.Println("----------------------------")
	for _, route := range r.Routes() {
		fmt.Printf("  %-6s %s\n", route.Method, route.Path)
	}
	fmt.Println("----------------------------")
}
