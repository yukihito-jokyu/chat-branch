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
}
