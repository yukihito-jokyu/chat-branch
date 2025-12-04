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

func (m *mockChatRepository) FindByID(ctx context.Context, uuid string) (*model.Chat, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chat), args.Error(1)
}

func (m *mockChatRepository) UpdateStatus(ctx context.Context, chatUUID string, status string) error {
	args := m.Called(ctx, chatUUID, status)
	return args.Error(0)
}

func (m *mockChatRepository) FindOldestByProjectUUID(ctx context.Context, projectUUID string) (*model.Chat, error) {
	args := m.Called(ctx, projectUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chat), args.Error(1)
}

func (m *mockChatRepository) CountByProjectUUID(ctx context.Context, projectUUID string) (int64, error) {
	args := m.Called(ctx, projectUUID)
	return args.Get(0).(int64), args.Error(1)
}

type mockMessageRepository struct {
	mock.Mock
}

func (m *mockMessageRepository) Create(ctx context.Context, message *model.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *mockMessageRepository) FindMessagesByChatID(ctx context.Context, chatUUID string) ([]*model.Message, error) {
	args := m.Called(ctx, chatUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Message), args.Error(1)
}

func (m *mockMessageRepository) UpdateContextSummary(ctx context.Context, messageUUID string, summary string) error {
	args := m.Called(ctx, messageUUID, summary)
	return args.Error(0)
}

func (m *mockMessageRepository) FindLatestMessageWithSummary(ctx context.Context, chatUUID string) (*model.Message, error) {
	args := m.Called(ctx, chatUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Message), args.Error(1)
}

func (m *mockMessageRepository) FindLatestMessageByRole(ctx context.Context, chatUUID string, role string) (*model.Message, error) {
	args := m.Called(ctx, chatUUID, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Message), args.Error(1)
}

func (m *mockMessageRepository) FindByID(ctx context.Context, uuid string) (*model.Message, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Message), args.Error(1)
}

type mockEdgeRepository struct {
	mock.Mock
}

func (m *mockEdgeRepository) FindEdgesByChatID(ctx context.Context, chatUUID string) ([]*model.Edge, error) {
	args := m.Called(ctx, chatUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Edge), args.Error(1)
}

