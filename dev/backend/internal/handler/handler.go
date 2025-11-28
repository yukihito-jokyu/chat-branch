package handler

import (
	"net/http"

	"backend/internal/handler/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler interface {
	HealthCheck(c echo.Context) error
}

type handler struct {
	db *gorm.DB
}

func New(db *gorm.DB) Handler {
	return &handler{db: db}
}

func (h *handler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, model.Response{
		Status: "ok",
	})
}
