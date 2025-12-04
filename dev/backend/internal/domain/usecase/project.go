package usecase

import (
	"backend/internal/domain/model"
	"context"
)

type ProjectUsecase interface {
	// プロジェクト取得処理
	GetProjects(ctx context.Context, userUUID string) ([]*model.Project, error)
	// プロジェクト作成処理
	CreateProject(ctx context.Context, userUUID, initialMessage string) (*model.Project, *model.Chat, *model.Message, error)
	// プロジェクトの親チャット取得処理
	GetParentChat(ctx context.Context, projectUUID string) (*model.Chat, error)
	// プロジェクトツリー取得処理
	GetProjectTree(ctx context.Context, projectUUID string) (*model.ProjectTree, error)
}
