package handler

import (
	"backend/internal/domain/model"
	handlerModel "backend/internal/handler/model"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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
					{ID: "p1", UserID: "user-1", Title: "Project 1", UpdatedAt: time.Now()},
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
