package worker

import (
	"backend/internal/domain/model"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/genai"
)

// Mocks
type MockSubscriber struct {
	mock.Mock
}

func (m *MockSubscriber) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	args := m.Called(ctx, topic)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(<-chan *message.Message), args.Error(1)
}

func (m *MockSubscriber) Close() error {
	args := m.Called()
	return args.Error(0)
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

func TestSummaryWorker_Handle(t *testing.T) {
	type mocks struct {
		subscriber  *MockSubscriber
		messageRepo *MockMessageRepository
		genaiClient *MockGenAIClient
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
			name: "正常系: 要約生成と保存が成功すること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				// 1. FindLatestMessageWithSummary
				m.messageRepo.On("FindLatestMessageWithSummary", mock.Anything, "chat-uuid").Return((*model.Message)(nil), nil)

				// 2. FindMessagesByChatID
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{
					{UUID: "msg-1", Content: "hello", Role: "user"},
					{UUID: "msg-2", Content: "world", Role: "assistant"},
				}, nil)

				// 3. GenerateContent
				resp := &genai.GenerateContentResponse{
					Candidates: []*genai.Candidate{
						{
							Content: genai.Text("summary content")[0],
						},
					},
				}
				m.genaiClient.On("GenerateContent", mock.Anything, "gemini-2.5-flash", mock.Anything, (*genai.GenerateContentConfig)(nil)).Return(resp, nil)

				// 4. UpdateContextSummary
				m.messageRepo.On("UpdateContextSummary", mock.Anything, "msg-2", "summary content").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "異常系: メッセージ取得失敗",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.messageRepo.On("FindLatestMessageWithSummary", mock.Anything, "chat-uuid").Return((*model.Message)(nil), nil)
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "正常系: メッセージが0件の場合は何もしない",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.messageRepo.On("FindLatestMessageWithSummary", mock.Anything, "chat-uuid").Return((*model.Message)(nil), nil)
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{}, nil)
			},
			wantErr: false,
		},
		{
			name: "異常系: GenAIエラー",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.messageRepo.On("FindLatestMessageWithSummary", mock.Anything, "chat-uuid").Return((*model.Message)(nil), nil)
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{
					{UUID: "msg-1", Content: "hello", Role: "user"},
				}, nil)

				m.genaiClient.On("GenerateContent", mock.Anything, "gemini-2.5-flash", mock.Anything, (*genai.GenerateContentConfig)(nil)).Return(nil, errors.New("genai error"))
			},
			wantErr: true,
		},
		{
			name: "異常系: 要約保存失敗",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.messageRepo.On("FindLatestMessageWithSummary", mock.Anything, "chat-uuid").Return((*model.Message)(nil), nil)
				m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{
					{UUID: "msg-1", Content: "hello", Role: "user"},
				}, nil)

				resp := &genai.GenerateContentResponse{
					Candidates: []*genai.Candidate{
						{
							Content: genai.Text("summary")[0],
						},
					},
				}
				m.genaiClient.On("GenerateContent", mock.Anything, "gemini-2.5-flash", mock.Anything, (*genai.GenerateContentConfig)(nil)).Return(resp, nil)

				m.messageRepo.On("UpdateContextSummary", mock.Anything, "msg-1", "summary").Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mocks{
				subscriber:  &MockSubscriber{},
				messageRepo: &MockMessageRepository{},
				genaiClient: &MockGenAIClient{},
			}
			tt.setupMock(m)

			w := NewSummaryWorker(m.subscriber, m.messageRepo, m.genaiClient)

			// JSON marshal the chatUUID
			payload, _ := json.Marshal(tt.args.chatUUID)
			msg := message.NewMessage("msg-uuid", payload)
			err := w.Handle(context.Background(), msg)

			if (err != nil) != tt.wantErr {
				t.Errorf("SummaryWorker.Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSummaryWorker_Run(t *testing.T) {
	// Runメソッドのテストは、Subscriberからのメッセージ受信とHandle呼び出しのループを確認する
	// 簡易的なテストとして、1つのメッセージを処理して終了することを確認する
	// ただし、Runは無限ループするので、Contextでキャンセルするか、チャネルを閉じる必要がある

	m := &mocks{
		subscriber:  &MockSubscriber{},
		messageRepo: &MockMessageRepository{},
		genaiClient: &MockGenAIClient{},
	}

	// Subscribeのモック
	msgChan := make(chan *message.Message, 1)
	payload, _ := json.Marshal("chat-uuid")
	msg := message.NewMessage("msg-uuid", payload)
	msgChan <- msg
	close(msgChan) // チャネルを閉じてループを終了させる

	m.subscriber.On("Subscribe", mock.Anything, "chat_summary").Return((<-chan *message.Message)(msgChan), nil)

	// Handleの内部処理のモック (FindMessagesByChatIDなど)
	m.messageRepo.On("FindLatestMessageWithSummary", mock.Anything, "chat-uuid").Return((*model.Message)(nil), nil)
	m.messageRepo.On("FindMessagesByChatID", mock.Anything, "chat-uuid").Return([]*model.Message{}, nil) // 0件で即終了

	w := NewSummaryWorker(m.subscriber, m.messageRepo, m.genaiClient)

	err := w.Run(context.Background())
	assert.NoError(t, err)

	// Ackが呼ばれたか確認 (Handleが成功した場合)
	select {
	case <-msg.Acked():
		// OK
	case <-time.After(100 * time.Millisecond):
		t.Error("message should be acked")
	}
}

// mocks struct for helper
type mocks struct {
	subscriber  *MockSubscriber
	messageRepo *MockMessageRepository
	genaiClient *MockGenAIClient
}
