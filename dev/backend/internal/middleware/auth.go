package middleware

import (
	"backend/config"
	"backend/internal/handler/model"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	cfg *config.Config
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		cfg: cfg,
	}
}

func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt_token")
		if err != nil {
			slog.WarnContext(c.Request().Context(), "JWTトークンが見つかりません", "error", err)
			return c.JSON(http.StatusUnauthorized, model.Response{
				Status:  "error",
				Message: "JWTトークンが見つかりません",
			})
		}

		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			slog.WarnContext(c.Request().Context(), "無効なJWTトークンです", "error", err)
			return c.JSON(http.StatusUnauthorized, model.Response{
				Status:  "error",
				Message: "無効なJWTトークンです",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			slog.WarnContext(c.Request().Context(), "無効なトークンクレームです")
			return c.JSON(http.StatusUnauthorized, model.Response{
				Status:  "error",
				Message: "無効なトークンクレームです",
			})
		}

		userUUID, ok := claims["user_uuid"].(string)
		if !ok {
			slog.WarnContext(c.Request().Context(), "トークンにユーザーUUIDが含まれていません")
			return c.JSON(http.StatusUnauthorized, model.Response{
				Status:  "error",
				Message: "トークンにユーザーUUIDが含まれていません",
			})
		}

		c.Set("user_uuid", userUUID)

		return next(c)
	}
}