func (m *mockEdgeRepository) Create(ctx context.Context, edge *model.Edge) error {
	args := m.Called(ctx, edge)
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
			mockEdgeRepo := new(mockEdgeRepository)
			mockTxManager := new(mockTransactionManager)
			tt.setupMock(mockRepo)

			u := NewProjectUsecase(mockRepo, mockChatRepo, mockMessageRepo, mockEdgeRepo, mockTxManager)
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
				initialMessage: "Hello",
			},
			setupMock: func(mRepo *mockProjectRepository, mChat *mockChatRepository, mMsg *mockMessageRepository, mTx *mockTransactionManager) {
				mRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *model.Project) bool {
					return p.Title == "Hello"
				})).Return(nil)
				mChat.On("Create", mock.Anything, mock.MatchedBy(func(c *model.Chat) bool {
					return c.Title == "Hello"
				})).Return(nil)
				mMsg.On("Create", mock.Anything, mock.MatchedBy(func(m *model.Message) bool {
					return m.Content == "Hello" && m.Role == "user"
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
			mockEdgeRepo := new(mockEdgeRepository)
			mockTxManager := new(mockTransactionManager)
			tt.setupMock(mockRepo, mockChatRepo, mockMessageRepo, mockTxManager)

			u := NewProjectUsecase(mockRepo, mockChatRepo, mockMessageRepo, mockEdgeRepo, mockTxManager)
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

func TestProjectUsecase_GetParentChat(t *testing.T) {
	type args struct {
		projectUUID string
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(mRepo *mockProjectRepository, mChat *mockChatRepository, mMsg *mockMessageRepository, mTx *mockTransactionManager)
		want      *model.Chat
		wantErr   bool
	}{
		{
			name: "正常系: 親チャットが取得できること",
			args: args{
				projectUUID: "project-uuid",
			},
			setupMock: func(mRepo *mockProjectRepository, mChat *mockChatRepository, mMsg *mockMessageRepository, mTx *mockTransactionManager) {
				mChat.On("FindOldestByProjectUUID", mock.Anything, "project-uuid").Return(&model.Chat{
					UUID:        "chat-uuid",
					ProjectUUID: "project-uuid",
				}, nil)
			},
			want: &model.Chat{
				UUID:        "chat-uuid",
				ProjectUUID: "project-uuid",
			},
			wantErr: false,
		},
		{
			name: "異常系: リポジトリでエラーが発生した場合エラーになること",
			args: args{
				projectUUID: "project-error",
			},
			setupMock: func(mRepo *mockProjectRepository, mChat *mockChatRepository, mMsg *mockMessageRepository, mTx *mockTransactionManager) {
				mChat.On("FindOldestByProjectUUID", mock.Anything, "project-error").Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockProjectRepository)
			mockChatRepo := new(mockChatRepository)
			mockMessageRepo := new(mockMessageRepository)
			mockEdgeRepo := new(mockEdgeRepository)
			mockTxManager := new(mockTransactionManager)
			tt.setupMock(mockRepo, mockChatRepo, mockMessageRepo, mockTxManager)

			u := NewProjectUsecase(mockRepo, mockChatRepo, mockMessageRepo, mockEdgeRepo, mockTxManager)
			got, err := u.GetParentChat(context.Background(), tt.args.projectUUID)

			if (err != nil) != tt.wantErr {
				t.Errorf("projectUsecase.GetParentChat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestProjectUsecase_GetProjectTree(t *testing.T) {
	type mocks struct {
		projectRepo *mockProjectRepository
		chatRepo    *mockChatRepository
		messageRepo *mockMessageRepository
		edgeRepo    *mockEdgeRepository
		txManager   *mockTransactionManager
	}
	type args struct {
		projectUUID string
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks)
		want      *model.ProjectTree
		wantErr   bool
	}{
		{
			name: "正常系: プロジェクトツリーが取得できること",
			args: args{
				projectUUID: "project-uuid",
			},
			setupMock: func(m *mocks) {
				// 1. Root Chat
				m.chatRepo.On("FindOldestByProjectUUID", mock.Anything, "project-uuid").Return(&model.Chat{UUID: "root-chat"}, nil)

				// 2. Root Chat Messages
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "root-chat").Return([]*model.Message{
					{
						UUID:    "msg-1",
						Role:    "user",
						Content: "user prompt",
						Forks:   []model.Fork{},
					},
					{
						UUID:      "msg-2",
						Role:      "assistant",
						Content:   "ai response",
						PositionX: 100,
						PositionY: 200,
						Forks: []model.Fork{
							{ChatUUID: "child-chat"},
						},
					},
				}, nil)

				// 3. Root Chat Edges
				m.edgeRepo.On("FindEdgesByChatID", mock.Anything, "root-chat").Return([]*model.Edge{
					{UUID: "edge-1", SourceMessageUUID: "msg-2", TargetMessageUUID: "msg-3"},
				}, nil)

				// 4. Child Chat Messages
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "child-chat").Return([]*model.Message{
					{
						UUID:      "msg-3",
						Role:      "assistant",
						Content:   "child response",
						PositionX: 300,
						PositionY: 400,
						Forks:     []model.Fork{},
					},
				}, nil)

				// 5. Child Chat Edges
				m.edgeRepo.On("FindEdgesByChatID", mock.Anything, "child-chat").Return([]*model.Edge{}, nil)
			},
			want: &model.ProjectTree{
				Nodes: []model.ProjectNode{
					{
						ID: "msg-2",
						Data: model.ProjectNodeData{
							UserMessage: func() *string { s := "user prompt"; return &s }(),
							Assistant:   "ai response",
						},
						Position: model.ProjectNodePosition{X: 100, Y: 200},
					},
					{
						ID: "msg-3",
						Data: model.ProjectNodeData{
							UserMessage: nil,
							Assistant:   "child response",
						},
						Position: model.ProjectNodePosition{X: 300, Y: 400},
					},
				},
				Edges: []model.ProjectEdge{
					{ID: "edge-1", Source: "msg-2", Target: "msg-3"},
				},
			},
			wantErr: false,
		},
		{
			name: "異常系: ルートチャットが見つからない場合エラー",
			args: args{
				projectUUID: "project-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindOldestByProjectUUID", mock.Anything, "project-uuid").Return(nil, errors.New("not found"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系: メッセージ取得失敗",
			args: args{
				projectUUID: "project-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindOldestByProjectUUID", mock.Anything, "project-uuid").Return(&model.Chat{UUID: "root-chat"}, nil)
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "root-chat").Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系: エッジ取得失敗",
			args: args{
				projectUUID: "project-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindOldestByProjectUUID", mock.Anything, "project-uuid").Return(&model.Chat{UUID: "root-chat"}, nil)
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "root-chat").Return([]*model.Message{}, nil)
				m.edgeRepo.On("FindEdgesByChatID", mock.Anything, "root-chat").Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mocks{
				projectRepo: &mockProjectRepository{},
				chatRepo:    &mockChatRepository{},
				messageRepo: &mockMessageRepository{},
				edgeRepo:    &mockEdgeRepository{},
				txManager:   &mockTransactionManager{},
			}
			tt.setupMock(m)

			u := NewProjectUsecase(m.projectRepo, m.chatRepo, m.messageRepo, m.edgeRepo, m.txManager)

			got, err := u.GetProjectTree(context.Background(), tt.args.projectUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("projectUsecase.GetProjectTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, len(tt.want.Nodes), len(got.Nodes))
				assert.Equal(t, len(tt.want.Edges), len(got.Edges))
			}
		})
	}
}
