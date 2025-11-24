package handler

import (
	"net/http"

	"backend/internal/model"
	"github.com/labstack/echo/v4"
)

func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, model.Response{
		Status: "ok",
	})
}
