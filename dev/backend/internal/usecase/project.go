package usecase

import (
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	"backend/internal/domain/usecase"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type projectUsecase struct {
	projectRepo repository.ProjectRepository
	chatRepo    repository.ChatRepository
	messageRepo repository.MessageRepository
	txManager   repository.TransactionManager
}

func NewProjectUsecase(
	projectRepo repository.ProjectRepository,
	chatRepo repository.ChatRepository,
	messageRepo repository.MessageRepository,
	txManager repository.TransactionManager,
) usecase.ProjectUsecase {
	return &projectUsecase{
		projectRepo: projectRepo,
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
		txManager:   txManager,
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

// プロジェクト作成処理
func (u *projectUsecase) CreateProject(ctx context.Context, userUUID, initialMessage string) (*model.Project, *model.Chat, *model.Message, error) {
	slog.InfoContext(ctx, "プロジェクト作成処理を開始", "user_uuid", userUUID)

	projectID := uuid.New().String()
	chatID := uuid.New().String()
	messageID := uuid.New().String()
	now := time.Now()

	project := &model.Project{
		UUID:      projectID,
		UserUUID:  userUUID,
		Title:     initialMessage, // タイトルは最初のメッセージとする（要件によるが一旦）
		CreatedAt: now,
		UpdatedAt: now,
	}

	chat := &model.Chat{
		UUID:        chatID,
		ProjectUUID: projectID,
		UserUUID:    userUUID,
		Title:       initialMessage, // チャットのタイトルも最初のメッセージとする
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	message := &model.Message{
		UUID:      messageID,
		ChatUUID:  chatID,
		UserUUID:  userUUID,
		Role:      "user",
		Content:   initialMessage,
		CreatedAt: now,
	}

	err := u.txManager.Do(ctx, func(ctx context.Context) error {
		if err := u.projectRepo.Create(ctx, project); err != nil {
			return fmt.Errorf("プロジェクト作成失敗: %w", err)
		}
		if err := u.chatRepo.Create(ctx, chat); err != nil {
			return fmt.Errorf("チャット作成失敗: %w", err)
		}
		if err := u.messageRepo.Create(ctx, message); err != nil {
			return fmt.Errorf("メッセージ作成失敗: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, nil, nil, err
	}

	slog.InfoContext(ctx, "プロジェクト作成処理を完了", "project_id", projectID)
	return project, chat, message, nil
}
