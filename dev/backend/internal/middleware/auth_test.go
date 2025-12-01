package middleware

import (
	"backend/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_Authenticate(t *testing.T) {
	// テスト用の設定
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			Expiration: time.Hour,
		},
	}

	// トークン生成ヘルパー
	createToken := func(secret string, claims jwt.MapClaims) string {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))
		return tokenString
	}

	tests := []struct {
		name       string
		setupReq   func(req *http.Request)
		wantStatus int
	}{
		{
			name: "正常系: 有効なトークンの場合リクエストが通ること",
			setupReq: func(req *http.Request) {
				token := createToken(cfg.JWT.Secret, jwt.MapClaims{
					"user_uuid": "test-user-uuid",
					"exp":       time.Now().Add(time.Hour).Unix(),
				})
				req.AddCookie(&http.Cookie{Name: "jwt_token", Value: token})
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "異常系: トークンがない場合401エラー",
			setupReq: func(req *http.Request) {
				// クッキーなし
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "異常系: 無効な署名のトークンの場合401エラー",
			setupReq: func(req *http.Request) {
				token := createToken("wrong-secret", jwt.MapClaims{
					"user_uuid": "test-user-uuid",
					"exp":       time.Now().Add(time.Hour).Unix(),
				})
				req.AddCookie(&http.Cookie{Name: "jwt_token", Value: token})
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "異常系: 期限切れのトークンの場合401エラー",
			setupReq: func(req *http.Request) {
				token := createToken(cfg.JWT.Secret, jwt.MapClaims{
					"user_uuid": "test-user-uuid",
					"exp":       time.Now().Add(-time.Hour).Unix(), // 過去の時間
				})
				req.AddCookie(&http.Cookie{Name: "jwt_token", Value: token})
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "異常系: user_uuidが含まれていないトークンの場合401エラー",
			setupReq: func(req *http.Request) {
				token := createToken(cfg.JWT.Secret, jwt.MapClaims{
					"exp": time.Now().Add(time.Hour).Unix(),
				})
				req.AddCookie(&http.Cookie{Name: "jwt_token", Value: token})
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			tt.setupReq(req)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			m := NewAuthMiddleware(cfg)
			h := m.Authenticate(func(c echo.Context) error {
				return c.String(http.StatusOK, "success")
			})

			if err := h(c); err != nil {
				// エラーハンドリングはMiddleware内で行われるため、ここには到達しないはず
				// ただし、c.JSONがエラーを返す可能性はある
			}

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
