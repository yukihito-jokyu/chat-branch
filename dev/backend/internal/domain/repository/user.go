package repository

import (
	"backend/internal/domain/model"
	"context"
)

type UserRepository interface {
	// 新しいユーザーをデータベースに作成する処理
	Create(ctx context.Context, user *model.User) error
	// 指定されたIDのユーザーを検索する処理
	FindByID(ctx context.Context, id string) (*model.User, error)
}
