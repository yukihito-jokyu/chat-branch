package router

import (
	"backend/config"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestInitRoutes(t *testing.T) {
	// テスト用のEchoインスタンスと依存関係を作成
	e := echo.New()
	db := &gorm.DB{}
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
		},
	}

	// ルーティングの初期化
	InitRoutes(e, db, cfg, nil, nil)

	// 期待されるルートの定義
	// 今後エンドポイントが増えた場合はここに追加する
	expectedRoutes := []struct {
		method string
		path   string
		name   string // ハンドラ名（完全一致でなくても、含まれているか確認する等でも可）
	}{
		{
			method: "GET",
			path:   "/health",
			name:   "HealthCheck",
		},
		{
			method: "POST",
			path:   "/api/auth/signup",
			name:   "Signup",
		},
		{
			method: "POST",
			path:   "/api/auth/login",
			name:   "Login",
		},
		{
			method: "GET",
			path:   "/api/projects",
			name:   "GetProjects",
		},
		{
			method: "POST",
			path:   "/api/projects",
			name:   "CreateProject",
		},
		{
			method: "GET",
			path:   "/api/chats/:chat_uuid",
			name:   "GetChat",
		},
		{
			method: "GET",
			path:   "/api/chats/:chat_uuid/messages",
			name:   "GetMessages",
		},
		{
			method: "POST",
			path:   "/api/chats/:chat_uuid/message",
			name:   "SendMessage",
		},
		{
			method: "GET",
			path:   "/api/chats/:chat_uuid/messages/stream",
			name:   "StreamMessage",
		},
		{
			method: "GET",
			path:   "/api/chats/:chat_uuid/stream",
			name:   "FirstStreamChat",
		},
		{
			method: "POST",
			path:   "/api/chats/:chat_uuid/fork/preview",
			name:   "GenerateForkPreview",
		},
		{
			method: "POST",
			path:   "/api/chats/:chat_uuid/fork",
			name:   "ForkChat",
		},
		{
			method: "POST",
			path:   "/api/chats/:chat_uuid/merge/preview",
			name:   "GetMergePreview",
		},
		{
			method: "POST",
			path:   "/api/chats/:chat_uuid/merge",
			name:   "MergeChat",
		},
		{
			method: "POST",
			path:   "/api/chats/:chat_uuid/close",
			name:   "CloseChat",
		},
		{
			method: "POST",
			path:   "/api/chats/:chat_uuid/open",
			name:   "OpenChat",
		},
	}

	// 登録されたルートを取得
	routes := e.Routes()

	// 期待されるルートが登録されているか検証
	for _, expected := range expectedRoutes {
		found := false
		for _, route := range routes {
			if route.Method == expected.method && route.Path == expected.path {
				// ハンドラ名の検証（リフレクション等で取得される名前は環境によって異なる可能性があるため、
				// ここでは簡易的にメソッド名が含まれているかなどを確認する方針も考えられるが、
				// EchoのRoute.Nameは登録時の関数名を返すため、それを検証する）
				if assert.Contains(t, route.Name, expected.name, "Route %s %s handler name mismatch", expected.method, expected.path) {
					found = true
					break
				}
			}
		}
		assert.True(t, found, "Route %s %s not found", expected.method, expected.path)
	}

	// 登録されたルートの総数が期待通りか検証（余計なルートが増えていないか）
	// ミドルウェアやデフォルトのルートが含まれる可能性があるため、厳密な一致ではなく
	// 「期待されるルート数以上であること」を確認するか、あるいは
	// アプリケーションで定義したルートのみをカウントするロジックを入れる。
	// ここではシンプルに、期待されるルートが全て見つかったことを以て良しとする。
}
