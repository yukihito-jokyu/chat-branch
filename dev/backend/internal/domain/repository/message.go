package repository

import (
	"backend/internal/domain/model"
	"context"
)

type MessageRepository interface {
	// メッセージを作成する処理
	Create(ctx context.Context, message *model.Message) error
	// チャットIDに紐づくメッセージを取得する処理
	FindMessagesByChatID(ctx context.Context, chatUUID string) ([]*model.Message, error)
	// メッセージのコンテキストサマリを更新する処理
	UpdateContextSummary(ctx context.Context, messageUUID string, summary string) error
	// 指定されたチャットIDの中で、コンテキストサマリを持つ最新のメッセージを取得する処理
	FindLatestMessageWithSummary(ctx context.Context, chatUUID string) (*model.Message, error)
	FindLatestMessageByRole(ctx context.Context, chatUUID string, role string) (*model.Message, error)
}
