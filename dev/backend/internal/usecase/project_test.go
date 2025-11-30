package usecase

import (
	"backend/internal/domain/model"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// モックの定義
type mockProjectRepository struct {
	mock.Mock
}

func (m *mockProjectRepository) FindAllByUserUUID(ctx context.Context, userUUID string) ([]*model.Project, error) {
	args := m.Called(ctx, userUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Project), args.Error(1)
}

func TestProjectUsecase_GetProjects(t *testing.T) {
	type args struct {
		userUUID string
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(m *mockProjectRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name: "正常系: プロジェクト一覧が取得できること",
			args: args{
				userUUID: "user-1",
			},
			setupMock: func(m *mockProjectRepository) {
				projects := []*model.Project{
					{ID: "p1", UserID: "user-1", Title: "Project 1", UpdatedAt: time.Now()},
					{ID: "p2", UserID: "user-1", Title: "Project 2", UpdatedAt: time.Now()},
				}
				m.On("FindAllByUserUUID", mock.Anything, "user-1").Return(projects, nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "異常系: リポジトリでエラーが発生した場合エラーになること",
			args: args{
				userUUID: "user-error",
			},
			setupMock: func(m *mockProjectRepository) {
				m.On("FindAllByUserUUID", mock.Anything, "user-error").Return(nil, errors.New("db error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockProjectRepository)
			tt.setupMock(mockRepo)

			u := NewProjectUsecase(mockRepo)
			got, err := u.GetProjects(context.Background(), tt.args.userUUID)

			if (err != nil) != tt.wantErr {
				t.Errorf("projectUsecase.GetProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Len(t, got, tt.wantCount)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
