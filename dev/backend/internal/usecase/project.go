package usecase

import (
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	"backend/internal/domain/usecase"
	"context"
	"fmt"
	"log/slog"
)

type projectUsecase struct {
	projectRepo repository.ProjectRepository
}

func NewProjectUsecase(projectRepo repository.ProjectRepository) usecase.ProjectUsecase {
	return &projectUsecase{
		projectRepo: projectRepo,
	}
}

// プロジェクト取得処理
func (u *projectUsecase) GetProjects(ctx context.Context, userUUID string) ([]*model.Project, error) {
	slog.InfoContext(ctx, "プロジェクト一覧取得処理を開始", "user_uuid", userUUID)
	projects, err := u.projectRepo.FindAllByUserUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("プロジェクト一覧の取得に失敗: %w", err)
	}
	slog.InfoContext(ctx, "プロジェクト一覧取得処理を完了", "user_uuid", userUUID)
	return projects, nil
}
