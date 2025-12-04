package usecase

import (
	"backend/internal/domain/model"
	"context"
	"errors"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/genai"
)

// Mocks
type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(topic string, messages ...*message.Message) error {
	args := m.Called(topic, messages)
	return args.Error(0)
}

func (m *MockPublisher) Close() error {
	args := m.Called()
	return args.Error(0)
}

type MockChatRepository struct {
	mock.Mock
}

func (m *MockChatRepository) Create(ctx context.Context, chat *model.Chat) error {
	args := m.Called(ctx, chat)
	return args.Error(0)
}

func (m *MockChatRepository) FindByID(ctx context.Context, uuid string) (*model.Chat, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chat), args.Error(1)
}

func (m *MockChatRepository) UpdateStatus(ctx context.Context, chatUUID string, status string) error {
	args := m.Called(ctx, chatUUID, status)
	return args.Error(0)
}

func (m *MockChatRepository) FindOldestByProjectUUID(ctx context.Context, projectUUID string) (*model.Chat, error) {
	args := m.Called(ctx, projectUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chat), args.Error(1)
}

func (m *MockChatRepository) CountByProjectUUID(ctx context.Context, projectUUID string) (int64, error) {
	args := m.Called(ctx, projectUUID)
	return args.Get(0).(int64), args.Error(1)
}

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Create(ctx context.Context, message *model.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) FindMessagesByChatID(ctx context.Context, chatUUID string) ([]*model.Message, error) {
	args := m.Called(ctx, chatUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Message), args.Error(1)
}

func (m *MockMessageRepository) UpdateContextSummary(ctx context.Context, messageUUID string, summary string) error {
	args := m.Called(ctx, messageUUID, summary)
	return args.Error(0)
}

func (m *MockMessageRepository) FindLatestMessageWithSummary(ctx context.Context, chatUUID string) (*model.Message, error) {
	args := m.Called(ctx, chatUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Message), args.Error(1)
}

func (m *MockMessageRepository) FindLatestMessageByRole(ctx context.Context, chatUUID string, role string) (*model.Message, error) {
	args := m.Called(ctx, chatUUID, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Message), args.Error(1)
}

func (m *MockMessageRepository) FindByID(ctx context.Context, uuid string) (*model.Message, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Message), args.Error(1)
}

type MockGenAIClient struct {
	mock.Mock
}

func (m *MockGenAIClient) GenerateContentStream(ctx context.Context, model string, parts []*genai.Content, config *genai.GenerateContentConfig) func(func(*genai.GenerateContentResponse, error) bool) {
	args := m.Called(ctx, model, parts, config)
	return args.Get(0).(func(func(*genai.GenerateContentResponse, error) bool))
}

func (m *MockGenAIClient) GenerateContent(ctx context.Context, model string, parts []*genai.Content, config *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	args := m.Called(ctx, model, parts, config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*genai.GenerateContentResponse), args.Error(1)
}

type MockMessageSelectionRepository struct {
	mock.Mock
}

