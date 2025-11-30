package usecase

import (
	"backend/internal/domain/model"
	"context"
)

type ProjectUsecase interface {
	// プロジェクト取得処理
	GetProjects(ctx context.Context, userUUID string) ([]*model.Project, error)
}
