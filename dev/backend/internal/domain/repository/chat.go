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
	// チャットのステータスを更新する処理
	UpdateStatus(ctx context.Context, chatUUID string, status string) error
	// プロジェクト内で最も古いチャットを取得する処理
	FindOldestByProjectUUID(ctx context.Context, projectUUID string) (*model.Chat, error)
}
