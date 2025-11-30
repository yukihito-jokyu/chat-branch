package repository

import (
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	"context"
	"errors"

	"log/slog"

	"gorm.io/gorm"
)

// userのORMモデル
type userORM struct {
	ID        string `gorm:"primaryKey;column:uuid;size:36"`
	Name      string `gorm:"size:255"`
	CreatedID string `gorm:"column:created_id;size:255"`
}

// テーブル名を返す処理
func (userORM) TableName() string {
	return "users"
}

// userORMをドメインモデルに変換する処理
func (orm *userORM) toDomain() *model.User {
	return &model.User{
		ID:   orm.ID,
		Name: orm.Name,
	}
}

// ドメインモデルをuserORMに変換する処理
func fromDomain(u *model.User) *userORM {
	return &userORM{
		ID:        u.ID,
		Name:      u.Name,
		CreatedID: u.ID,
	}
}

type userRepository struct {
	db *gorm.DB
}

// UserRepository の新しいインスタンスを作成する処理
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// 新しいユーザーをデータベースに作成する処理
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	slog.DebugContext(ctx, "ユーザー作成処理を開始", "user_id", user.ID)
	orm := fromDomain(user)
	err := r.db.WithContext(ctx).Create(orm).Error
	if err != nil {
		return err
	}
	return nil
}

// 指定されたUUIDのユーザーを検索する処理
func (r *userRepository) FindByUUID(ctx context.Context, uuid string) (*model.User, error) {
	slog.DebugContext(ctx, "ユーザー検索処理を開始", "user_uuid", uuid)
	var orm userORM
	err := r.db.WithContext(ctx).First(&orm, "uuid = ?", uuid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return orm.toDomain(), nil
}
