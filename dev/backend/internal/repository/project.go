package repository

import (
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

// projectのORMモデル
type projectORM struct {
	UUID      string    `gorm:"primaryKey;column:uuid;size:36"`
	UserID    string    `gorm:"column:user_id;size:36"`
	Title     string    `gorm:"size:255"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// projectORMのテーブル名
func (projectORM) TableName() string {
	return "projects"
}

// projectORMをドメインモデルに変換する処理
func (orm *projectORM) toDomain() *model.Project {
	return &model.Project{
		ID:        orm.UUID,
		UserID:    orm.UserID,
		Title:     orm.Title,
		UpdatedAt: orm.UpdatedAt,
	}
}

type projectRepository struct {
	db *gorm.DB
}

// projectRepositoryの新しいインスタンスを作成する処理
func NewProjectRepository(db *gorm.DB) repository.ProjectRepository {
	return &projectRepository{db: db}
}

// ユーザーUUIDでプロジェクト一覧を取得する処理
func (r *projectRepository) FindAllByUserUUID(ctx context.Context, userUUID string) ([]*model.Project, error) {
	slog.DebugContext(ctx, "プロジェクト一覧取得処理を開始", "user_uuid", userUUID)
	var orms []projectORM
	err := r.db.WithContext(ctx).Where("user_id = ?", userUUID).Find(&orms).Error
	if err != nil {
		return nil, err
	}

	projects := make([]*model.Project, len(orms))
	for i, orm := range orms {
		projects[i] = orm.toDomain()
	}

	return projects, nil
}
