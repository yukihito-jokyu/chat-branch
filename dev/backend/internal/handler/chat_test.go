package handler

import (
	"backend/internal/domain/model"
	handlerModel "backend/internal/handler/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

func (m *MockChatUsecase) StreamMessage(ctx context.Context, chatUUID string, outputChan chan<- string) error {
	args := m.Called(ctx, chatUUID, outputChan)
	if fn, ok := args.Get(0).(func(chan<- string)); ok && fn != nil {
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

func (m *MockChatUsecase) GetMessages(ctx context.Context, chatUUID string) ([]*model.Message, error) {
	args := m.Called(ctx, chatUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Message), args.Error(1)
}

func (m *MockChatUsecase) SendMessage(ctx context.Context, chatUUID string, content string) (*model.Message, error) {
	args := m.Called(ctx, chatUUID, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Message), args.Error(1)
}

func (m *MockChatUsecase) GenerateForkPreview(ctx context.Context, chatUUID string, req model.ForkPreviewRequest) (*model.ForkPreviewResponse, error) {
	args := m.Called(ctx, chatUUID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ForkPreviewResponse), args.Error(1)
}

func (m *MockChatUsecase) ForkChat(ctx context.Context, params model.ForkChatParams) (string, error) {
	args := m.Called(ctx, params)
	return args.String(0), args.Error(1)
}

func (m *MockChatUsecase) GetMergePreview(ctx context.Context, chatUUID string) (*model.MergePreview, error) {
	args := m.Called(ctx, chatUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MergePreview), args.Error(1)
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
			wantStatus: http.StatusOK,
			wantBody:   "",
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

func TestChatHandler_GetMessages(t *testing.T) {
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
			name: "正常系: メッセージ一覧が取得できること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatUsecase.On("GetMessages", mock.Anything, "chat-uuid").Return([]*model.Message{
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
			wantStatus: http.StatusOK,
			wantBody:   `[{"uuid":"msg-1","role":"user","content":"hello","forks":[{"chat_uuid":"child-chat","selected_text":"hello","range_start":0,"range_end":5}]}]`,
		},
		{
			name: "異常系: Usecaseがエラーを返した場合",
			args: args{
				chatUUID: "error-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatUsecase.On("GetMessages", mock.Anything, "error-uuid").Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"status":"error","message":"db error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/chats/"+tt.args.chatUUID+"/messages", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/chats/:chat_uuid/messages")
			c.SetParamNames("chat_uuid")
			c.SetParamValues(tt.args.chatUUID)

			m := &mocks{
				chatUsecase: &MockChatUsecase{},
			}
			tt.setupMock(m)

			h := NewChatHandler(m.chatUsecase)
			err := h.GetMessages(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

func TestChatHandler_SendMessage(t *testing.T) {
	type mocks struct {
		chatUsecase *MockChatUsecase
	}
	type args struct {
		chatUUID string
		body     string
	}
	tests := []struct {
		name       string
		args       args
		setupMock  func(m *mocks)
		wantStatus int
		wantBody   string
	}{
		{
			name: "正常系: メッセージ送信が成功すること",
			args: args{
				chatUUID: "chat-uuid",
				body:     `{"content": "hello"}`,
			},
			setupMock: func(m *mocks) {
				m.chatUsecase.On("SendMessage", mock.Anything, "chat-uuid", "hello").Return(&model.Message{
					UUID:           "msg-uuid",
					ChatUUID:       "chat-uuid",
					Role:           "user",
					Content:        "hello",
					SourceChatUUID: nil,
				}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"uuid":"msg-uuid","role":"user","content":"hello","forks":[]}`,
		},
		{
			name: "異常系: リクエストボディが不正な場合",
			args: args{
				chatUUID: "chat-uuid",
				body:     `invalid json`,
			},
			setupMock: func(m *mocks) {
				// 呼ばれない
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"status":"error","message":"リクエストボディのバインドに失敗しました"}`,
		},
		{
			name: "異常系: コンテンツが空の場合",
			args: args{
				chatUUID: "chat-uuid",
				body:     `{"content": ""}`,
			},
			setupMock: func(m *mocks) {
				// 呼ばれない
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"status":"error","message":"content が空です"}`,
		},
		{
			name: "異常系: Usecaseがエラーを返した場合",
			args: args{
				chatUUID: "chat-uuid",
				body:     `{"content": "hello"}`,
			},
			setupMock: func(m *mocks) {
				m.chatUsecase.On("SendMessage", mock.Anything, "chat-uuid", "hello").Return(nil, errors.New("db error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"status":"error","message":"db error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/chats/"+tt.args.chatUUID+"/message", strings.NewReader(tt.args.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/chats/:chat_uuid/message")
			c.SetParamNames("chat_uuid")
			c.SetParamValues(tt.args.chatUUID)

			m := &mocks{
				chatUsecase: &MockChatUsecase{},
			}
			tt.setupMock(m)

			h := NewChatHandler(m.chatUsecase)
			err := h.SendMessage(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

func TestChatHandler_StreamMessage(t *testing.T) {
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
				m.chatUsecase.On("StreamMessage", mock.Anything, "chat-uuid", mock.Anything).Return(func(ch chan<- string) {
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
				m.chatUsecase.On("StreamMessage", mock.Anything, "error-uuid", mock.Anything).Return(nil, errors.New("usecase error"))
			},
			wantStatus: http.StatusOK,
			wantBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/chats/"+tt.args.chatUUID+"/stream_message", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/chats/:chat_uuid/stream_message")
			c.SetParamNames("chat_uuid")
			c.SetParamValues(tt.args.chatUUID)

			m := &mocks{
				chatUsecase: &MockChatUsecase{},
			}
			tt.setupMock(m)

			h := NewChatHandler(m.chatUsecase)

			// ハンドラの実行
			err := h.StreamMessage(c)

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

func TestChatHandler_GenerateForkPreview(t *testing.T) {
	type mocks struct {
		chatUsecase *MockChatUsecase
	}
	type args struct {
		chatUUID string
		body     string
	}
	tests := []struct {
		name       string
		args       args
		setupMock  func(m *mocks)
		wantStatus int
		wantBody   string
	}{
		{
			name: "正常系: プレビュー生成が成功すること",
			args: args{
				chatUUID: "chat-uuid",
				body:     `{"target_message_uuid": "msg-uuid", "selected_text": "hello", "range_start": 0, "range_end": 5}`,
			},
			setupMock: func(m *mocks) {
				m.chatUsecase.On("GenerateForkPreview", mock.Anything, "chat-uuid", model.ForkPreviewRequest{
					TargetMessageUUID: "msg-uuid",
					SelectedText:      "hello",
					RangeStart:        0,
					RangeEnd:          5,
				}).Return(&model.ForkPreviewResponse{
					SuggestedTitle:   "Suggested Title",
					GeneratedContext: "Generated Context",
				}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"suggested_title":"Suggested Title","generated_context":"Generated Context"}`,
		},
		{
			name: "異常系: リクエストボディが不正な場合",
			args: args{
				chatUUID: "chat-uuid",
				body:     `invalid json`,
			},
			setupMock: func(m *mocks) {
				// 呼ばれない
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"status":"error","message":"リクエストボディのバインドに失敗しました"}`,
		},
		{
			name: "異常系: Usecaseがエラーを返した場合",
			args: args{
				chatUUID: "chat-uuid",
				body:     `{"target_message_uuid": "msg-uuid", "selected_text": "hello", "range_start": 0, "range_end": 5}`,
			},
			setupMock: func(m *mocks) {
				m.chatUsecase.On("GenerateForkPreview", mock.Anything, "chat-uuid", mock.Anything).Return(nil, errors.New("genai error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"status":"error","message":"genai error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/chats/"+tt.args.chatUUID+"/fork/preview", strings.NewReader(tt.args.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/chats/:chat_uuid/fork/preview")
			c.SetParamNames("chat_uuid")
			c.SetParamValues(tt.args.chatUUID)

			m := &mocks{
				chatUsecase: &MockChatUsecase{},
			}
			tt.setupMock(m)

			h := NewChatHandler(m.chatUsecase)
			err := h.GenerateForkPreview(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

func TestChatHandler_ForkChat(t *testing.T) {
	e := echo.New()
	type mocks struct {
		chatUsecase *MockChatUsecase
	}
	type args struct {
		chatUUID string
		req      handlerModel.ForkChatRequest
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks)
		want      *handlerModel.ForkChatResponse
		wantErr   bool
	}{
		{
			name: "正常系: チャットフォーク成功",
			args: args{
				chatUUID: "parent-chat-uuid",
				req: handlerModel.ForkChatRequest{
					TargetMessageUUID: "msg-uuid",
					ParentChatUUID:    "parent-chat-uuid",
					SelectedText:      "selected",
					RangeStart:        0,
					RangeEnd:          5,
					Title:             "New Chat",
					ContextSummary:    "Summary",
				},
			},
			setupMock: func(m *mocks) {
				m.chatUsecase.On("ForkChat", mock.Anything, mock.MatchedBy(func(params model.ForkChatParams) bool {
					return params.ParentChatUUID == "parent-chat-uuid" && params.Title == "New Chat"
				})).Return("new-chat-uuid", nil)
			},
			want: &handlerModel.ForkChatResponse{
				NewChatID: "new-chat-uuid",
				Message:   "子チャットを作成しました",
			},
			wantErr: false,
		},
		{
			name: "異常系: 親チャットID不一致",
			args: args{
				chatUUID: "parent-chat-uuid",
				req: handlerModel.ForkChatRequest{
					ParentChatUUID: "other-chat-uuid",
				},
			},
			setupMock: func(m *mocks) {
				// 呼ばれないはず
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系: Usecaseエラー",
			args: args{
				chatUUID: "parent-chat-uuid",
				req: handlerModel.ForkChatRequest{
					ParentChatUUID: "parent-chat-uuid",
				},
			},
			setupMock: func(m *mocks) {
				m.chatUsecase.On("ForkChat", mock.Anything, mock.Anything).Return("", errors.New("usecase error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mocks{
				chatUsecase: &MockChatUsecase{},
			}
			tt.setupMock(m)

			h := NewChatHandler(m.chatUsecase)

			reqBody, _ := json.Marshal(tt.args.req)
			req := httptest.NewRequest(http.MethodPost, "/api/chats/"+tt.args.chatUUID+"/fork", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/chats/:chat_uuid/fork")
			c.SetParamNames("chat_uuid")
			c.SetParamValues(tt.args.chatUUID)

			err := h.ForkChat(c)
			assert.NoError(t, err) // Handler should not return error, but write to response

			if tt.wantErr {
				assert.NotEqual(t, http.StatusOK, rec.Code)
			} else {
				assert.Equal(t, http.StatusOK, rec.Code)
				var res handlerModel.ForkChatResponse
				err = json.Unmarshal(rec.Body.Bytes(), &res)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, &res)
			}
		})
	}
}

func TestChatHandler_GetMergePreview(t *testing.T) {
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
			name: "正常系: マージプレビューが取得できること",
			args: args{
				chatUUID: "chat-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatUsecase.On("GetMergePreview", mock.Anything, "chat-uuid").Return(&model.MergePreview{
					SuggestedSummary: "summary",
				}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"suggested_summary":"summary"}`,
		},
		{
			name: "異常系: Usecaseがエラーを返した場合",
			args: args{
				chatUUID: "error-uuid",
			},
			setupMock: func(m *mocks) {
				m.chatUsecase.On("GetMergePreview", mock.Anything, "error-uuid").Return(nil, errors.New("usecase error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"status":"error","message":"usecase error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/chats/"+tt.args.chatUUID+"/merge/preview", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/chats/:chat_uuid/merge/preview")
			c.SetParamNames("chat_uuid")
			c.SetParamValues(tt.args.chatUUID)

			m := &mocks{
				chatUsecase: &MockChatUsecase{},
			}
			tt.setupMock(m)

			h := NewChatHandler(m.chatUsecase)
			err := h.GetMergePreview(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}