func (m *MockMessageSelectionRepository) Create(ctx context.Context, selection *model.MessageSelection) error {
	args := m.Called(ctx, selection)
	return args.Error(0)
}

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func TestChatUsecase_FirstStreamChat(t *testing.T) {
	type mocks struct {
		chatRepo             *MockChatRepository
		messageRepo          *MockMessageRepository
		messageSelectionRepo *MockMessageSelectionRepository
		edgeRepo             *mockEdgeRepository
		transactionManager   *MockTransactionManager
		genaiClient          *MockGenAIClient
		publisher            *MockPublisher
	}
	type args struct {
		chatUUID string
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks)
		wantErr   bool
	}{
		{
			name: "正常系: ストリームチャットが成功すること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				// 1. FindByID
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(&model.Chat{UUID: "chat-uuid"}, nil)
				// 2. FindMessagesByChatID
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{
					{Content: "hello", Role: "user"},
				}, nil)
				// 3. GenerateContentStream
				// モックイテレータ関数を作成
				mockIter := func(yield func(*genai.GenerateContentResponse, error) bool) {
					yield(&genai.GenerateContentResponse{
						Candidates: []*genai.Candidate{
							{
								Content: genai.Text("world")[0],
							},
						},
					}, nil)
				}
				m.genaiClient.On("GenerateContentStream", mock.Anything, "gemini-2.5-flash", mock.Anything, (*genai.GenerateContentConfig)(nil)).Return(mockIter)
				// 4. Create (Assistant Message)
				m.messageRepo.On("Create", mock.Anything, mock.MatchedBy(func(msg *model.Message) bool {
					// PositionY should be equal to chat.PositionY (0)
					return msg.Role == "assistant" && msg.Content == "world" && msg.ChatUUID == "chat-uuid" && msg.PositionY == 0
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "異常系: チャットが存在しない場合エラー",
			args: args{
				chatUUID: "non-existent",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "non-existent").Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name: "異常系: メッセージ履歴が不正（0件）の場合エラー",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(&model.Chat{UUID: "chat-uuid"}, nil)
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{}, nil)
			},
			wantErr: true,
		},
		{
			name: "異常系: メッセージ履歴の取得に失敗した場合エラー",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(&model.Chat{UUID: "chat-uuid"}, nil)
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "異常系: GenAI APIからの受信エラー",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(&model.Chat{UUID: "chat-uuid"}, nil)
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{
					{Content: "hello", Role: "user"},
				}, nil)

				mockIter := func(yield func(*genai.GenerateContentResponse, error) bool) {
					yield(nil, errors.New("genai error"))
				}
				m.genaiClient.On("GenerateContentStream", mock.Anything, "gemini-2.5-flash", mock.Anything, (*genai.GenerateContentConfig)(nil)).Return(mockIter)
			},
			wantErr: true,
		},
		{
			name: "異常系: アシスタントメッセージの保存に失敗した場合エラー",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(&model.Chat{UUID: "chat-uuid"}, nil)
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{
					{Content: "hello", Role: "user"},
				}, nil)

				mockIter := func(yield func(*genai.GenerateContentResponse, error) bool) {
					yield(&genai.GenerateContentResponse{
						Candidates: []*genai.Candidate{
							{
								Content: genai.Text("world")[0],
							},
						},
					}, nil)
				}
				m.genaiClient.On("GenerateContentStream", mock.Anything, "gemini-2.5-flash", mock.Anything, (*genai.GenerateContentConfig)(nil)).Return(mockIter)

				m.messageRepo.On("Create", mock.Anything, mock.MatchedBy(func(msg *model.Message) bool {
					return msg.Role == "assistant" && msg.Content == "world"
				})).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mocks{
				chatRepo:             &MockChatRepository{},
				messageRepo:          &MockMessageRepository{},
				messageSelectionRepo: &MockMessageSelectionRepository{},
				edgeRepo:             &mockEdgeRepository{},
				transactionManager:   &MockTransactionManager{},
				genaiClient:          &MockGenAIClient{},
				publisher:            &MockPublisher{},
			}
			tt.setupMock(m)

			u := NewChatUsecase(m.chatRepo, m.messageRepo, m.messageSelectionRepo, m.edgeRepo, m.transactionManager, m.genaiClient, m.publisher)

			outputChan := make(chan string, 10)
			err := u.FirstStreamChat(context.Background(), tt.args.chatUUID, outputChan)

			if (err != nil) != tt.wantErr {
				t.Errorf("chatUsecase.FirstStreamChat() error = %v, wantErr %v", err, tt.wantErr)
			}

			// チャネルの読み出し（正常系の場合）
			if !tt.wantErr {
				close(outputChan)
				var output string
				for s := range outputChan {
					output += s
				}
				assert.Equal(t, "world", output)
			}
		})
	}
}

