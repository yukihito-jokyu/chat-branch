package usecase

import (
	"backend/config"
	"backend/internal/domain/model"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// モックの定義
type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func TestAuthUsecase_GuestSignup(t *testing.T) {
	// テスト用の設定
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			Expiration: time.Hour,
		},
	}

	tests := []struct {
		name      string
		setupMock func(m *mockUserRepository)
		wantUser  bool // ユーザーが返されることを期待するか
		wantToken bool // トークンが返されることを期待するか
		wantErr   bool
	}{
		{
			name: "正常系: ゲストユーザー作成とトークン生成が成功すること",
			setupMock: func(m *mockUserRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
					return u.ID != "" && u.Name != "" // IDと名前が生成されていること
				})).Return(nil)
			},
			wantUser:  true,
			wantToken: true,
			wantErr:   false,
		},
		{
			name: "異常系: ユーザー作成に失敗した場合エラーになること",
			setupMock: func(m *mockUserRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantUser:  false,
			wantToken: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockUserRepository)
			tt.setupMock(mockRepo)

			u := NewAuthUsecase(mockRepo, cfg)
			user, token, err := u.GuestSignup(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("authUsecase.GuestSignup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantUser {
				assert.NotNil(t, user)
				assert.NotEmpty(t, user.ID)
			}
			if tt.wantToken {
				assert.NotEmpty(t, token)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthUsecase_GuestLogin(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			Expiration: time.Hour,
		},
	}

	tests := []struct {
		name      string
		userID    string
		setupMock func(m *mockUserRepository)
		wantToken bool
		wantErr   bool
	}{
		{
			name:   "正常系: 存在するユーザーでログインできること",
			userID: "test-user-id",
			setupMock: func(m *mockUserRepository) {
				m.On("FindByID", mock.Anything, "test-user-id").Return(&model.User{
					ID:   "test-user-id",
					Name: "Test User",
				}, nil)
			},
			wantToken: true,
			wantErr:   false,
		},
		{
			name:   "異常系: 存在しないユーザーの場合エラーになること",
			userID: "non-existent-id",
			setupMock: func(m *mockUserRepository) {
				m.On("FindByID", mock.Anything, "non-existent-id").Return(nil, errors.New("user not found"))
			},
			wantToken: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockUserRepository)
			tt.setupMock(mockRepo)

			u := NewAuthUsecase(mockRepo, cfg)
			token, err := u.GuestLogin(context.Background(), tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("authUsecase.GuestLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantToken {
				assert.NotEmpty(t, token)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
