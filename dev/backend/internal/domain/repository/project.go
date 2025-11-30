package repository

import (
	"backend/internal/domain/model"
	"context"
)

type ProjectRepository interface {
	// ユーザーUUIDでプロジェクト一覧を取得する処理
	FindAllByUserUUID(ctx context.Context, userUUID string) ([]*model.Project, error)
}
