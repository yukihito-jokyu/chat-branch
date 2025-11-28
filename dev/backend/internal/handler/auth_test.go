package handler

import (
	"backend/config"
	"backend/internal/domain/model"
	"bytes"
	"context"

	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// モックの定義
type mockAuthUsecase struct {
	mock.Mock
}

func (m *mockAuthUsecase) GuestSignup(ctx context.Context) (*model.User, string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.String(1), args.Error(2)
	}
	return args.Get(0).(*model.User), args.String(1), args.Error(2)
}

func (m *mockAuthUsecase) GuestLogin(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func TestAuthHandler_Signup(t *testing.T) {
	cfg := &config.Config{}

	tests := []struct {
		name       string
		setupMock  func(m *mockAuthUsecase)
		wantStatus int
		wantBody   string // JSONの一部が含まれているか確認
	}{
		{
			name: "正常系: サインアップが成功すること",
			setupMock: func(m *mockAuthUsecase) {
				m.On("GuestSignup", mock.Anything).Return(&model.User{
					ID:   "test-uuid",
					Name: "Guest-test",
				}, "test-token", nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `"token":"test-token"`,
		},
		{
			name: "異常系: Usecaseでエラーが発生した場合500になること",
			setupMock: func(m *mockAuthUsecase) {
				m.On("GuestSignup", mock.Anything).Return(nil, "", errors.New("internal error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `"status":"error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockUsecase := new(mockAuthUsecase)
			tt.setupMock(mockUsecase)

			h := NewAuthHandler(mockUsecase, cfg)
			if err := h.Signup(c); err != nil {
				t.Errorf("AuthHandler.Signup() error = %v", err)
			}

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.wantBody)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	cfg := &config.Config{}

	tests := []struct {
		name       string
		reqBody    string
		setupMock  func(m *mockAuthUsecase)
		wantStatus int
		wantBody   string
	}{
		{
			name:    "正常系: ログインが成功すること",
			reqBody: `{"user_id": "test-uuid"}`,
			setupMock: func(m *mockAuthUsecase) {
				m.On("GuestLogin", mock.Anything, "test-uuid").Return("test-token", nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `"token":"test-token"`,
		},
		{
			name:       "異常系: リクエストボディが不正な場合400になること",
			reqBody:    `invalid-json`,
			setupMock:  func(m *mockAuthUsecase) {},
			wantStatus: http.StatusBadRequest,
			wantBody:   `"status":"error"`,
		},
		{
			name:       "異常系: user_idが空の場合400になること",
			reqBody:    `{"user_id": ""}`,
			setupMock:  func(m *mockAuthUsecase) {},
			wantStatus: http.StatusBadRequest,
			wantBody:   `"message":"user_id is required"`,
		},
		{
			name:    "異常系: Usecaseでエラーが発生した場合500になること",
			reqBody: `{"user_id": "test-uuid"}`,
			setupMock: func(m *mockAuthUsecase) {
				m.On("GuestLogin", mock.Anything, "test-uuid").Return("", errors.New("internal error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `"status":"error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(tt.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockUsecase := new(mockAuthUsecase)
			tt.setupMock(mockUsecase)

			h := NewAuthHandler(mockUsecase, cfg)
			if err := h.Login(c); err != nil {
				t.Errorf("AuthHandler.Login() error = %v", err)
			}

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.wantBody)
			mockUsecase.AssertExpectations(t)
		})
	}
}
