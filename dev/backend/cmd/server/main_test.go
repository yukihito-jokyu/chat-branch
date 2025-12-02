package main

import (
	"backend/config"
	"backend/internal/worker"
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/glebarez/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/genai"
	"gorm.io/gorm"
)

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

type MockSubscriber struct {
	mock.Mock
}

func (m *MockSubscriber) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	args := m.Called(ctx, topic)
	return args.Get(0).(<-chan *message.Message), args.Error(1)
}

func (m *MockSubscriber) Close() error {
	args := m.Called()
	return args.Error(0)
}

func Test_setupServer(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T) (*config.Config, *gorm.DB, *genai.Client, message.Publisher)
		assertion func(t *testing.T, e *echo.Echo)
	}{
		{
			name: "正常系: サーバーが正しく初期化される",
			setup: func(t *testing.T) (*config.Config, *gorm.DB, *genai.Client, message.Publisher) {
				cfg := &config.Config{
					Server: config.ServerConfig{Address: ":8080"},
					Gemini: config.GeminiConfig{APIKey: "dummy"},
				}
				db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				assert.NoError(t, err)
				genaiClient, err := genai.NewClient(context.Background(), &genai.ClientConfig{APIKey: "dummy"})
				assert.NoError(t, err)
				mockPublisher := new(MockPublisher)
				return cfg, db, genaiClient, mockPublisher
			},
			assertion: func(t *testing.T, e *echo.Echo) {
				assert.NotNil(t, e)
				routes := e.Routes()
				assert.NotEmpty(t, routes)

				// InitRoutesが呼び出されたことを確認するために特定のルートをチェックする
				found := false
				for _, r := range routes {
					if r.Path == "/health" {
						found = true
						break
					}
				}
				assert.True(t, found, "ヘルスチェックルートが登録されている必要があります")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, db, genaiClient, publisher := tt.setup(t)
			e := setupServer(cfg, db, genaiClient, publisher)
			tt.assertion(t, e)
		})
	}
}

func Test_setupWorker(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T) (*gorm.DB, *genai.Client, message.Subscriber)
		assertion func(t *testing.T, w *worker.SummaryWorker)
	}{
		{
			name: "正常系: Workerが正しく初期化される",
			setup: func(t *testing.T) (*gorm.DB, *genai.Client, message.Subscriber) {
				db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				assert.NoError(t, err)
				genaiClient, err := genai.NewClient(context.Background(), &genai.ClientConfig{APIKey: "dummy"})
				assert.NoError(t, err)
				mockSubscriber := new(MockSubscriber)
				return db, genaiClient, mockSubscriber
			},
			assertion: func(t *testing.T, w *worker.SummaryWorker) {
				assert.NotNil(t, w)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, genaiClient, subscriber := tt.setup(t)
			w := setupWorker(db, genaiClient, subscriber)
			tt.assertion(t, w)
		})
	}
}
