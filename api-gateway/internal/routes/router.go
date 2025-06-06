package routes

import (
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	//authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	//todoServiceURL := os.Getenv("TODO_SERVICE_URL")
	// Auth service routes (open)
	e.Any("/auth/*", handlers.ProxyHandler("AUTH_SERVICE_URL"))

	// Todo service routes (protected)
	group := e.Group("/todo", middleware.JWTMiddleware)
	group.Any("/*", handlers.ProxyHandler("TODO_SERVICE_URL"))
}
