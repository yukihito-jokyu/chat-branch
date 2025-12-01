package usecase

import (
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	domainUsecase "backend/internal/domain/usecase"
	"backend/internal/infrastructure/queue"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"google.golang.org/genai"
)

// GenAIClientWrapper は GenAIClient インターフェースの実装
type GenAIClientWrapper struct {
	client *genai.Client
}

func NewGenAIClientWrapper(client *genai.Client) domainUsecase.GenAIClient {
	return &GenAIClientWrapper{client: client}
}

func (w *GenAIClientWrapper) GenerateContentStream(ctx context.Context, model string, parts []*genai.Content, config *genai.GenerateContentConfig) func(func(*genai.GenerateContentResponse, error) bool) {
	return w.client.Models.GenerateContentStream(ctx, model, parts, config)
}

func (w *GenAIClientWrapper) GenerateContent(ctx context.Context, model string, parts []*genai.Content, config *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	return w.client.Models.GenerateContent(ctx, model, parts, config)
}

type chatUsecase struct {
	chatRepo    repository.ChatRepository
	messageRepo repository.MessageRepository
	genaiClient domainUsecase.GenAIClient
	publisher   message.Publisher
}

func NewChatUsecase(chatRepo repository.ChatRepository, messageRepo repository.MessageRepository, genaiClient domainUsecase.GenAIClient, publisher message.Publisher) domainUsecase.ChatUsecase {
	return &chatUsecase{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
		genaiClient: genaiClient,
		publisher:   publisher,
	}
}

// ユーザーのメッセージを元に、GenAI にストリームを送信する
func (u *chatUsecase) StreamMessage(ctx context.Context, chatUUID string, outputChan chan<- string) error {
	slog.InfoContext(ctx, "メッセージストリーム処理開始", "chat_uuid", chatUUID)

	// 1. 最新のサマリを持つメッセージを取得
	latestSummaryMessage, err := u.messageRepo.FindLatestMessageWithSummary(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "最新サマリ取得失敗", "chat_uuid", chatUUID, "error", err)
		return err
	}

	// 2. メッセージ履歴の取得
	allMessages, err := u.messageRepo.FindMessagesByChatID(ctx, chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "メッセージ履歴取得失敗", "chat_uuid", chatUUID, "error", err)
		return err
	}

	var contextMessages []*model.Message
	if latestSummaryMessage != nil {
		// サマリ以降のメッセージを抽出
		found := false
		for _, msg := range allMessages {
			if found {
				contextMessages = append(contextMessages, msg)
			}
			if msg.UUID == latestSummaryMessage.UUID {
				found = true
			}
		}
	} else {
		contextMessages = allMessages
	}

	// 3. プロンプト構築
	var parts []*genai.Content

	// サマリがあれば追加
	if latestSummaryMessage != nil && latestSummaryMessage.ContextSummary != nil {
		parts = append(parts, &genai.Content{
			Role: "user",
			Parts: []*genai.Part{
				{Text: "以下の会話の要約を踏まえて回答してください:\n" + *latestSummaryMessage.ContextSummary},
			},
		})
	}

	for _, msg := range contextMessages {
		role := "user"
		if msg.Role == "assistant" {
			role = "model"
		}
		parts = append(parts, &genai.Content{
			Role: role,
			Parts: []*genai.Part{
				{Text: msg.Content},
			},
		})
	}

	// 4. GenAI 呼び出し
	client := u.genaiClient
	modelName := "gemini-2.5-flash"

	if len(contextMessages) == 0 && latestSummaryMessage == nil {
		return errors.New("no context to generate response")
	}

	if len(parts) == 0 {
		return errors.New("empty prompt")
	}

	iter := client.GenerateContentStream(ctx, modelName, parts, nil)

	var fullResponse string
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

	// 5. 生成された文章の保存
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

	// 6. サマリ生成タスクのPublish
	topic := "chat_summary"
	payload, err := json.Marshal(chatUUID)
	if err != nil {
		slog.ErrorContext(ctx, "payloadのJSON変換に失敗しました", "error", err)
		return nil // 非同期タスクの失敗はメイン処理のエラーにはしない
	}

	if err := queue.PublishTask(u.publisher, topic, payload); err != nil {
		slog.ErrorContext(ctx, "サマリ生成タスクの登録に失敗しました", "error", err)
	} else {
		slog.InfoContext(ctx, "サマリ生成タスクを登録しました", "chat_uuid", chatUUID)
	}

	slog.InfoContext(ctx, "メッセージストリーム処理完了", "chat_uuid", chatUUID)
	return nil
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
