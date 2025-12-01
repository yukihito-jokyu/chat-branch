package handler

import (
	"backend/config"
	"backend/internal/domain/usecase"
	"backend/internal/handler/model"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	cfg         *config.Config
}

// AuthHandler の新しいインスタンスを作成する処理
func NewAuthHandler(authUsecase usecase.AuthUsecase, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		cfg:         cfg,
	}
}

// ゲストサインアップリクエストの処理
func (h *AuthHandler) Signup(c echo.Context) error {
	ctx := c.Request().Context()
	user, token, err := h.authUsecase.GuestSignup(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "ゲストサインアップに失敗", "error", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	slog.InfoContext(ctx, "ゲストサインアップに成功", "user_uuid", user.UUID)
	return c.JSON(http.StatusOK, model.SignupResponse{
		Token: token,
		User: struct {
			UUID string `json:"uuid"`
			Name string `json:"name"`
		}{
			UUID: user.UUID,
			Name: user.Name,
		},
	})
}

// ゲストログインリクエストの処理
func (h *AuthHandler) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		slog.WarnContext(c.Request().Context(), "ログインリクエストのバインドに失敗", "error", err)
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  "error",
			Message: "invalid request",
		})
	}

	if req.UserUUID == "" {
		slog.WarnContext(c.Request().Context(), "ユーザーUUIDが指定されていません")
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  "error",
			Message: "user_uuid is required",
		})
	}

	ctx := c.Request().Context()
	token, err := h.authUsecase.GuestLogin(ctx, req.UserUUID)
	if err != nil {
		slog.ErrorContext(ctx, "ゲストログインに失敗", "error", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	slog.InfoContext(ctx, "ゲストログインに成功", "user_uuid", req.UserUUID)
	return c.JSON(http.StatusOK, model.LoginResponse{
		Token: token,
	})
}
