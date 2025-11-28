package usecase

import (
	"backend/config"
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	"backend/internal/domain/usecase"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type authUsecase struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

// AuthUsecase の新しいインスタンスを作成する処理
func NewAuthUsecase(userRepo repository.UserRepository, cfg *config.Config) usecase.AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// ゲストユーザーのサインアップ処理
func (u *authUsecase) GuestSignup(ctx context.Context) (*model.User, string, error) {
	slog.InfoContext(ctx, "ゲストサインアップ処理を開始")
	// ランダムなユーザーを生成
	userID := uuid.New().String()
	user := &model.User{
		ID:   userID,
		Name: "Guest-" + userID[:8],
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, "", fmt.Errorf("ユーザー作成に失敗: %w", err)
	}

	token, err := u.generateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("トークン生成に失敗: %w", err)
	}

	slog.InfoContext(ctx, "ゲストユーザーを作成しました", "user_id", user.ID)
	return user, token, nil
}

// ゲストユーザーのログイン処理
func (u *authUsecase) GuestLogin(ctx context.Context, userID string) (string, error) {
	slog.InfoContext(ctx, "ゲストログイン処理を開始", "user_id", userID)
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("ユーザー検索に失敗: %w", err)
	}

	token, err := u.generateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("トークン生成に失敗: %w", err)
	}

	return token, nil
}

// 指定されたユーザーIDのJWTトークンを生成する処理
func (u *authUsecase) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(u.cfg.JWT.Expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.cfg.JWT.Secret))
}
