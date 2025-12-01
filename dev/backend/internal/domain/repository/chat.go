package repository

import (
	"backend/internal/domain/model"
	"context"
)

type ChatRepository interface {
	// チャットを作成する処理
	Create(ctx context.Context, chat *model.Chat) error
	// チャットを取得する処理
	FindByID(ctx context.Context, uuid string) (*model.Chat, error)
}
