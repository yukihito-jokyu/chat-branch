package usecase

import (
	"backend/internal/domain/model"
	"context"
)

type AuthUsecase interface {
	// ゲストユーザーのサインアップ処理
	GuestSignup(ctx context.Context) (*model.User, string, error)
	// ゲストユーザーのログイン処理
	GuestLogin(ctx context.Context, userID string) (string, error)
}