func TestChatUsecase_GetChat(t *testing.T) {
	type mocks struct {
		chatRepo             *MockChatRepository
		messageRepo          *MockMessageRepository
		messageSelectionRepo *MockMessageSelectionRepository
		edgeRepo             *mockEdgeRepository
		transactionManager   *MockTransactionManager
		genaiClient          *MockGenAIClient
		publisher            *MockPublisher
	}
	type args struct {
		chatUUID string
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks)
		want      *model.Chat
		wantErr   bool
	}{
		{
			name: "正常系: チャットが取得できること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(&model.Chat{
					UUID: "chat-uuid",
				}, nil)
			},
			want: &model.Chat{
				UUID: "chat-uuid",
			},
			wantErr: false,
		},
		{
			name: "異常系: Repositoryがエラーを返した場合",
			args: args{
				chatUUID: "error-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "error-uuid").Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mocks{
				chatRepo:             &MockChatRepository{},
				messageRepo:          &MockMessageRepository{},
				messageSelectionRepo: &MockMessageSelectionRepository{},
				edgeRepo:             &mockEdgeRepository{},
				transactionManager:   &MockTransactionManager{},
				genaiClient:          &MockGenAIClient{},
				publisher:            &MockPublisher{},
			}
			tt.setupMock(m)

			u := NewChatUsecase(m.chatRepo, m.messageRepo, m.messageSelectionRepo, m.edgeRepo, m.transactionManager, m.genaiClient, m.publisher)

			got, err := u.GetChat(context.Background(), tt.args.chatUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("chatUsecase.GetChat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestChatUsecase_GetMessages(t *testing.T) {
	type mocks struct {
		chatRepo             *MockChatRepository
		messageRepo          *MockMessageRepository
		messageSelectionRepo *MockMessageSelectionRepository
		edgeRepo             *mockEdgeRepository
		transactionManager   *MockTransactionManager
		genaiClient          *MockGenAIClient
		publisher            *MockPublisher
	}
	type args struct {
		chatUUID string
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks)
		want      []*model.Message
		wantErr   bool
	}{
		{
			name: "正常系: メッセージ一覧が取得できること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{
					{
						UUID:     "msg-1",
						ChatUUID: "chat-uuid",
						Role:     "user",
						Content:  "hello",
						Forks: []model.Fork{
							{
								ChatUUID:     "child-chat",
								SelectedText: "hello",
								RangeStart:   0,
								RangeEnd:     5,
							},
						},
					},
				}, nil)
			},
			want: []*model.Message{
				{
					UUID:     "msg-1",
					ChatUUID: "chat-uuid",
					Role:     "user",
					Content:  "hello",
					Forks: []model.Fork{
						{
							ChatUUID:     "child-chat",
							SelectedText: "hello",
							RangeStart:   0,
							RangeEnd:     5,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "正常系: マージレポートが親メッセージにネストされること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				parentUUID := "msg-1"
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{
					{
						UUID:     "msg-1",
						ChatUUID: "chat-uuid",
						Role:     "user",
						Content:  "parent message",
					},
					{
						UUID:              "report-1",
						ChatUUID:          "chat-uuid",
						Role:              "merge_report",
						Content:           "merge report content",
						ParentMessageUUID: &parentUUID,
					},
				}, nil)
			},
			want: []*model.Message{
				{
					UUID:     "msg-1",
					ChatUUID: "chat-uuid",
					Role:     "user",
					Content:  "parent message",
					MergeReports: []*model.Message{
						{
							UUID:              "report-1",
							ChatUUID:          "chat-uuid",
							Role:              "merge_report",
							Content:           "merge report content",
							ParentMessageUUID: func() *string { s := "msg-1"; return &s }(),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "異常系: Repositoryがエラーを返した場合",
			args: args{
				chatUUID: "error-uuid",
			},
			setupMock: func(m *mocks) {
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "error-uuid").Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mocks{
				chatRepo:             &MockChatRepository{},
				messageRepo:          &MockMessageRepository{},
				messageSelectionRepo: &MockMessageSelectionRepository{},
				edgeRepo:             &mockEdgeRepository{},
				transactionManager:   &MockTransactionManager{},
				genaiClient:          &MockGenAIClient{},
				publisher:            &MockPublisher{},
			}
			tt.setupMock(m)

			u := NewChatUsecase(m.chatRepo, m.messageRepo, m.messageSelectionRepo, m.edgeRepo, m.transactionManager, m.genaiClient, m.publisher)

			got, err := u.GetMessages(context.Background(), tt.args.chatUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("chatUsecase.GetMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestChatUsecase_SendMessage(t *testing.T) {
	type mocks struct {
		chatRepo             *MockChatRepository
		messageRepo          *MockMessageRepository
		messageSelectionRepo *MockMessageSelectionRepository
		edgeRepo             *mockEdgeRepository
		transactionManager   *MockTransactionManager
		genaiClient          *MockGenAIClient
		publisher            *MockPublisher
	}
	type args struct {
		chatUUID string
		content  string
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks)
		want      *model.Message
		wantErr   bool
	}{
		{
			name: "正常系: メッセージ送信が成功すること",
			args: args{
				chatUUID: "chat-uuid",
				content:  "hello",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(&model.Chat{UUID: "chat-uuid"}, nil)
				m.messageRepo.On("Create", mock.Anything, mock.MatchedBy(func(msg *model.Message) bool {
					return msg.Role == "user" && msg.Content == "hello" && msg.ChatUUID == "chat-uuid"
				})).Return(nil)
			},
			want: &model.Message{
				ChatUUID: "chat-uuid",
				Role:     "user",
				Content:  "hello",
			},
			wantErr: false,
		},
		{
			name: "異常系: チャットが存在しない場合エラー",
			args: args{
				chatUUID: "non-existent",
				content:  "hello",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "non-existent").Return(nil, errors.New("not found"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系: メッセージ保存に失敗した場合エラー",
			args: args{
				chatUUID: "chat-uuid",
				content:  "hello",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(&model.Chat{UUID: "chat-uuid"}, nil)
				m.messageRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mocks{
				chatRepo:             &MockChatRepository{},
				messageRepo:          &MockMessageRepository{},
				messageSelectionRepo: &MockMessageSelectionRepository{},
				edgeRepo:             &mockEdgeRepository{},
				transactionManager:   &MockTransactionManager{},
				genaiClient:          &MockGenAIClient{},
				publisher:            &MockPublisher{},
			}
			tt.setupMock(m)

			u := NewChatUsecase(m.chatRepo, m.messageRepo, m.messageSelectionRepo, m.edgeRepo, m.transactionManager, m.genaiClient, m.publisher)

			got, err := u.SendMessage(context.Background(), tt.args.chatUUID, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("chatUsecase.SendMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want.ChatUUID, got.ChatUUID)
				assert.Equal(t, tt.want.Role, got.Role)
				assert.Equal(t, tt.want.Content, got.Content)
				assert.NotEmpty(t, got.UUID)
				assert.False(t, got.CreatedAt.IsZero())
			}
		})
	}
}

func TestChatUsecase_StreamMessage(t *testing.T) {
	type mocks struct {
		chatRepo             *MockChatRepository
		messageRepo          *MockMessageRepository
		messageSelectionRepo *MockMessageSelectionRepository
		edgeRepo             *mockEdgeRepository
		transactionManager   *MockTransactionManager
		genaiClient          *MockGenAIClient
		publisher            *MockPublisher
	}
	type args struct {
		chatUUID string
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks)
		wantErr   bool
	}{
		{
			name: "正常系: ストリームメッセージが成功すること（サマリなし）",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				// 1. FindLatestMessageWithSummary
				m.messageRepo.On("FindLatestMessageWithSummary", mock.Anything, "chat-uuid").Return(nil, nil)
				// 2. FindMessagesByChatID
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{
					{Content: "hello", Role: "user"},
				}, nil)
				// 3. GenerateContentStream
				mockIter := func(yield func(*genai.GenerateContentResponse, error) bool) {
					yield(&genai.GenerateContentResponse{
						Candidates: []*genai.Candidate{
							{
								Content: genai.Text("world")[0],
							},
						},
					}, nil)
				}
				m.genaiClient.On("GenerateContentStream", mock.Anything, "gemini-2.5-flash", mock.Anything, (*genai.GenerateContentConfig)(nil)).Return(mockIter)
				// 4. FindByID (for position)
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(&model.Chat{UUID: "chat-uuid", PositionX: 0, PositionY: 0}, nil)
				// 5. Create (Assistant Message)
				m.messageRepo.On("Create", mock.Anything, mock.MatchedBy(func(msg *model.Message) bool {
					// assistantCount is 0, so PositionY = 0
					return msg.Role == "assistant" && msg.Content == "world" && msg.ChatUUID == "chat-uuid" && msg.PositionY == 0
				})).Return(nil)
				// 6. PublishTask
				m.publisher.On("Publish", "chat_summary", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "正常系: ストリームメッセージが成功すること（サマリあり）",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				summary := "previous summary"
				summaryMsg := &model.Message{UUID: "msg-summary", ContextSummary: &summary}
				// 1. FindLatestMessageWithSummary
				m.messageRepo.On("FindLatestMessageWithSummary", mock.Anything, "chat-uuid").Return(summaryMsg, nil)
				// 2. FindMessagesByChatID
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{
					{UUID: "msg-old", Content: "old", Role: "user"},
					summaryMsg,
					{UUID: "msg-new", Content: "new", Role: "user"},
				}, nil)
				// 3. GenerateContentStream
				// プロンプトにサマリが含まれているか確認したいが、mock.Anythingで簡略化
				mockIter := func(yield func(*genai.GenerateContentResponse, error) bool) {
					yield(&genai.GenerateContentResponse{
						Candidates: []*genai.Candidate{
							{
								Content: genai.Text("response")[0],
							},
						},
					}, nil)
				}
				m.genaiClient.On("GenerateContentStream", mock.Anything, "gemini-2.5-flash", mock.Anything, (*genai.GenerateContentConfig)(nil)).Return(mockIter)
				// 4. FindByID (for position)
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(&model.Chat{UUID: "chat-uuid", PositionX: 0, PositionY: 0}, nil)
				// 5. Create (Assistant Message)
				m.messageRepo.On("Create", mock.Anything, mock.MatchedBy(func(msg *model.Message) bool {
					// assistantCount is 0, so PositionY = 0
					return msg.Role == "assistant" && msg.Content == "response" && msg.PositionY == 0
				})).Return(nil)
				// 6. PublishTask
				m.publisher.On("Publish", "chat_summary", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "異常系: 最新サマリ取得失敗",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.messageRepo.On("FindLatestMessageWithSummary", mock.Anything, "chat-uuid").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "異常系: メッセージ履歴取得失敗",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.messageRepo.On("FindLatestMessageWithSummary", mock.Anything, "chat-uuid").Return(nil, nil)
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "異常系: コンテキストなし（メッセージ0件かつサマリなし）",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.messageRepo.On("FindLatestMessageWithSummary", mock.Anything, "chat-uuid").Return(nil, nil)
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{}, nil)
			},
			wantErr: true,
		},
		{
			name: "異常系: GenAI APIからの受信エラー",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.messageRepo.On("FindLatestMessageWithSummary", mock.Anything, "chat-uuid").Return(nil, nil)
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{
					{Content: "hello", Role: "user"},
				}, nil)

				mockIter := func(yield func(*genai.GenerateContentResponse, error) bool) {
					yield(nil, errors.New("genai error"))
				}
				m.genaiClient.On("GenerateContentStream", mock.Anything, "gemini-2.5-flash", mock.Anything, (*genai.GenerateContentConfig)(nil)).Return(mockIter)
			},
			wantErr: true,
		},
		{
			name: "異常系: アシスタントメッセージ保存失敗",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.messageRepo.On("FindLatestMessageWithSummary", mock.Anything, "chat-uuid").Return(nil, nil)
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{
					{Content: "hello", Role: "user"},
				}, nil)

				mockIter := func(yield func(*genai.GenerateContentResponse, error) bool) {
					yield(&genai.GenerateContentResponse{
						Candidates: []*genai.Candidate{
							{
								Content: genai.Text("world")[0],
							},
						},
					}, nil)
				}
				m.genaiClient.On("GenerateContentStream", mock.Anything, "gemini-2.5-flash", mock.Anything, (*genai.GenerateContentConfig)(nil)).Return(mockIter)
				// FindByID (for position)
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(&model.Chat{UUID: "chat-uuid", PositionX: 0, PositionY: 0}, nil)

				m.messageRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mocks{
				chatRepo:             &MockChatRepository{},
				messageRepo:          &MockMessageRepository{},
				messageSelectionRepo: &MockMessageSelectionRepository{},
				edgeRepo:             &mockEdgeRepository{},
				transactionManager:   &MockTransactionManager{},
				genaiClient:          &MockGenAIClient{},
				publisher:            &MockPublisher{},
			}
			tt.setupMock(m)

			u := NewChatUsecase(m.chatRepo, m.messageRepo, m.messageSelectionRepo, m.edgeRepo, m.transactionManager, m.genaiClient, m.publisher)

			outputChan := make(chan string, 10)
			err := u.StreamMessage(context.Background(), tt.args.chatUUID, outputChan)

			if (err != nil) != tt.wantErr {
				t.Errorf("chatUsecase.StreamMessage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				close(outputChan)
				var output string
				for s := range outputChan {
					output += s
				}
				if tt.name == "正常系: ストリームメッセージが成功すること（サマリあり）" {
					assert.Equal(t, "response", output)
				} else {
					assert.Equal(t, "world", output)
				}
			}
		})
	}
}

func TestChatUsecase_GenerateForkPreview(t *testing.T) {
	type mocks struct {
		chatRepo             *MockChatRepository
		messageRepo          *MockMessageRepository
		messageSelectionRepo *MockMessageSelectionRepository
		edgeRepo             *mockEdgeRepository
		transactionManager   *MockTransactionManager
		genaiClient          *MockGenAIClient
		publisher            *MockPublisher
	}
	type args struct {
		chatUUID string
		req      model.ForkPreviewRequest
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks)
		want      *model.ForkPreviewResponse
		wantErr   bool
	}{
		{
			name: "正常系: フォークプレビュー生成成功",
			args: args{
				chatUUID: "chat-uuid",
				req: model.ForkPreviewRequest{
					TargetMessageUUID: "msg-2",
					SelectedText:      "selected",
					RangeStart:        0,
					RangeEnd:          8,
				},
			},
			setupMock: func(m *mocks) {
				// 1. FindMessagesByChatID
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{
					{UUID: "msg-1", Content: "hello", Role: "user"},
					{UUID: "msg-2", Content: "world", Role: "assistant"},
					{UUID: "msg-3", Content: "ignored", Role: "user"},
				}, nil)

				// 2. GenerateContent
				m.genaiClient.On("GenerateContent", mock.Anything, "gemini-2.5-flash", mock.Anything, mock.Anything).Return(&genai.GenerateContentResponse{
					Candidates: []*genai.Candidate{
						{
							Content: &genai.Content{
								Parts: []*genai.Part{
									{Text: `{"suggested_title": "New Title", "generated_context": "New Context"}`},
								},
							},
						},
					},
				}, nil)
			},
			want: &model.ForkPreviewResponse{
				SuggestedTitle:   "New Title",
				GeneratedContext: "New Context",
			},
			wantErr: false,
		},
		{
			name: "異常系: メッセージ履歴取得失敗",
			args: args{
				chatUUID: "chat-uuid",
				req: model.ForkPreviewRequest{
					TargetMessageUUID: "msg-2",
				},
			},
			setupMock: func(m *mocks) {
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系: 対象メッセージが見つからない",
			args: args{
				chatUUID: "chat-uuid",
				req: model.ForkPreviewRequest{
					TargetMessageUUID: "msg-999",
				},
			},
			setupMock: func(m *mocks) {
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{
					{UUID: "msg-1", Content: "hello", Role: "user"},
				}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系: GenAI呼び出し失敗",
			args: args{
				chatUUID: "chat-uuid",
				req: model.ForkPreviewRequest{
					TargetMessageUUID: "msg-1",
				},
			},
			setupMock: func(m *mocks) {
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{
					{UUID: "msg-1", Content: "hello", Role: "user"},
				}, nil)

				m.genaiClient.On("GenerateContent", mock.Anything, "gemini-2.5-flash", mock.Anything, mock.Anything).Return(nil, errors.New("genai error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系: JSONパース失敗",
			args: args{
				chatUUID: "chat-uuid",
				req: model.ForkPreviewRequest{
					TargetMessageUUID: "msg-1",
				},
			},
			setupMock: func(m *mocks) {
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{
					{UUID: "msg-1", Content: "hello", Role: "user"},
				}, nil)

				m.genaiClient.On("GenerateContent", mock.Anything, "gemini-2.5-flash", mock.Anything, mock.Anything).Return(&genai.GenerateContentResponse{
					Candidates: []*genai.Candidate{
						{
							Content: &genai.Content{
								Parts: []*genai.Part{
									{Text: `invalid json`},
								},
							},
						},
					},
				}, nil)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mocks{
				chatRepo:             &MockChatRepository{},
				messageRepo:          &MockMessageRepository{},
				messageSelectionRepo: &MockMessageSelectionRepository{},
				edgeRepo:             &mockEdgeRepository{},
				transactionManager:   &MockTransactionManager{},
				genaiClient:          &MockGenAIClient{},
				publisher:            &MockPublisher{},
			}
			tt.setupMock(m)

			u := NewChatUsecase(m.chatRepo, m.messageRepo, m.messageSelectionRepo, m.edgeRepo, m.transactionManager, m.genaiClient, m.publisher)

			got, err := u.GenerateForkPreview(context.Background(), tt.args.chatUUID, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("chatUsecase.GenerateForkPreview() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestChatUsecase_ForkChat(t *testing.T) {
	type mocks struct {
		chatRepo             *MockChatRepository
		messageRepo          *MockMessageRepository
		messageSelectionRepo *MockMessageSelectionRepository
		edgeRepo             *mockEdgeRepository
		transactionManager   *MockTransactionManager
		genaiClient          *MockGenAIClient
		publisher            *MockPublisher
	}
	type args struct {
		params model.ForkChatParams
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks)
		want      string
		wantErr   bool
	}{
		{
			name: "正常系: チャットフォーク成功",
			args: args{
				params: model.ForkChatParams{
					TargetMessageUUID: "msg-1",
					ParentChatUUID:    "parent-chat",
					SelectedText:      "selected",
					RangeStart:        0,
					RangeEnd:          8,
					Title:             "New Chat",
					ContextSummary:    "Summary",
				},
			},
			setupMock: func(m *mocks) {
				// 1. FindByID (Parent Chat)
				m.chatRepo.On("FindByID", mock.Anything, "parent-chat").Return(&model.Chat{
					UUID:        "parent-chat",
					ProjectUUID: "project-1",
				}, nil)

				// 2. FindByID (Target Message)
				m.messageRepo.On("FindByID", mock.Anything, "msg-1").Return(&model.Message{
					UUID:      "msg-1",
					PositionY: 100,
				}, nil)

				// 3. CountByProjectUUID
				m.chatRepo.On("CountByProjectUUID", mock.Anything, "project-1").Return(int64(5), nil)

				// 4. Transaction
				m.transactionManager.On("Do", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(context.Context) error)
					fn(context.Background())
				})

				// 5. Create MessageSelection
				m.messageSelectionRepo.On("Create", mock.Anything, mock.MatchedBy(func(s *model.MessageSelection) bool {
					return s.SelectedText == "selected"
				})).Return(nil)

				// 6. Create Chat
				m.chatRepo.On("Create", mock.Anything, mock.MatchedBy(func(c *model.Chat) bool {
					// PositionX = 5 * 250 = 1250, PositionY = 100
					return c.Title == "New Chat" && c.ProjectUUID == "project-1" && *c.ParentUUID == "parent-chat" && c.PositionX == 1250 && c.PositionY == 100
				})).Return(nil)

				// 7. Create Message
				m.messageRepo.On("Create", mock.Anything, mock.MatchedBy(func(msg *model.Message) bool {
					return msg.Role == "assistant" && msg.Content == "New Chat\n\nSummary" && msg.PositionX == 1250 && msg.PositionY == 100
				})).Return(nil)

				// 8. Create Edge
				m.edgeRepo.On("Create", mock.Anything, mock.MatchedBy(func(edge *model.Edge) bool {
					// TargetMessageUUID is "msg-1"
					return edge.TargetMessageUUID == "msg-1"
				})).Return(nil)
			},
			want:    "new-chat-id", // UUIDはランダム生成なので、空文字でないことを確認する
			wantErr: false,
		},
		{
			name: "異常系: 親チャットが見つからない",
			args: args{
				params: model.ForkChatParams{
					ParentChatUUID: "parent-chat",
				},
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "parent-chat").Return(nil, errors.New("not found"))
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "異常系: トランザクションエラー",
			args: args{
				params: model.ForkChatParams{
					ParentChatUUID:    "parent-chat",
					TargetMessageUUID: "msg-1",
				},
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "parent-chat").Return(&model.Chat{
					UUID:        "parent-chat",
					ProjectUUID: "project-1",
				}, nil)

				m.messageRepo.On("FindByID", mock.Anything, "msg-1").Return(&model.Message{
					UUID:      "msg-1",
					PositionY: 100,
				}, nil)

				m.chatRepo.On("CountByProjectUUID", mock.Anything, "project-1").Return(int64(5), nil)

				m.transactionManager.On("Do", mock.Anything, mock.Anything).Return(errors.New("tx error"))
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mocks{
				chatRepo:             &MockChatRepository{},
				messageRepo:          &MockMessageRepository{},
				messageSelectionRepo: &MockMessageSelectionRepository{},
				edgeRepo:             &mockEdgeRepository{},
				transactionManager:   &MockTransactionManager{},
				genaiClient:          &MockGenAIClient{},
				publisher:            &MockPublisher{},
			}
			tt.setupMock(m)

			u := NewChatUsecase(m.chatRepo, m.messageRepo, m.messageSelectionRepo, m.edgeRepo, m.transactionManager, m.genaiClient, m.publisher)

			got, err := u.ForkChat(context.Background(), tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("chatUsecase.ForkChat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotEmpty(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestChatUsecase_GetMergePreview(t *testing.T) {
	type mocks struct {
		chatRepo             *MockChatRepository
		messageRepo          *MockMessageRepository
		messageSelectionRepo *MockMessageSelectionRepository
		edgeRepo             *mockEdgeRepository
		transactionManager   *MockTransactionManager
		genaiClient          *MockGenAIClient
		publisher            *MockPublisher
	}
	type args struct {
		chatUUID string
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks)
		want      *model.MergePreview
		wantErr   bool
	}{
		{
			name: "正常系: マージプレビューが生成できること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(&model.Chat{
					UUID:           "chat-uuid",
					ContextSummary: "parent context",
				}, nil)
				m.messageRepo.On("FindLatestMessageWithSummary", mock.Anything, "chat-uuid").Return(&model.Message{
					ContextSummary: func() *string { s := "child summary"; return &s }(),
				}, nil)
				m.messageRepo.On("FindLatestMessageByRole", mock.Anything, "chat-uuid", "assistant").Return(&model.Message{
					Content: "latest assistant message",
				}, nil)

				m.genaiClient.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&genai.GenerateContentResponse{
					Candidates: []*genai.Candidate{
						{
							Content: &genai.Content{
								Parts: []*genai.Part{
									{Text: "summary"},
								},
							},
						},
					},
				}, nil)
			},
			want: &model.MergePreview{
				SuggestedSummary: "summary",
			},
			wantErr: false,
		},
		{
			name: "異常系: チャット取得失敗",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mocks{
				chatRepo:             &MockChatRepository{},
				messageRepo:          &MockMessageRepository{},
				messageSelectionRepo: &MockMessageSelectionRepository{},
				edgeRepo:             &mockEdgeRepository{},
				transactionManager:   &MockTransactionManager{},
				genaiClient:          &MockGenAIClient{},
				publisher:            &MockPublisher{},
			}
			tt.setupMock(m)

			u := NewChatUsecase(m.chatRepo, m.messageRepo, m.messageSelectionRepo, m.edgeRepo, m.transactionManager, m.genaiClient, m.publisher)

			got, err := u.GetMergePreview(context.Background(), tt.args.chatUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("chatUsecase.GetMergePreview() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestChatUsecase_MergeChat(t *testing.T) {
	type mocks struct {
		chatRepo             *MockChatRepository
		messageRepo          *MockMessageRepository
		messageSelectionRepo *MockMessageSelectionRepository
		edgeRepo             *mockEdgeRepository
		transactionManager   *MockTransactionManager
		genaiClient          *MockGenAIClient
		publisher            *MockPublisher
	}
	type args struct {
		chatUUID string
		params   model.MergeChatParams
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks)
		want      *model.MergeChatResult
		wantErr   bool
	}{
		{
			name: "正常系: チャットマージ成功",
			args: args{
				chatUUID: "child-chat-uuid",
				params: model.MergeChatParams{
					ParentChatUUID: "parent-chat-uuid",
					SummaryContent: "summary content",
				},
			},
			setupMock: func(m *mocks) {
				sourceMsgUUID := "source-msg-uuid"
				// 1. FindByID (Child Chat)
				m.chatRepo.On("FindByID", mock.Anything, "child-chat-uuid").Return(&model.Chat{
					UUID:              "child-chat-uuid",
					SourceMessageUUID: &sourceMsgUUID,
				}, nil)

				// 2. Transaction
				m.transactionManager.On("Do", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(context.Context) error)
					fn(context.Background())
				})

				// 3. Create (Report Message)
				m.messageRepo.On("Create", mock.Anything, mock.MatchedBy(func(msg *model.Message) bool {
					return msg.Role == "merge_report" &&
						msg.ChatUUID == "parent-chat-uuid" &&
						msg.ParentMessageUUID != nil && *msg.ParentMessageUUID == "source-msg-uuid" &&
						*msg.SourceChatUUID == "child-chat-uuid"
				})).Return(nil)

				// 4. UpdateStatus
				m.chatRepo.On("UpdateStatus", mock.Anything, "child-chat-uuid", "merged").Return(nil)
			},
			want: &model.MergeChatResult{
				SummaryContent: "summary content",
			},
			wantErr: false,
		},
		{
			name: "異常系: 子チャット取得失敗",
			args: args{
				chatUUID: "child-chat-uuid",
				params: model.MergeChatParams{
					ParentChatUUID: "parent-chat-uuid",
				},
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "child-chat-uuid").Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系: ソースメッセージがない場合",
			args: args{
				chatUUID: "child-chat-uuid",
				params: model.MergeChatParams{
					ParentChatUUID: "parent-chat-uuid",
				},
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "child-chat-uuid").Return(&model.Chat{
					UUID:              "child-chat-uuid",
					SourceMessageUUID: nil,
				}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系: トランザクションエラー",
			args: args{
				chatUUID: "child-chat-uuid",
				params: model.MergeChatParams{
					ParentChatUUID: "parent-chat-uuid",
					SummaryContent: "summary content",
				},
			},
			setupMock: func(m *mocks) {
				sourceMsgUUID := "source-msg-uuid"
				m.chatRepo.On("FindByID", mock.Anything, "child-chat-uuid").Return(&model.Chat{
					UUID:              "child-chat-uuid",
					SourceMessageUUID: &sourceMsgUUID,
				}, nil)

				m.transactionManager.On("Do", mock.Anything, mock.Anything).Return(errors.New("tx error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mocks{
				chatRepo:             &MockChatRepository{},
				messageRepo:          &MockMessageRepository{},
				messageSelectionRepo: &MockMessageSelectionRepository{},
				edgeRepo:             &mockEdgeRepository{},
				transactionManager:   &MockTransactionManager{},
				genaiClient:          &MockGenAIClient{},
				publisher:            &MockPublisher{},
			}
			tt.setupMock(m)

			u := NewChatUsecase(m.chatRepo, m.messageRepo, m.messageSelectionRepo, m.edgeRepo, m.transactionManager, m.genaiClient, m.publisher)

			got, err := u.MergeChat(context.Background(), tt.args.chatUUID, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("chatUsecase.MergeChat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want.SummaryContent, got.SummaryContent)
				assert.NotEmpty(t, got.ReportMessageID)
			}
		})
	}
}

func TestChatUsecase_CloseChat(t *testing.T) {
	type mocks struct {
		chatRepo             *MockChatRepository
		messageRepo          *MockMessageRepository
		messageSelectionRepo *MockMessageSelectionRepository
		edgeRepo             *mockEdgeRepository
		transactionManager   *MockTransactionManager
		genaiClient          *MockGenAIClient
		publisher            *MockPublisher
	}
	type args struct {
		chatUUID string
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks)
		want      string
		wantErr   bool
	}{
		{
			name: "正常系: チャットのクローズが成功すること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(&model.Chat{UUID: "chat-uuid"}, nil)
				m.chatRepo.On("UpdateStatus", mock.Anything, "chat-uuid", "closed").Return(nil)
			},
			want:    "chat-uuid",
			wantErr: false,
		},
		{
			name: "異常系: チャットが存在しない場合エラー",
			args: args{
				chatUUID: "non-existent",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "non-existent").Return(nil, errors.New("not found"))
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "異常系: ステータス更新に失敗した場合エラー",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(&model.Chat{UUID: "chat-uuid"}, nil)
				m.chatRepo.On("UpdateStatus", mock.Anything, "chat-uuid", "closed").Return(errors.New("db error"))
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mocks{
				chatRepo:             &MockChatRepository{},
				messageRepo:          &MockMessageRepository{},
				messageSelectionRepo: &MockMessageSelectionRepository{},
				edgeRepo:             &mockEdgeRepository{},
				transactionManager:   &MockTransactionManager{},
				genaiClient:          &MockGenAIClient{},
				publisher:            &MockPublisher{},
			}
			tt.setupMock(m)

			u := NewChatUsecase(m.chatRepo, m.messageRepo, m.messageSelectionRepo, m.edgeRepo, m.transactionManager, m.genaiClient, m.publisher)

			got, err := u.CloseChat(context.Background(), tt.args.chatUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("chatUsecase.CloseChat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestChatUsecase_OpenChat(t *testing.T) {
	type mocks struct {
		chatRepo             *MockChatRepository
		messageRepo          *MockMessageRepository
		messageSelectionRepo *MockMessageSelectionRepository
		edgeRepo             *mockEdgeRepository
		transactionManager   *MockTransactionManager
		genaiClient          *MockGenAIClient
		publisher            *MockPublisher
	}
	type args struct {
		chatUUID string
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks)
		want      string
		wantErr   bool
	}{
		{
			name: "正常系: チャットオープンが成功すること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(&model.Chat{UUID: "chat-uuid"}, nil)
				m.chatRepo.On("UpdateStatus", mock.Anything, "chat-uuid", "open").Return(nil)
			},
			want:    "chat-uuid",
			wantErr: false,
		},
		{
			name: "異常系: チャットが存在しない場合エラー",
			args: args{
				chatUUID: "non-existent",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "non-existent").Return(nil, errors.New("not found"))
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "異常系: ステータス更新に失敗した場合エラー",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatRepo.On("FindByID", mock.Anything, "chat-uuid").Return(&model.Chat{UUID: "chat-uuid"}, nil)
				m.chatRepo.On("UpdateStatus", mock.Anything, "chat-uuid", "open").Return(errors.New("db error"))
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mocks{
				chatRepo:             &MockChatRepository{},
				messageRepo:          &MockMessageRepository{},
				messageSelectionRepo: &MockMessageSelectionRepository{},
				edgeRepo:             &mockEdgeRepository{},
				transactionManager:   &MockTransactionManager{},
				genaiClient:          &MockGenAIClient{},
				publisher:            &MockPublisher{},
			}
			tt.setupMock(m)

			u := NewChatUsecase(m.chatRepo, m.messageRepo, m.messageSelectionRepo, m.edgeRepo, m.transactionManager, m.genaiClient, m.publisher)

			got, err := u.OpenChat(context.Background(), tt.args.chatUUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("chatUsecase.OpenChat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
