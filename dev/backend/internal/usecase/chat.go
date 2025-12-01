package usecase

import (
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"google.golang.org/genai"
)

// GenAIClient は GenAI クライアントのインターフェース
// モック化のために定義
type GenAIClient interface {
	GenerateContentStream(ctx context.Context, model string, parts []*genai.Content, config *genai.GenerateContentConfig) func(func(*genai.GenerateContentResponse, error) bool)
}

// GenAIClientWrapper は GenAIClient インターフェースの実装
type GenAIClientWrapper struct {
	client *genai.Client
}

func NewGenAIClientWrapper(client *genai.Client) GenAIClient {
	return &GenAIClientWrapper{client: client}
}

func (w *GenAIClientWrapper) GenerateContentStream(ctx context.Context, model string, parts []*genai.Content, config *genai.GenerateContentConfig) func(func(*genai.GenerateContentResponse, error) bool) {
	return w.client.Models.GenerateContentStream(ctx, model, parts, config)
}

type chatUsecase struct {
	chatRepo    repository.ChatRepository
	messageRepo repository.MessageRepository
	genaiClient GenAIClient
}

func NewChatUsecase(chatRepo repository.ChatRepository, messageRepo repository.MessageRepository, genaiClient GenAIClient) *chatUsecase {
	return &chatUsecase{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
		genaiClient: genaiClient,
	}
}

// チャットの最初のメッセージを元に、GenAI にストリームを送信する
func (u *chatUsecase) FirstStreamChat(ctx context.Context, chatUUID string, outputChan chan<- string) error {
	slog.InfoContext(ctx, "チャットストリーム処理開始", "chat_uuid", chatUUID)

	// 1. チャットの存在確認
	_, err := u.chatRepo.FindByID(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "チャットが見つかりません", "chat_uuid", chatUUID, "error", err)
		return err
	}

	// 2. メッセージ履歴の取得
	messages, err := u.messageRepo.FindMessagesByChatID(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "メッセージ履歴の取得に失敗しました", "chat_uuid", chatUUID, "error", err)
		return err
	}

	if len(messages) != 1 {
		return errors.New("invalid message state: expected exactly one initial message")
	}

	// 3. GenAI クライアント (注入されたものを使用)
	client := u.genaiClient

	modelName := "gemini-2.5-flash"

	// 4. プロンプトの構築とストリーム送信
	targetMessage := messages[0]

	// GenerateContentStream の呼び出し
	iter := client.GenerateContentStream(ctx, modelName, genai.Text(targetMessage.Content), nil)

	var fullResponse string

	// 5. ストリーム処理
	for resp, err := range iter {
		if err != nil {
			slog.ErrorContext(ctx, "GenAI APIからの受信エラー", "error", err)
			return err
		}

		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					chunk := part.Text
					fullResponse += chunk
					outputChan <- chunk
				}
			}
		}
	}

	// 6. 生成された文章の保存
	assistantMessage := &model.Message{
		UUID:      uuid.New().String(),
		ChatUUID:  chatUUID,
		Role:      "assistant",
		Content:   fullResponse,
		CreatedAt: time.Now(),
	}

	if err := u.messageRepo.Create(ctx, assistantMessage); err != nil {
		slog.ErrorContext(ctx, "アシスタントメッセージの保存に失敗しました", "error", err)
		return err
	}

	slog.InfoContext(ctx, "チャットストリーム処理完了", "chat_uuid", chatUUID)
	return nil
}

// チャットを取得する
func (u *chatUsecase) GetChat(ctx context.Context, chatUUID string) (*model.Chat, error) {
	slog.InfoContext(ctx, "チャット取得処理開始", "chat_uuid", chatUUID)
	chat, err := u.chatRepo.FindByID(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "チャット取得失敗", "chat_uuid", chatUUID, "error", err)
		return nil, err
	}
	return chat, nil
}

// チャットのメッセージ一覧を取得する
func (u *chatUsecase) GetMessages(ctx context.Context, chatUUID string) ([]*model.Message, error) {
	slog.InfoContext(ctx, "メッセージ一覧取得処理開始", "chat_uuid", chatUUID)
	messages, err := u.messageRepo.FindMessagesByChatID(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "メッセージ一覧取得失敗", "chat_uuid", chatUUID, "error", err)
		return nil, err
	}
	return messages, nil
}

// メッセージを送信する
func (u *chatUsecase) SendMessage(ctx context.Context, chatUUID string, content string) (*model.Message, error) {
	slog.InfoContext(ctx, "メッセージ送信処理開始", "chat_uuid", chatUUID)

	// チャットの存在確認
	_, err := u.chatRepo.FindByID(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "チャットが見つかりません", "chat_uuid", chatUUID, "error", err)
		return nil, err
	}

	message := &model.Message{
		UUID:      uuid.New().String(),
		ChatUUID:  chatUUID,
		Role:      "user",
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := u.messageRepo.Create(ctx, message); err != nil {
		slog.ErrorContext(ctx, "メッセージ保存失敗", "chat_uuid", chatUUID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "メッセージ送信成功", "message_uuid", message.UUID)
	return message, nil
}
