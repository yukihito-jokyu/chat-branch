package main

import (
	"backend/internal/router"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	router.InitRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
