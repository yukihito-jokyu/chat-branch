package router

import (
	"backend/internal/handler"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	e.GET("/health", handler.HealthCheck)

	// API group example
	// api := e.Group("/api")
	// api.GET("/users", handler.GetUsers)
}
