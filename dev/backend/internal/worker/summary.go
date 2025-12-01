package worker

import (
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	"backend/internal/domain/usecase"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
	"google.golang.org/genai"
)

type SummaryWorker struct {
	subscriber  message.Subscriber
	messageRepo repository.MessageRepository
	genaiClient usecase.GenAIClient
}

func NewSummaryWorker(subscriber message.Subscriber, messageRepo repository.MessageRepository, genaiClient usecase.GenAIClient) *SummaryWorker {
	return &SummaryWorker{
		subscriber:  subscriber,
		messageRepo: messageRepo,
		genaiClient: genaiClient,
	}
}

// 要約生成タスクの起動
func (w *SummaryWorker) Run(ctx context.Context) error {
	messages, err := w.subscriber.Subscribe(ctx, "chat_summary")
	if err != nil {
		return fmt.Errorf("failed to subscribe to chat_summary: %w", err)
	}

	for msg := range messages {
		if err := w.Handle(ctx, msg); err != nil {
			slog.ErrorContext(ctx, "要約生成タスクの処理に失敗", "error", err)
			msg.Nack()
		} else {
			msg.Ack()
		}
	}

	return nil
}

// 要約生成タスクの処理
func (w *SummaryWorker) Handle(ctx context.Context, msg *message.Message) error {
	var chatUUID string
	if err := json.Unmarshal(msg.Payload, &chatUUID); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	slog.InfoContext(ctx, "要約生成タスク開始", "chat_uuid", chatUUID)

	// 1. 最新のサマリを持つメッセージを取得
	latestSummaryMessage, err := w.messageRepo.FindLatestMessageWithSummary(ctx, chatUUID)
	if err != nil {
		return fmt.Errorf("failed to fetch latest summary message: %w", err)
	}

	// 2. メッセージ履歴の取得
	allMessages, err := w.messageRepo.FindMessagesByChatID(ctx, chatUUID)
	if err != nil {
		return fmt.Errorf("failed to fetch messages: %w", err)
	}

	if len(allMessages) == 0 {
		slog.WarnContext(ctx, "メッセージが存在しません", "chat_uuid", chatUUID)
		return nil
	}

	// 3. サマリ以降のメッセージを抽出
	var targetMessages []*model.Message
	if latestSummaryMessage != nil {
		found := false
		for _, m := range allMessages {
			if found {
				targetMessages = append(targetMessages, m)
			}
			if m.UUID == latestSummaryMessage.UUID {
				found = true
			}
		}
	} else {
		targetMessages = allMessages
	}

	if len(targetMessages) == 0 {
		slog.InfoContext(ctx, "要約対象の新規メッセージがありません", "chat_uuid", chatUUID)
		return nil
	}

	// 4. プロンプト構築
	var parts []*genai.Content

	// ベースとなるサマリがある場合
	if latestSummaryMessage != nil && latestSummaryMessage.ContextSummary != nil {
		parts = append(parts, &genai.Content{
			Role: "user",
			Parts: []*genai.Part{
				{Text: "これまでの会話の要約:\n" + *latestSummaryMessage.ContextSummary},
			},
		})
	}

	for _, m := range targetMessages {
		role := "user"
		if m.Role == "assistant" {
			role = "model"
		}
		parts = append(parts, &genai.Content{
			Role: role,
			Parts: []*genai.Part{
				{Text: m.Content},
			},
		})
	}

	// 要約指示
	prompt := "上記の会話（これまでの要約を含む）を、次の会話のコンテキストとして使用できるように要約してください。"
	parts = append(parts, &genai.Content{
		Role: "user",
		Parts: []*genai.Part{
			{Text: prompt},
		},
	})

	// 5. GenAI 呼び出し
	client := w.genaiClient
	modelName := "gemini-2.5-flash"

	resp, err := client.GenerateContent(ctx, modelName, parts, nil)
	if err != nil {
		return fmt.Errorf("genai error: %w", err)
	}

	var summary string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				summary += part.Text
			}
		}
	}

	if summary == "" {
		return fmt.Errorf("empty summary generated")
	}

	// 6. 最新のメッセージの context_summary を更新
	// targetMessagesの最後ではなく、allMessagesの最後（＝最新のメッセージ）に紐づける
	lastMessage := allMessages[len(allMessages)-1]

	if err := w.messageRepo.UpdateContextSummary(ctx, lastMessage.UUID, summary); err != nil {
		return fmt.Errorf("failed to update context summary: %w", err)
	}

	slog.InfoContext(ctx, "要約生成完了", "chat_uuid", chatUUID, "summary_length", len(summary))
	return nil
}
