package repository

import (
	"backend/internal/domain/model"
	"context"
)

type UserRepository interface {
	// 新しいユーザーをデータベースに作成する処理
	Create(ctx context.Context, user *model.User) error
	// 指定されたUUIDのユーザーを検索する処理
	FindByUUID(ctx context.Context, uuid string) (*model.User, error)
}
