package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func PrintRoutes(r *chi.Mux) {
	fmt.Println("\nЗарегистрированные маршруты:")
	fmt.Println("----------------------------")
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*", "", -1)
		fmt.Printf("  %-6s %s\n", method, route)
		return nil
	})
	fmt.Println("----------------------------\n")
}
