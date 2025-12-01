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

func TestChatUsecase_FirstStreamChat(t *testing.T) {
	type mocks struct {
		chatRepo    *MockChatRepository
		messageRepo *MockMessageRepository
		genaiClient *MockGenAIClient
		publisher   *MockPublisher
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
					return msg.Role == "assistant" && msg.Content == "world" && msg.ChatUUID == "chat-uuid"
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
				chatRepo:    &MockChatRepository{},
				messageRepo: &MockMessageRepository{},
				genaiClient: &MockGenAIClient{},
				publisher:   &MockPublisher{},
			}
			tt.setupMock(m)

			u := NewChatUsecase(m.chatRepo, m.messageRepo, m.genaiClient, m.publisher)

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
		chatRepo    *MockChatRepository
		messageRepo *MockMessageRepository
		genaiClient *MockGenAIClient
		publisher   *MockPublisher
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
				chatRepo:    &MockChatRepository{},
				messageRepo: &MockMessageRepository{},
				genaiClient: &MockGenAIClient{},
				publisher:   &MockPublisher{},
			}
			tt.setupMock(m)

			u := NewChatUsecase(m.chatRepo, m.messageRepo, m.genaiClient, m.publisher)

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
		chatRepo    *MockChatRepository
		messageRepo *MockMessageRepository
		genaiClient *MockGenAIClient
		publisher   *MockPublisher
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
				chatRepo:    &MockChatRepository{},
				messageRepo: &MockMessageRepository{},
				genaiClient: &MockGenAIClient{},
				publisher:   &MockPublisher{},
			}
			tt.setupMock(m)

			u := NewChatUsecase(m.chatRepo, m.messageRepo, m.genaiClient, m.publisher)

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
		chatRepo    *MockChatRepository
		messageRepo *MockMessageRepository
		genaiClient *MockGenAIClient
		publisher   *MockPublisher
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
				chatRepo:    &MockChatRepository{},
				messageRepo: &MockMessageRepository{},
				genaiClient: &MockGenAIClient{},
				publisher:   &MockPublisher{},
			}
			tt.setupMock(m)

			u := NewChatUsecase(m.chatRepo, m.messageRepo, m.genaiClient, m.publisher)

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
		chatRepo    *MockChatRepository
		messageRepo *MockMessageRepository
		genaiClient *MockGenAIClient
		publisher   *MockPublisher
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
				// 4. Create (Assistant Message)
				m.messageRepo.On("Create", mock.Anything, mock.MatchedBy(func(msg *model.Message) bool {
					return msg.Role == "assistant" && msg.Content == "world" && msg.ChatUUID == "chat-uuid"
				})).Return(nil)
				// 5. PublishTask
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
				// 4. Create (Assistant Message)
				m.messageRepo.On("Create", mock.Anything, mock.MatchedBy(func(msg *model.Message) bool {
					return msg.Role == "assistant" && msg.Content == "response"
				})).Return(nil)
				// 5. PublishTask
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

				m.messageRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mocks{
				chatRepo:    &MockChatRepository{},
				messageRepo: &MockMessageRepository{},
				genaiClient: &MockGenAIClient{},
				publisher:   &MockPublisher{},
			}
			tt.setupMock(m)

			u := NewChatUsecase(m.chatRepo, m.messageRepo, m.genaiClient, m.publisher)

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
