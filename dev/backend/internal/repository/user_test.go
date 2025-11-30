package repository

import (
	"backend/internal/domain/model"
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserRepository_Create(t *testing.T) {
	type args struct {
		user *model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "正常系: ユーザーが作成できること",
			args: args{
				user: &model.User{
					ID:   "test-uuid",
					Name: "test-user",
				},
			},
			wantErr: false,
		},
		{
			name: "異常系: IDが重複している場合エラーになること",
			args: args{
				user: &model.User{
					ID:   "duplicate-uuid",
					Name: "duplicate-user",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// インメモリDBのセットアップ
			db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to connect database: %v", err)
			}
			// マイグレーション
			if err := db.AutoMigrate(&userORM{}); err != nil {
				t.Fatalf("failed to migrate database: %v", err)
			}

			r := NewUserRepository(db)

			// 重複エラーのテストケースのための事前データ投入
			if tt.name == "異常系: IDが重複している場合エラーになること" {
				if err := r.Create(context.Background(), tt.args.user); err != nil {
					t.Fatalf("failed to create pre-existing user: %v", err)
				}
			}

			if err := r.Create(context.Background(), tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name      string
		args      args
		setupData func(db *gorm.DB)
		want      *model.User
		wantErr   bool
	}{
		{
			name: "正常系: 存在するユーザーを取得できること",
			args: args{
				id: "test-uuid",
			},
			setupData: func(db *gorm.DB) {
				db.Create(&userORM{
					ID:        "test-uuid",
					Name:      "test-user",
					CreatedID: "test-uuid",
				})
			},
			want: &model.User{
				ID:   "test-uuid",
				Name: "test-user",
			},
			wantErr: false,
		},
		{
			name: "異常系: 存在しないユーザーの場合エラーになること",
			args: args{
				id: "non-existent-uuid",
			},
			setupData: func(db *gorm.DB) {},
			want:      nil,
			wantErr:   true,
		},
		{
			name: "異常系: データベースエラーの場合エラーになること",
			args: args{
				id: "test-uuid",
			},
			setupData: func(db *gorm.DB) {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// インメモリDBのセットアップ
			db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to connect database: %v", err)
			}
			// マイグレーション
			if err := db.AutoMigrate(&userORM{}); err != nil {
				t.Fatalf("failed to migrate database: %v", err)
			}

			// データセットアップ
			tt.setupData(db)

			r := NewUserRepository(db)
			got, err := r.FindByUUID(context.Background(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepository.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
