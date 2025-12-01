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
	UserUUID  string    `gorm:"column:user_uuid;size:36"`
	Title     string    `gorm:"column:title;size:255"`
	CreatedID string    `gorm:"column:created_id;size:255"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// projectORMのテーブル名
func (projectORM) TableName() string {
	return "projects"
}

// projectORMをドメインモデルに変換する処理
func (orm *projectORM) toDomain() *model.Project {
	return &model.Project{
		UUID:      orm.UUID,
		UserUUID:  orm.UserUUID,
		Title:     orm.Title,
		CreatedAt: orm.CreatedAt,
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
	db := getDB(ctx, r.db)
	err := db.WithContext(ctx).Where("user_uuid = ?", userUUID).Find(&orms).Error
	if err != nil {
		return nil, err
	}

	projects := make([]*model.Project, len(orms))
	for i, orm := range orms {
		projects[i] = orm.toDomain()
	}

	return projects, nil
}

// プロジェクトを作成する処理
func (r *projectRepository) Create(ctx context.Context, project *model.Project) error {
	slog.DebugContext(ctx, "プロジェクト作成処理を開始", "project_uuid", project.UUID)
	orm := projectORM{
		UUID:      project.UUID,
		UserUUID:  project.UserUUID,
		Title:     project.Title,
		CreatedID: project.UserUUID,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	}
	db := getDB(ctx, r.db)
	return db.WithContext(ctx).Create(&orm).Error
}
