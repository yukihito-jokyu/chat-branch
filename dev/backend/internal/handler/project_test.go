package handler

import (
	"backend/internal/domain/model"
	handlerModel "backend/internal/handler/model"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// モックの定義
type mockProjectUsecase struct {
	mock.Mock
}

func (m *mockProjectUsecase) GetProjects(ctx context.Context, userUUID string) ([]*model.Project, error) {
	args := m.Called(ctx, userUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Project), args.Error(1)
}

func (m *mockProjectUsecase) CreateProject(ctx context.Context, userUUID, initialMessage string) (*model.Project, *model.Chat, *model.Message, error) {
	args := m.Called(ctx, userUUID, initialMessage)
	if args.Get(0) == nil {
		return nil, nil, nil, args.Error(3)
	}
	return args.Get(0).(*model.Project), args.Get(1).(*model.Chat), args.Get(2).(*model.Message), args.Error(3)
}

func TestProjectHandler_GetProjects(t *testing.T) {
	type args struct {
		userUUID interface{} // コンテキストにセットする値
	}
	tests := []struct {
		name           string
		args           args
		setupMock      func(m *mockProjectUsecase)
		wantStatus     int
		wantBodySubstr string
	}{
		{
			name: "正常系: プロジェクト一覧が取得できること",
			args: args{
				userUUID: "user-1",
			},
			setupMock: func(m *mockProjectUsecase) {
				projects := []*model.Project{
					{UUID: "p1", UserUUID: "user-1", Title: "Project 1", UpdatedAt: time.Now()},
				}
				m.On("GetProjects", mock.Anything, "user-1").Return(projects, nil)
			},
			wantStatus:     http.StatusOK,
			wantBodySubstr: "Project 1",
		},
		{
			name: "異常系: ユーザーUUIDがコンテキストにない場合401エラー",
			args: args{
				userUUID: nil,
			},
			setupMock: func(m *mockProjectUsecase) {
				// 呼び出されない
			},
			wantStatus:     http.StatusUnauthorized,
			wantBodySubstr: "ユーザーUUIDの取得に失敗しました",
		},
		{
			name: "異常系: Usecaseでエラーが発生した場合500エラー",
			args: args{
				userUUID: "user-error",
			},
			setupMock: func(m *mockProjectUsecase) {
				m.On("GetProjects", mock.Anything, "user-error").Return(nil, errors.New("usecase error"))
			},
			wantStatus:     http.StatusInternalServerError,
			wantBodySubstr: "usecase error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Echoのセットアップ
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// コンテキストの設定
			if tt.args.userUUID != nil {
				c.Set("user_uuid", tt.args.userUUID)
			}

			// モックのセットアップ
			mockUsecase := new(mockProjectUsecase)
			tt.setupMock(mockUsecase)

			h := NewProjectHandler(mockUsecase)
			err := h.GetProjects(c)

			// エラーチェック (Handlerはエラーを返さない場合もあるので、ステータスコードで判断)
			if err != nil {
				// Echoのエラーハンドリングに任せる場合、ここでエラーが返るかもしれないが
				// 今回の実装ではc.JSONで返しているのでnilが返るはず
				// ただし、c.JSON自体がエラーを返す可能性はある
			}

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.wantBodySubstr)

			// レスポンスボディの構造チェック (正常系のみ)
			if tt.wantStatus == http.StatusOK {
				var res []*handlerModel.ProjectResponse
				err := json.Unmarshal(rec.Body.Bytes(), &res)
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestProjectHandler_CreateProject(t *testing.T) {
	type args struct {
		userUUID interface{}
		body     string
	}
	tests := []struct {
		name           string
		args           args
		setupMock      func(m *mockProjectUsecase)
		wantStatus     int
		wantBodySubstr string
	}{
		{
			name: "正常系: プロジェクト作成が成功すること",
			args: args{
				userUUID: "user-1",
				body:     `{"initial_message": "Hello"}`,
			},
			setupMock: func(m *mockProjectUsecase) {
				project := &model.Project{UUID: "p1", UpdatedAt: time.Now()}
				chat := &model.Chat{UUID: "c1"}
				message := &model.Message{UUID: "m1", Content: "Hello"}
				m.On("CreateProject", mock.Anything, "user-1", "Hello").Return(project, chat, message, nil)
			},
			wantStatus:     http.StatusCreated,
			wantBodySubstr: "p1",
		},
		{
			name: "異常系: ユーザーUUIDがコンテキストにない場合401エラー",
			args: args{
				userUUID: nil,
				body:     `{"initial_message": "Hello"}`,
			},
			setupMock: func(m *mockProjectUsecase) {
				// 呼び出されない
			},
			wantStatus:     http.StatusUnauthorized,
			wantBodySubstr: "ユーザーUUIDの取得に失敗しました",
		},
		{
			name: "異常系: リクエストボディが不正な場合400エラー",
			args: args{
				userUUID: "user-1",
				body:     `invalid json`,
			},
			setupMock: func(m *mockProjectUsecase) {
				// 呼び出されない
			},
			wantStatus:     http.StatusBadRequest,
			wantBodySubstr: "リクエストボディの形式が正しくありません",
		},
		{
			name: "異常系: initial_messageが空の場合400エラー",
			args: args{
				userUUID: "user-1",
				body:     `{"initial_message": ""}`,
			},
			setupMock: func(m *mockProjectUsecase) {
				// 呼び出されない
			},
			wantStatus:     http.StatusBadRequest,
			wantBodySubstr: "initial_messageは必須です",
		},
		{
			name: "異常系: Usecaseでエラーが発生した場合500エラー",
			args: args{
				userUUID: "user-error",
				body:     `{"initial_message": "Hello"}`,
			},
			setupMock: func(m *mockProjectUsecase) {
				m.On("CreateProject", mock.Anything, "user-error", "Hello").Return(nil, nil, nil, errors.New("usecase error"))
			},
			wantStatus:     http.StatusInternalServerError,
			wantBodySubstr: "usecase error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Echoのセットアップ
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/projects", strings.NewReader(tt.args.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// コンテキストの設定
			if tt.args.userUUID != nil {
				c.Set("user_uuid", tt.args.userUUID)
			}

			// モックのセットアップ
			mockUsecase := new(mockProjectUsecase)
			tt.setupMock(mockUsecase)

			h := NewProjectHandler(mockUsecase)
			err := h.CreateProject(c)

			if err != nil {
				// エラーハンドリング
			}

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.wantBodySubstr)

			mockUsecase.AssertExpectations(t)
		})
	}
}
