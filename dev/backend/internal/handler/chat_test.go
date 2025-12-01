package handler

import (
	"backend/internal/domain/model"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockChatUsecase
type MockChatUsecase struct {
	mock.Mock
}

func (m *MockChatUsecase) FirstStreamChat(ctx context.Context, chatUUID string, outputChan chan<- string) error {
	args := m.Called(ctx, chatUUID, outputChan)
	// ストリーム処理のシミュレーション
	if fn, ok := args.Get(0).(func(chan<- string)); ok && fn != nil {
		// 非同期で実行しないと、ハンドラがチャネル待ちでブロックする可能性があるが、
		// ハンドラの実装では go routine で usecase を呼んでいるので、
		// ここでは同期的に書き込んでもいいかもしれないが、
		// ハンドラのテストでは usecase は go routine で呼ばれる。
		// しかし、FirstStreamChat はエラーを返すかどうかなので、
		// 成功時はチャネルに書き込む処理をここで実行する。
		fn(outputChan)
	}
	return args.Error(1)
}

func (m *MockChatUsecase) GetChat(ctx context.Context, chatUUID string) (*model.Chat, error) {
	args := m.Called(ctx, chatUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chat), args.Error(1)
}

func TestChatHandler_FirstStreamChat(t *testing.T) {
	type mocks struct {
		chatUsecase *MockChatUsecase
	}
	type args struct {
		chatUUID string
	}
	tests := []struct {
		name       string
		args       args
		setupMock  func(m *mocks)
		wantStatus int
		wantBody   string
	}{
		{
			name: "正常系: ストリームレスポンスが返ること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatUsecase.On("FirstStreamChat", mock.Anything, "chat-uuid", mock.Anything).Return(func(ch chan<- string) {
					ch <- "hello"
					ch <- "world"
				}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `data: {"chunk":"hello","status":"processing"}` + "\n\n\n" + `data: {"chunk":"world","status":"processing"}` + "\n\n\n" + `data: {"status":"done"}` + "\n\n\n",
		},
		{
			name: "異常系: Usecaseがエラーを返した場合",
			args: args{
				chatUUID: "error-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatUsecase.On("FirstStreamChat", mock.Anything, "error-uuid", mock.Anything).Return(nil, errors.New("usecase error"))
			},
			// ハンドラの実装では、エラーチャネルからエラーを受け取ると return err する。
			// Echo のハンドラがエラーを返すと、ミドルウェアがそれを処理するが、
			// ここではハンドラ関数自体の戻り値エラーは検証しにくい（ServeHTTP経由だとレスポンスに反映されるか？）
			// 今回の実装では、SSEのループ内でエラーが発生するとループを抜けてエラーを返す。
			// ただし、すでにヘッダーは200で送信済みなので、ステータスコードは200になる可能性がある。
			// エラー時の挙動は実装依存だが、ここではステータスコード200で、ボディが途中で終わるか、あるいはエラーが返るか。
			// httptest.Recorder は WriteHeader が呼ばれた後の変更を記録する。
			wantStatus: http.StatusOK,
			// エラー時のボディは実装によるが、今回はエラーハンドリングが SSE のストリーム中に行われるため、
			// クライアントにはエラーが伝わらない（接続が切れる）か、ログが出るだけかもしれない。
			// テストとしては、関数がエラーを返すことを確認すべきだが、EchoのServeHTTPを使うとエラーは飲み込まれるか、HTTPエラーになる。
			// ここでは簡略化のため、ステータスコードのみ確認するが、厳密にはハンドラ関数の戻り値をテストすべき。
			// しかし Echo のテストパラダイムでは ServeHTTP を使うのが一般的。
			wantBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/chats/"+tt.args.chatUUID+"/stream", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/chats/:chat_uuid/stream")
			c.SetParamNames("chat_uuid")
			c.SetParamValues(tt.args.chatUUID)

			m := &mocks{
				chatUsecase: &MockChatUsecase{},
			}
			tt.setupMock(m)

			h := NewChatHandler(m.chatUsecase)

			// ハンドラの実行
			// エラーが返る場合もある
			err := h.FirstStreamChat(c)

			if tt.name == "異常系: Usecaseがエラーを返した場合" {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatus, rec.Code)
				// ボディの検証
				assert.Equal(t, tt.wantBody, rec.Body.String())
			}
		})
	}
}

func TestChatHandler_GetChat(t *testing.T) {
	type mocks struct {
		chatUsecase *MockChatUsecase
	}
	type args struct {
		chatUUID string
	}
	tests := []struct {
		name       string
		args       args
		setupMock  func(m *mocks)
		wantStatus int
		wantBody   string
	}{
		{
			name: "正常系: チャットが取得できること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				parentUUID := "parent-uuid"
				m.chatUsecase.On("GetChat", mock.Anything, "chat-uuid").Return(&model.Chat{
					UUID:           "chat-uuid",
					ProjectUUID:    "project-uuid",
					ParentUUID:     &parentUUID,
					Title:          "test chat",
					Status:         "active",
					ContextSummary: "summary",
				}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"uuid":"chat-uuid","project_uuid":"project-uuid","parent_uuid":"parent-uuid","title":"test chat","status":"active","context_summary":"summary"}`,
		},
		{
			name: "異常系: Usecaseがエラーを返した場合",
			args: args{
				chatUUID: "error-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatUsecase.On("GetChat", mock.Anything, "error-uuid").Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"status":"error","message":"db error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/chats/"+tt.args.chatUUID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/chats/:chat_uuid")
			c.SetParamNames("chat_uuid")
			c.SetParamValues(tt.args.chatUUID)

			m := &mocks{
				chatUsecase: &MockChatUsecase{},
			}
			tt.setupMock(m)

			h := NewChatHandler(m.chatUsecase)
			err := h.GetChat(c)

			if tt.wantStatus >= 400 {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatus, rec.Code)
				assert.JSONEq(t, tt.wantBody, rec.Body.String())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatus, rec.Code)
				assert.JSONEq(t, tt.wantBody, rec.Body.String())
			}
		})
	}
}
