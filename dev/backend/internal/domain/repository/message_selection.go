package repository

import (
	"backend/internal/domain/model"
	"context"
)

type MessageSelectionRepository interface {
	// メッセージ選択範囲を作成する処理
	Create(ctx context.Context, selection *model.MessageSelection) error
}
