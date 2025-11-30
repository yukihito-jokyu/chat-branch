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

func (m *mockProjectRepository) Create(ctx context.Context, project *model.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

type mockChatRepository struct {
	mock.Mock
}

func (m *mockChatRepository) Create(ctx context.Context, chat *model.Chat) error {
	args := m.Called(ctx, chat)
	return args.Error(0)
}

type mockMessageRepository struct {
	mock.Mock
}

func (m *mockMessageRepository) Create(ctx context.Context, message *model.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

type mockTransactionManager struct {
	mock.Mock
}

func (m *mockTransactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	// 単純にfnを実行する
	return fn(ctx)
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
					{UUID: "p1", UserUUID: "user-1", Title: "Project 1", UpdatedAt: time.Now()},
					{UUID: "p2", UserUUID: "user-1", Title: "Project 2", UpdatedAt: time.Now()},
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
			mockChatRepo := new(mockChatRepository)
			mockMessageRepo := new(mockMessageRepository)
			mockTxManager := new(mockTransactionManager)
			tt.setupMock(mockRepo)

			u := NewProjectUsecase(mockRepo, mockChatRepo, mockMessageRepo, mockTxManager)
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

func TestProjectUsecase_CreateProject(t *testing.T) {
	type args struct {
		userUUID       string
		initialMessage string
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(mRepo *mockProjectRepository, mChat *mockChatRepository, mMsg *mockMessageRepository, mTx *mockTransactionManager)
		wantErr   bool
	}{
		{
			name: "正常系: プロジェクト作成が成功すること",
			args: args{
				userUUID:       "user-1",
				initialMessage: "Hello",
			},
			setupMock: func(mRepo *mockProjectRepository, mChat *mockChatRepository, mMsg *mockMessageRepository, mTx *mockTransactionManager) {
				mRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *model.Project) bool {
					return p.UserUUID == "user-1" && p.Title == "Hello"
				})).Return(nil)
				mChat.On("Create", mock.Anything, mock.MatchedBy(func(c *model.Chat) bool {
					return c.UserUUID == "user-1" && c.Title == "Hello"
				})).Return(nil)
				mMsg.On("Create", mock.Anything, mock.MatchedBy(func(m *model.Message) bool {
					return m.UserUUID == "user-1" && m.Content == "Hello" && m.Role == "user"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "異常系: プロジェクト作成失敗時にエラーになること",
			args: args{
				userUUID:       "user-1",
				initialMessage: "Hello",
			},
			setupMock: func(mRepo *mockProjectRepository, mChat *mockChatRepository, mMsg *mockMessageRepository, mTx *mockTransactionManager) {
				mRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "異常系: チャット作成失敗時にエラーになること",
			args: args{
				userUUID:       "user-1",
				initialMessage: "Hello",
			},
			setupMock: func(mRepo *mockProjectRepository, mChat *mockChatRepository, mMsg *mockMessageRepository, mTx *mockTransactionManager) {
				mRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				mChat.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "異常系: メッセージ作成失敗時にエラーになること",
			args: args{
				userUUID:       "user-1",
				initialMessage: "Hello",
			},
			setupMock: func(mRepo *mockProjectRepository, mChat *mockChatRepository, mMsg *mockMessageRepository, mTx *mockTransactionManager) {
				mRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				mChat.On("Create", mock.Anything, mock.Anything).Return(nil)
				mMsg.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockProjectRepository)
			mockChatRepo := new(mockChatRepository)
			mockMessageRepo := new(mockMessageRepository)
			mockTxManager := new(mockTransactionManager)
			tt.setupMock(mockRepo, mockChatRepo, mockMessageRepo, mockTxManager)

			u := NewProjectUsecase(mockRepo, mockChatRepo, mockMessageRepo, mockTxManager)
			p, c, m, err := u.CreateProject(context.Background(), tt.args.userUUID, tt.args.initialMessage)

			if (err != nil) != tt.wantErr {
				t.Errorf("projectUsecase.CreateProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, p)
				assert.NotNil(t, c)
				assert.NotNil(t, m)
				assert.Equal(t, tt.args.userUUID, p.UserUUID)
				assert.Equal(t, tt.args.initialMessage, p.Title)
			}
		})
	}
}
