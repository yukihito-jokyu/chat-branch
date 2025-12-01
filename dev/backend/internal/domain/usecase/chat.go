package usecase

import (
	"backend/internal/domain/model"
	"context"

	"google.golang.org/genai"
)

type ChatUsecase interface {
	// チャットの最初のメッセージを元に、GenAI にストリームを送信する
	FirstStreamChat(ctx context.Context, chatUUID string, outputChan chan<- string) error
	// チャットを取得する
	GetChat(ctx context.Context, chatUUID string) (*model.Chat, error)
	// チャットのメッセージ一覧を取得する
	GetMessages(ctx context.Context, chatUUID string) ([]*model.Message, error)
	// メッセージを送信する
	SendMessage(ctx context.Context, chatUUID string, content string) (*model.Message, error)
	// メッセージをストリーミング送信する
	StreamMessage(ctx context.Context, chatUUID string, outputChan chan<- string) error
	// フォークプレビューを生成する
	GenerateForkPreview(ctx context.Context, chatUUID string, req model.ForkPreviewRequest) (*model.ForkPreviewResponse, error)
	// チャットをフォークする
	ForkChat(ctx context.Context, params model.ForkChatParams) (string, error)
}

// GenAIClient は GenAI クライアントのインターフェース
// モック化のために定義
type GenAIClient interface {
	GenerateContentStream(ctx context.Context, model string, parts []*genai.Content, config *genai.GenerateContentConfig) func(func(*genai.GenerateContentResponse, error) bool)
	GenerateContent(ctx context.Context, model string, parts []*genai.Content, config *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error)
}
